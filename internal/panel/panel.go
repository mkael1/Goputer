package panel

import tea "github.com/charmbracelet/bubbletea"

type Panel interface {
	tea.Model
	ToggleActive()
}
