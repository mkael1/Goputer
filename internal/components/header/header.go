package header

import (
	"fmt"
	"os/user"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Header struct {
	width int
	os    string
	user  string
}

func NewHeader(width int) *Header {
	user, _ := user.Current()

	return &Header{
		width: width,
		os:    runtime.GOOS,
		user:  user.Username,
	}
}

func (h *Header) Init() tea.Cmd {
	return nil
}

func (h *Header) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.width = msg.Width
		return h, nil
	}
	var cmd tea.Cmd
	return h, cmd
}

func (h *Header) View() string {

	var headerStyle = lipgloss.NewStyle().
		Width(h.width)

	leftHeader := "System Monitor Resources"
	rightHeader := fmt.Sprintf("%s | %s | %s", h.user, strings.ToUpper(h.os), time.Now().Format(time.Kitchen))

	// Create the full header with left and right content
	header := headerStyle.Render(lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftHeader,
		strings.Repeat(" ", h.width-lipgloss.Width(leftHeader)-lipgloss.Width(rightHeader)-headerStyle.GetHorizontalPadding()),
		rightHeader,
	))

	return header
}
