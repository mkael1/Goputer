package panel

import (
	"goputer/internal/keys"
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DisplayType int

const (
	DisplayGrid DisplayType = iota
)

type PanelManager struct {
	Grid
	selectedPanelIndex int
	width, height      int
	Type               DisplayType
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
			p.Cells[p.selectedPanelIndex].ToggleActive()
			if p.selectedPanelIndex == len(p.Cells)-1 {
				p.selectedPanelIndex = 0
			} else {
				p.selectedPanelIndex += 1
			}
			p.Cells[p.selectedPanelIndex].ToggleActive()
		case key.Matches(msg, keys.Panes.Left):
			p.Cells[p.selectedPanelIndex].ToggleActive()
			if p.selectedPanelIndex == 0 {
				p.selectedPanelIndex = len(p.Cells) - 1
			} else {
				p.selectedPanelIndex -= 1
			}
			p.Cells[p.selectedPanelIndex].ToggleActive()
		}
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		p.calulateCellsWidth(msg)
		return p, nil
	}

	if _, isKeyMsg := msg.(tea.KeyMsg); !isKeyMsg {
		for _, val := range p.Cells {
			_, cmd := val.Update(msg)
			cmds = append(cmds, cmd)
		}
	} else {
		// Only send the key messages to the active pane
		_, cmd := p.Cells[p.selectedPanelIndex].Update(msg)
		cmds = append(cmds, cmd)
	}

	return p, tea.Batch(cmds...)
}

func (p *PanelManager) calulateCellsWidth(msg tea.WindowSizeMsg) []tea.Cmd {
	var cmds []tea.Cmd

	widthPerCols := msg.Width / p.Cols
	for _, c := range p.Cells {
		msg.Width = widthPerCols * c.ColSpan
		_, cmd := c.Update(msg)
		cmds = append(cmds, cmd)
	}

	return cmds
}

func (p *PanelManager) View() string {
	var content string

	colsLeft := p.Cols
	var rows []string
	rows = append(rows, "")
	rowIndex := 0
	for _, c := range p.Cells {
		colsLeft -= c.ColSpan
		if colsLeft >= 0 {
			rows[rowIndex] = lipgloss.JoinHorizontal(lipgloss.Top, rows[rowIndex], c.Panel.View())
		} else {
			log.Println(colsLeft)
			colsLeft = p.Cols - c.ColSpan
			rows = append(rows, c.Panel.View())
		}
	}

	for _, row := range rows {
		content = lipgloss.JoinVertical(lipgloss.Left, content, row)
	}

	return content
}

func (p *PanelManager) Init() tea.Cmd {
	var batch []tea.Cmd

	for _, c := range p.Cells {
		batch = append(batch, c.Init())
	}

	return tea.Batch(
		batch...,
	)
}

type PanelOption func(*PanelManager)

func NewPanelManager(width, height int, opts ...PanelOption) *PanelManager {
	p := &PanelManager{
		width:              width,
		height:             height,
		selectedPanelIndex: 0,
		Type:               DisplayGrid,
		Grid: Grid{
			Cols:        1,
			Cells:       []Cell{},
			AutoSpacing: false,
		},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func WithCols(count int) PanelOption {
	return func(pm *PanelManager) {
		pm.Cols = count
	}
}

func WithCells(cc []Cell) PanelOption {
	return func(pm *PanelManager) {
		pm.Cells = cc
	}
}

// Calculates the available spacing between each cell in the grid so there's no leftover at the end.
func WithAutomaticSpacing() PanelOption {
	return func(pm *PanelManager) {
		pm.AutoSpacing = true
	}
}

type Grid struct {
	Cols        int
	Cells       []Cell
	AutoSpacing bool
}

type Cell struct {
	Panel
	ColSpan int
}
