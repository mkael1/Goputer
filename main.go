package main

import (
	"fmt"
	"goputer/internal/components/cpu"
	"goputer/internal/components/disk"
	"goputer/internal/components/keybindings"
	"goputer/internal/components/memory"
	"goputer/internal/components/processes"
	"goputer/internal/keys"
	"goputer/internal/panel"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type model struct {
	OS            string
	User          string
	width, height int
	panelManager  *panel.PanelManager
	showHelp      bool
}

func initialModel() model {
	user, _ := user.Current()
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	var panels []panel.Panel
	panels = append(
		panels,
		cpu.MakeCpuModel(width/2, height),
		memory.MakeMemoryModel(width/2, height),
		disk.MakeDiskModel(width/2, height),
		processes.MakeProcessesModel(width/2, height),
	)

	var cells []panel.Cell
	for _, p := range panels {
		cells = append(cells, panel.Cell{
			Panel:   p,
			ColSpan: 6,
		})
	}

	model := model{
		OS:           runtime.GOOS,
		User:         user.Username,
		width:        width,
		height:       height,
		panelManager: panel.NewPanelManager(width, height, panel.WithCols(12), panel.WithCells(cells)),
		showHelp:     false,
	}

	return model
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Global.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Global.Help):
			m.showHelp = !m.showHelp
		default:
			_, cmd := m.panelManager.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		_, cmd := m.panelManager.Update(msg)
		cmds = append(cmds, cmd)
	}

	if _, isKeyMsg := msg.(tea.KeyMsg); !isKeyMsg {
		_, cmd := m.panelManager.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var headerStyle = lipgloss.NewStyle().
		Width(m.width)

	leftHeader := fmt.Sprintf("System Monitor Resources")
	rightHeader := fmt.Sprintf("%s | %s | %s", m.User, strings.ToUpper(m.OS), time.Now().Format(time.Kitchen))

	// Create the full header with left and right content
	header := headerStyle.Render(lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftHeader,
		strings.Repeat(" ", m.width-lipgloss.Width(leftHeader)-lipgloss.Width(rightHeader)-headerStyle.GetHorizontalPadding()),
		rightHeader,
	))

	var footer string
	if !m.showHelp {
		footer = lipgloss.NewStyle().BorderTop(true).BorderBottom(false).Render("? toggle help * q quit")
	} else {
		footer = keybindings.MakeKeybindingsModel(m.width, m.height).View()
	}

	content := m.panelManager.View()
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		footer,
	)
}

func (m model) Init() tea.Cmd {
	var batch []tea.Cmd

	batch = append(batch, m.panelManager.Init())

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
