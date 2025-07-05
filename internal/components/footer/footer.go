package footer

import (
	"fmt"
	"goputer/internal/components"
	"goputer/internal/keys"
	"goputer/internal/styles"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Footer struct {
	width    int
	card     components.Card
	ShowHelp bool
}

func NewFooter(width int) *Footer {
	card := components.NewCard("Keybindings", "").ShowHeader(false)

	return &Footer{
		ShowHelp: false,
		card:     card,
		width:    width,
	}
}

func (m *Footer) Init() tea.Cmd {
	return nil
}

func (m *Footer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	}
	var cmd tea.Cmd
	return m, cmd
}

func (m *Footer) View() string {
	var footer string
	if !m.ShowHelp {
		footer = lipgloss.NewStyle().BorderTop(true).BorderBottom(false).Render("? toggle help • q quit")
	} else {
		var bindings []key.Binding
		bindings = append(bindings, keys.KeyMapToSlice(keys.Global)...)
		bindings = append(bindings, keys.KeyMapToSlice(keys.Panes)...)

		content := ""

		for _, val := range bindings {
			content += fmt.Sprintf("%s%s\n", styles.LabelStyle.Render(val.Help().Key), descStyle.Render(val.Help().Desc))
		}

		header := "? help"

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			headerStyle.Width(m.width).Render(header),
			contentStyle.Render(content),
		)
		footer = m.card.SetWidth(m.width).SetContent(content).Render()
	}

	return footer
}

var headerStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Background(lipgloss.Color("249")).
	Foreground(lipgloss.Color("235"))

var contentStyle = lipgloss.NewStyle().
	Padding(0, 1)

var descStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("252"))
