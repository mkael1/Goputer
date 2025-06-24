package main

import (
	"fmt"
	"goputer/internal/components"
	"goputer/internal/system"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type model struct {
	OS     string
	Time   time.Time
	User   string
	Memory system.Memory
	Cpu    system.Cpu
	Disk   system.Disk
}

func initialModel() model {
	user, _ := user.Current()
	return model{
		OS:   runtime.GOOS,
		Time: time.Now(),
		User: user.Username,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case system.MemoryMsg:
		m.Memory = system.Memory(msg)
		return m, system.CheckMemory()

	case system.CpuMsg:
		m.Cpu = system.Cpu(msg)
		return m, system.CheckCpu()

	case system.DiskMsg:
		m.Disk = system.Disk(msg)
		return m, system.CheckDisk()
	}

	return m, nil
}

func (m model) View() string {
	rightSideHeader := strings.ToUpper(m.OS) + " | " + m.Time.Format(time.Kitchen)

	// Get actual terminal width
	termWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	if termWidth == 0 {
		termWidth = 80 // fallback width
	}
	var headerStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ffffff")).
		BorderBottom(true).
		Width(termWidth).
		MarginBottom(1)

	s := fmt.Sprintf("System Monitor Resources - %s", m.User)
	// Create the full header with left and right content
	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		s,
		strings.Repeat(" ", termWidth-lipgloss.Width(s)-lipgloss.Width(rightSideHeader)),
		rightSideHeader,
	)

	header = headerStyle.Render(header)
	return header + lipgloss.JoinHorizontal(
		lipgloss.Top,
		components.RenderCpuCard(m.Cpu, getTermWidth()/2-2),
		components.RenderMemoryCard(m.Memory, getTermWidth()/2-2),
	) + lipgloss.JoinHorizontal(
		lipgloss.Top,
		components.RenderDiskCard(m.Disk, getTermWidth()/2-2),
		"",
	)

}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		system.CheckMemory(),
		system.CheckCpu(),
		system.CheckDisk(),
	)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Ahhh... there's been an error: %v", err)
		os.Exit(1)
	}
}

func getTermWidth() int {
	termWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	return termWidth
}
