package main

import (
	"fmt"
	"goputer/internal/components/cpu"
	"goputer/internal/components/disk"
	"goputer/internal/components/footer"
	"goputer/internal/components/header"
	"goputer/internal/components/memory"
	"goputer/internal/components/processes"
	"goputer/internal/keys"
	"goputer/internal/panel"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type model struct {
	OS            string
	User          string
	width, height int
	header        *header.Header
	panelManager  *panel.PanelManager
	footer        *footer.Footer
}

func initialModel() model {
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
		width:        width,
		height:       height,
		header:       header.NewHeader(width),
		panelManager: panel.NewPanelManager(panel.WithCols(12), panel.WithCells(cells)),
		footer:       footer.NewFooter(width),
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
			m.footer.ShowHelp = !m.footer.ShowHelp
			m.calculateWindowDimensions(tea.WindowSizeMsg{Width: m.width, Height: m.height})
			return m, nil
		default:
			_, cmd := m.panelManager.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.calculateWindowDimensions(msg)
		return m, nil
	default:
		_, cmd := m.panelManager.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.header.View(),
		m.panelManager.View(),
		m.footer.View(),
	)
}

func (m model) calculateWindowDimensions(msg tea.WindowSizeMsg) {
	m.footer.Update(msg)
	m.header.Update(msg)

	panelMsg := msg
	panelMsg.Height = m.height - lipgloss.Height(m.footer.View()) - lipgloss.Height(m.header.View())
	m.panelManager.Update(panelMsg)
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
