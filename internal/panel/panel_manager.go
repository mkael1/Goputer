package panel

import (
	"goputer/internal/keys"

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
	Type               DisplayType
	Width              int
	Height             int
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
		p.Width = msg.Width
		p.Height = msg.Height
		p.calculateCells(msg)
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

func (p *PanelManager) calculateCells(msg tea.WindowSizeMsg) []tea.Cmd {
	var cmds []tea.Cmd
	originalMsg := msg
	rows := p.calculateGridRows()

	targetCellHeight := originalMsg.Height / len(rows)
	leftoverHeightPx := originalMsg.Height % len(rows)
	for index, row := range rows {
		msg.Height = targetCellHeight

		if leftoverHeightPx > 0 {
			msg.Height += 1
			leftoverHeightPx -= 1
		}

		leftoverWidthPx := originalMsg.Width % len(rows[index])
		for _, c := range row {
			msg.Width = originalMsg.Width / (p.Cols / c.ColSpan)

			// This it to accomodate leftover pixels if it's an odd width for the terminal
			if leftoverWidthPx > 0 {
				msg.Width += 1
				leftoverWidthPx -= 1
			}

			_, cmd := c.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return cmds
}

func (p *PanelManager) View() string {
	var finalContent string

	rows := p.calculateGridRows()

	for i, row := range rows {
		rowContent := ""
		for _, cell := range row {
			rowContent = lipgloss.JoinHorizontal(lipgloss.Top, rowContent, cell.View())
		}
		if i > 0 {
			finalContent = lipgloss.JoinVertical(lipgloss.Left, finalContent, rowContent)
		} else {
			finalContent = rowContent
		}
	}

	return finalContent
}

func (p *PanelManager) calculateGridRows() [][]Cell {
	colsLeft := p.Cols
	var rows [][]Cell
	rowIndex := 0
	for _, c := range p.Cells {
		colsLeft -= c.ColSpan
		if colsLeft >= 0 {
			for len(rows) <= rowIndex {
				rows = append(rows, []Cell{})
			}

			rows[rowIndex] = append(rows[rowIndex], c)
		} else {
			colsLeft = p.Cols - c.ColSpan
			rowIndex++
			for len(rows) <= rowIndex {
				rows = append(rows, []Cell{})
			}
			rows[rowIndex] = append(rows[rowIndex], c)
		}
	}
	return rows
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

func NewPanelManager(opts ...PanelOption) *PanelManager {
	p := &PanelManager{
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
		pm.Cells[pm.selectedPanelIndex].ToggleActive()
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
