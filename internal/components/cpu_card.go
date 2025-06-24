package components

import (
	"fmt"
	"goputer/internal/card"
	"goputer/internal/system"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

func RenderCpuCard(cpu system.Cpu, cardWidth int) string {

	// Calculate average CPU usage
	percentage := 0.0
	if len(cpu.Usage) > 0 {
		var sum float64
		for _, val := range cpu.Usage {
			sum += val
		}
		percentage = (sum / float64(len(cpu.Usage))) / 100
	}
	overallBar := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C")).ViewAs(percentage)
	overallRow := lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Overall:"), overallBar)
	content := lipgloss.JoinVertical(lipgloss.Left,
		fmt.Sprintf("CPU: %d cores (%d threads)", cpu.Cores, cpu.Threads),
		overallRow,
	)

	// Show the first 4 threads instead of everything
	// TODO: add the ability to view more
	for index, val := range cpu.Usage {
		if index > 3 {
			continue
		}
		bar := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C")).ViewAs(val / 100)
		text := labelStyle.Render(fmt.Sprintf("Thread %d: ", index+1))
		row := lipgloss.JoinHorizontal(lipgloss.Left, text, bar)
		content = lipgloss.JoinVertical(lipgloss.Left, content, row)
	}

	content = lipgloss.JoinVertical(lipgloss.Left, content, "", getUptimeString(cpu.Uptime))

	return card.New("CPU Usage", content).SetWidth(cardWidth).Render()
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
