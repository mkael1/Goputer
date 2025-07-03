package components

import (
	"github.com/charmbracelet/lipgloss"
)

type Card struct {
	header     string
	content    string
	CardStyle  lipgloss.Style
	active     bool
	showHeader bool
}

func NewCard(header string, content string) Card {
	return Card{
		header:     header,
		content:    content,
		CardStyle:  cardStyle,
		showHeader: true,
	}
}

func (c Card) SetWidth(i int) Card {
	c.CardStyle = c.CardStyle.Width(i)
	return c
}

func (c Card) SetHeight(i int) Card {
	c.CardStyle = c.CardStyle.Height(i)
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

func (c Card) ShowHeader(b bool) Card {
	c.showHeader = b

	return c
}

func (c Card) Render() string {

	if c.active {
		c.CardStyle = c.CardStyle.BorderForeground(lipgloss.Color("10"))
	}
	w := c.CardStyle.GetWidth() - c.CardStyle.GetHorizontalBorderSize() // This is the actual width of the card
	h := c.CardStyle.GetHeight() - c.CardStyle.GetVerticalBorderSize()  // This is the actual height of the card

	// We remove the padding, because we have to take it into account
	// to calculate the actual length of the header content.
	if c.showHeader {
		header := headerStyle.Width(w-c.CardStyle.GetHorizontalPadding()).Render(c.header) + "\n"
		c.content = header + c.content
	}

	return c.CardStyle.Width(w).Height(h).Render(c.content)
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
