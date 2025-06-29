package card

import (
	"github.com/charmbracelet/lipgloss"
)

type Card struct {
	header    string
	content   string
	CardStyle lipgloss.Style
	active    bool
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

func (c Card) SetContent(content string) Card {
	c.content = content
	return c
}

func (c Card) ToggleActive() Card {
	c.active = !c.active

	return c
}

func (c Card) Render() string {

	if c.active {
		c.CardStyle = c.CardStyle.BorderForeground(lipgloss.Color("10"))
	}
	w := c.CardStyle.GetWidth() - c.CardStyle.GetHorizontalBorderSize() // This is the actual width of the card

	// We remove the padding, because we have to take it into account
	// to calculate the actual length of the header content.
	header := headerStyle.Width(w - c.CardStyle.GetHorizontalPadding()).Render(c.header)

	return c.CardStyle.Width(w).Render(header + "\n" + c.content)
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
