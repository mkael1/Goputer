package panel

import (
	"goputer/internal/keys"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PanelManager struct {
	panels             []Panel
	selectedPanelIndex int
	width, height      int
}

type Panel interface {
	tea.Model
	ToggleActive()
}

func (p *PanelManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Panes.Right):
			p.panels[p.selectedPanelIndex].ToggleActive()
			if p.selectedPanelIndex == len(p.panels)-1 {
				p.selectedPanelIndex = 0
			} else {
				p.selectedPanelIndex += 1
			}
			p.panels[p.selectedPanelIndex].ToggleActive()
		case key.Matches(msg, keys.Panes.Left):
			p.panels[p.selectedPanelIndex].ToggleActive()
			if p.selectedPanelIndex == 0 {
				p.selectedPanelIndex = len(p.panels) - 1
			} else {
				p.selectedPanelIndex -= 1
			}
			p.panels[p.selectedPanelIndex].ToggleActive()
		}
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		for _, pane := range p.panels {
			_, cmd := pane.Update(msg)
			cmds = append(cmds, cmd)
		}
		return p, nil
	}

	if _, isKeyMsg := msg.(tea.KeyMsg); !isKeyMsg {
		for _, val := range p.panels {
			_, cmd := val.Update(msg)
			cmds = append(cmds, cmd)
		}
	} else {
		// Only send the key messages to the active pane
		_, cmd := p.panels[p.selectedPanelIndex].Update(msg)
		cmds = append(cmds, cmd)
	}

	return p, tea.Batch(cmds...)
}

func (p *PanelManager) View() string {
	var panelsView string
	if len(p.panels) >= 4 {
		topRow := lipgloss.JoinHorizontal(
			lipgloss.Top,
			p.panels[0].View(),
			p.panels[1].View(),
		)

		bottomRow := lipgloss.JoinHorizontal(
			lipgloss.Top,
			p.panels[2].View(),
			p.panels[3].View(),
		)

		panelsView = lipgloss.JoinVertical(
			lipgloss.Left,
			topRow,
			bottomRow,
		)
	}

	return panelsView
}

func (p *PanelManager) AddPanel(panel Panel) {
	p.panels = append(p.panels, panel)
	// In the scenario where we're adding a pane that is equivalent to the current active pane index,
	// We want to set it as active automatically.
	if len(p.panels)-1 == p.selectedPanelIndex {
		p.panels[p.selectedPanelIndex].ToggleActive()
	}
}

func (p *PanelManager) Init() tea.Cmd {
	var batch []tea.Cmd

	for _, pane := range p.panels {
		batch = append(batch, pane.Init())
	}

	return tea.Batch(
		batch...,
	)
}

func NewPanelManager(width, height int) *PanelManager {
	return &PanelManager{
		width:              width,
		height:             height,
		selectedPanelIndex: 0,
	}
}
