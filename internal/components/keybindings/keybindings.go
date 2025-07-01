package keybindings

import (
	"fmt"
	"goputer/internal/components"
	"goputer/internal/keys"
	"goputer/internal/styles"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeybindingsModel struct {
	width  int
	height int
	card   components.Card
}

func MakeKeybindingsModel(width, height int) *KeybindingsModel {
	card := components.NewCard("Keybindings", "").ShowHeader(false)

	model := KeybindingsModel{
		width:  width,
		height: height,
		card:   card,
	}
	return &model
}

func (m *KeybindingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	var cmd tea.Cmd
	return m, cmd
}

func (m *KeybindingsModel) View() string {
	var bindings []key.Binding
	bindings = append(bindings, keys.KeyMapToSlice(keys.Global)...)
	bindings = append(bindings, keys.KeyMapToSlice(keys.Panes)...)

	content := ""

	for _, val := range bindings {
		content += fmt.Sprintf("%s%s\n", styles.LabelStyle.Render(val.Help().Key), descStyle.Render(val.Help().Desc))
	}

	header := "? help"

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerStyle.Width(m.width).Render(header),
		contentStyle.Render(content),
	)
}

func (m *KeybindingsModel) Init() tea.Cmd {
	return nil
}

var headerStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Background(lipgloss.Color("249")).
	Foreground(lipgloss.Color("235"))

var contentStyle = lipgloss.NewStyle().
	Padding(0, 1)

var descStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("252"))
