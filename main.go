package main

import (
	"fmt"
	"goputer/internal/components"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type Panel interface {
	tea.Model
	ToggleActive()
}

type model struct {
	OS                 string
	Time               time.Time
	User               string
	panels             []Panel
	selectedPanelIndex int
	debugMode          bool
	width, height      int
}

func initialModel() model {
	user, _ := user.Current()
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	var panels []Panel
	panels = append(
		panels,
		components.MakeCpuModel(width/2, height),
		components.MakeMemoryModel(width/2, height),
		components.MakeDiskModel(width/2, height),
		components.MakeProcessesModel(width/2, height),
	)
	panels[0].ToggleActive()

	return model{
		OS:        runtime.GOOS,
		Time:      time.Now(),
		User:      user.Username,
		debugMode: false,
		width:     width,
		height:    height,
		panels:    panels,
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
			m.panels[m.selectedPanelIndex].ToggleActive()
			if len(m.panels) == m.selectedPanelIndex+1 {
				m.selectedPanelIndex = 0
			} else {
				m.selectedPanelIndex += 1
			}
			m.panels[m.selectedPanelIndex].ToggleActive()
		case "left":
			m.panels[m.selectedPanelIndex].ToggleActive()
			if m.selectedPanelIndex == 0 {
				m.selectedPanelIndex = len(m.panels) - 1
			} else {
				m.selectedPanelIndex -= 1
			}
			m.panels[m.selectedPanelIndex].ToggleActive()
		case "del":
			m.selectedPanelIndex = 0
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		for _, pane := range m.panels {
			pane.Update(msg)
		}
	}

	if _, isKeyMsg := msg.(tea.KeyMsg); !isKeyMsg {
		for i, panel := range m.panels {
			var cmd tea.Cmd
			model, cmd := panel.Update(msg)
			m.panels[i] = model.(Panel)

			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	} else {
		// Only send the key messages to the active pane
		_, cmd := m.panels[m.selectedPanelIndex].Update(msg)
		cmds = append(cmds, cmd)
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

	// Create 2x2 grid
	var panelsView string
	if len(m.panels) >= 4 {
		// First row: panels 0 and 1
		topRow := lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.panels[0].View(),
			m.panels[1].View(),
		)

		// Second row: panels 2 and 3
		bottomRow := lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.panels[2].View(),
			m.panels[3].View(),
		)

		// Join rows vertically
		panelsView = lipgloss.JoinVertical(
			lipgloss.Left,
			topRow,
			bottomRow,
		)
	}

	return header + "\n" + panelsView
}

func (m model) Init() tea.Cmd {
	var batch []tea.Cmd
	batch = append(
		batch,
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
