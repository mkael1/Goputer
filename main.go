package main

import (
	"fmt"
	"goputer/internal/components"
	"goputer/internal/system"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type model struct {
	OS            string
	Time          time.Time
	User          string
	Cpu           system.Cpu
	panels        []tea.Model
	selectedPanel tea.Model
	debugMode     bool
	width, height int
}

func initialModel() model {
	user, _ := user.Current()
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	var panels []tea.Model
	return model{
		OS:        runtime.GOOS,
		Time:      time.Now(),
		User:      user.Username,
		debugMode: false,
		width:     width,
		height:    height,
		panels: append(
			panels,
			components.MakeMemoryModel(width, height),
			components.MakeDiskModel(width, height),
			components.MakeProcessesModel(width, height),
		),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		case "d":
			m.debugMode = !m.debugMode
			return m, nil
		case "right":
			m.selectedPanel = m.panels[1]
		case "del":
			m.selectedPanel = nil
		}

	case system.CpuMsg:
		m.Cpu = system.Cpu(msg)
		return m, system.CheckCpu()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	for i, panel := range m.panels {
		var cmd tea.Cmd

		if _, isKeyMsg := msg.(tea.KeyMsg); !isKeyMsg {
			// Not a key message, send to all panels
			m.panels[i], cmd = panel.Update(msg)
		} else if panel == m.selectedPanel {
			// Is a key message, only send to selected panel
			m.panels[i], cmd = panel.Update(msg)
		}

		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	rightSideHeader := strings.ToUpper(m.OS) + " | " + m.Time.Format(time.Kitchen)

	if m.debugMode {
		var b strings.Builder

		return b.String()
	}

	var headerStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ffffff")).
		BorderBottom(true).
		Width(m.width).
		MarginBottom(1)

	s := fmt.Sprintf("System Monitor Resources - %s", m.User)

	// Create the full header with left and right content
	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		s,
		strings.Repeat(" ", m.width-lipgloss.Width(s)-lipgloss.Width(rightSideHeader)),
		rightSideHeader,
	)

	header = headerStyle.Render(header)
	return header + lipgloss.JoinVertical(
		lipgloss.Top,
		m.panels[0].View(),
		m.panels[1].View(),
		m.panels[2].View(),
	)
}

func (m model) Init() tea.Cmd {
	var batch []tea.Cmd
	batch = append(
		batch,
		system.CheckCpu(),
	)

	for _, panel := range m.panels {
		batch = append(batch, panel.Init())
	}

	return tea.Batch(
		batch...,
	)
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Ahhh... there's been an error: %v", err)
		os.Exit(1)
	}
}
