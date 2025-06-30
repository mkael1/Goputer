package components

import (
	"fmt"
	"goputer/internal/card"
	"goputer/internal/keys"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeybindingsModel struct {
	width  int
	height int
	card   card.Card
}

func MakeKeybindingsModel(width, height int) *KeybindingsModel {
	card := card.New("Keybindings", "").ShowHeader(false)
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
		m.card = m.card.SetWidth(m.width)
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
		content += fmt.Sprintf("%s - %s\n", val.Help().Key, val.Help().Desc)
	}

	return m.card.SetContent(content).Render()
}

func (m *KeybindingsModel) Init() tea.Cmd {
	return nil
}

func (m *KeybindingsModel) ToggleActive() {
	m.card = m.card.ToggleActive()
}
