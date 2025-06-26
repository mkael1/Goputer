package card

import "github.com/charmbracelet/lipgloss"

func New(header string, content string) card {
	return card{
		header:    header,
		content:   content,
		cardStyle: cardStyle,
	}
}

func (c card) SetWidth(i int) card {
	c.cardStyle = c.cardStyle.Width(i)
	return c
}

func (c card) Render() string {
	// First, render the content to determine the card's inner width
	tempCard := c.cardStyle.Render(c.header + "\n" + c.content)
	cardWidth := lipgloss.Width(tempCard) - cardStyle.GetHorizontalPadding() - cardStyle.GetHorizontalBorderSize()

	// Now render the header with the calculated width
	header := headerStyle.Width(cardWidth).Render(c.header)

	return cardStyle.Render(header + "\n" + c.content)
}

type card struct {
	header    string
	content   string
	cardStyle lipgloss.Style
}

var (
	cardStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			MarginTop(1).
			Padding(1)
	headerStyle = lipgloss.NewStyle().
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			Foreground(lipgloss.Color("#00ffff")).
			Bold(true)
)
