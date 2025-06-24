package main

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v4/mem"
	"golang.org/x/term"
)

type model struct {
	OS     string
	Time   time.Time
	User   string
	Memory Memory
}

type Memory struct {
	Ram  mem.VirtualMemoryStat
	Swap mem.SwapMemoryStat
}

func initialModel() model {
	user, _ := user.Current()
	// v, _ := mem.VirtualMemory()
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

	case memoryMsg:
		m.Memory = Memory(msg)
		return m, checkMemory()
	}

	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("System Monitor Resources - %s", m.User)
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

	// Create the full header with left and right content
	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		s,
		strings.Repeat(" ", termWidth-lipgloss.Width(s)-lipgloss.Width(rightSideHeader)),
		rightSideHeader,
	)

	header = headerStyle.Render(header)

	return header + m.getRamCard()
}

func (m model) Init() tea.Cmd {
	return checkMemory()
}

func checkMemory() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		ram, err := mem.VirtualMemory()
		swap, err := mem.SwapMemory()
		if err != nil {
			return err
		}

		m := Memory{
			Ram:  *ram,
			Swap: *swap,
		}

		return memoryMsg(m)
	})
}

type memoryMsg Memory

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

var (
	cardStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			MarginTop(1)

	greenStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff00"))
)

func (m model) getRamCard() string {
	totalRamGb := m.Memory.Ram.Total >> 30
	usedRamGb := m.Memory.Ram.Used >> 30
	freeRamGb := m.Memory.Ram.Free >> 30
	cachedRamGb := m.Memory.Ram.Cached >> 30

	swapUsedGb := m.Memory.Swap.Used >> 30
	swapTotalGb := m.Memory.Swap.Used >> 30
	swapUsed := fmt.Sprintf("%v GB / %v GB", swapUsedGb, swapTotalGb)

	prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
	return cardStyle.Render(fmt.Sprintf(
		"RAM: %s\nTotal: %s\nUsed: %s\nFree: %s\nCached: %s\n\nSwap: %s\nSwap Used: %s",
		prog.ViewAs(m.Memory.Ram.UsedPercent/100),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Render(fmt.Sprintf("%dGB", totalRamGb)),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Render(fmt.Sprintf("%dGB", usedRamGb)),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Render(fmt.Sprintf("%dGB", freeRamGb)),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Render(fmt.Sprintf("%dGB", cachedRamGb)),
		prog.ViewAs(m.Memory.Swap.UsedPercent/100),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Render(swapUsed),
	))
}
