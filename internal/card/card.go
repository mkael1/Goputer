package card

import "github.com/charmbracelet/lipgloss"

type Card struct {
	header    string
	content   string
	CardStyle lipgloss.Style
}

func New(header string, content string) Card {
	return Card{
		header:    header,
		content:   content,
		CardStyle: cardStyle,
	}
}

func (c Card) SetWidth(i int) Card {
	c.CardStyle = c.CardStyle.Width(i)
	return c
}

func (c Card) Render() string {
	w := c.CardStyle.GetWidth() - c.CardStyle.GetHorizontalBorderSize() - c.CardStyle.GetHorizontalPadding()
	header := headerStyle.Width(w).Render(c.header)

	return cardStyle.Render(header + "\n" + c.content)
}

var (
	cardStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1).
			MarginTop(1)
	headerStyle = lipgloss.NewStyle().
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			Foreground(lipgloss.Color("#00ffff")).
			Bold(true)
)
