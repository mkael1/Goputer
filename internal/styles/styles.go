package styles

import "github.com/charmbracelet/lipgloss"

var ContentTextStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#00ff00"))

var LabelStyle = lipgloss.NewStyle().Width(10).Align(lipgloss.Left)
