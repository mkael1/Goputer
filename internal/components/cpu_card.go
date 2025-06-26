package components

import (
	"fmt"
	"goputer/internal/card"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
)

type CPUModel struct {
	cpu    cpuData
	width  int
	height int
}

func MakeCpuModel(width, height int) *CPUModel {
	model := CPUModel{
		width:  width,
		height: height,
	}
	return &model
}

func (m *CPUModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case CpuMsg:
		m.cpu = cpuData(msg)
		return m, checkCpu()
	case tea.WindowSizeMsg:
		m.width = msg.Width / 2
		m.height = msg.Height
		return m, nil
	}
	var cmd tea.Cmd
	return m, cmd
}

func (m *CPUModel) View() string {

	// Calculate average CPU usage
	percentage := 0.0
	if len(m.cpu.Usage) > 0 {
		var sum float64
		for _, val := range m.cpu.Usage {
			sum += val
		}
		percentage = (sum / float64(len(m.cpu.Usage))) / 100
	}
	overallBar := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C")).ViewAs(percentage)
	overallRow := lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Overall:"), overallBar)
	content := lipgloss.JoinVertical(lipgloss.Left,
		fmt.Sprintf("CPU: %d cores (%d threads)", m.cpu.Cores, m.cpu.Threads),
		overallRow,
	)

	// Show the first 4 threads instead of everything
	// TODO: add the ability to view more
	for index, val := range m.cpu.Usage {
		if index > 3 {
			continue
		}
		bar := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C")).ViewAs(val / 100)
		text := labelStyle.Render(fmt.Sprintf("Thread %d: ", index+1))
		row := lipgloss.JoinHorizontal(lipgloss.Left, text, bar)
		content = lipgloss.JoinVertical(lipgloss.Left, content, row)
	}

	content = lipgloss.JoinVertical(lipgloss.Left, content, "", getUptimeString(m.cpu.Uptime))

	return card.New("CPU Usage", content).SetWidth(m.width).Render()
}

func (m *CPUModel) Init() tea.Cmd {
	return tea.Batch(
		checkCpu(),
	)
}

func getUptimeString(uptime uint64) string {
	days := uptime / (24 * 3600)
	hours := (uptime % (24 * 3600)) / 3600
	minutes := (uptime % 3600) / 60
	seconds := uptime % 60

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		labelStyle.Render("Uptime:"),
		fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds),
	)
}

var labelStyle = lipgloss.NewStyle().Width(10).Align(lipgloss.Left)

type cpuData struct {
	Info    cpu.InfoStat
	Cores   int
	Threads int
	Usage   []float64
	Uptime  uint64
}

func checkCpu() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		info, err := cpu.Info()
		if err != nil {
			panic(err)
		}
		cpuUsage, err := cpu.Percent(0, true)
		if err != nil {
			panic(err)
		}

		threads, _ := cpu.Counts(true)
		cores, _ := cpu.Counts(false)
		uptime, _ := host.Uptime()

		cpu := cpuData{
			Info:    info[0],
			Usage:   cpuUsage,
			Threads: threads,
			Cores:   cores,
			Uptime:  uptime,
		}

		return CpuMsg(cpu)
	})
}

type CpuMsg cpuData
