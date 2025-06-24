package components

import (
	"goputer/internal/card"
	"goputer/internal/system"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func RenderProcessesCard(processes []system.ProcessInfo, cardWidth int) string {
	columns := []table.Column{
		{Title: "PID", Width: 4},
		{Title: "User", Width: 10},
		{Title: "Process", Width: 10},
		{Title: "CPU%", Width: 10},
		{Title: "MEM%", Width: 10},
		{Title: "Command", Width: cardWidth - 10 - 10 - 10 - 10 - 4 - 14}, // TODO: calculate the width better...
	}

	rows := []table.Row{}

	for _, process := range processes {
		r := table.Row{
			strconv.Itoa(int(process.PID)),
			process.Username,
			process.Name,
			strconv.FormatFloat(process.CPUPercent/10, 'f', -1, 32),
			strconv.FormatFloat(float64(process.MemPercent), 'f', -1, 32),
			process.Command,
		}
		rows = append(rows, r)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
		// Width(20)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return card.New("Top Processes", t.View()).SetWidth(cardWidth).Render()
}
