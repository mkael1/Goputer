package components

import (
	"fmt"
	"goputer/internal/card"
	"goputer/internal/styles"
	"goputer/internal/system"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

func RenderMemoryCard(memory system.Memory, cardWidth int) string {
	totalRamGb := bytesToGB(memory.Ram.Total)
	usedRamGb := bytesToGB(memory.Ram.Used)
	freeRamGb := bytesToGB(memory.Ram.Free)
	cachedRamGb := bytesToGB(memory.Ram.Cached)

	swapUsedGb := bytesToGB(memory.Swap.Used)
	swapTotalGb := bytesToGB(memory.Swap.Total)
	swapUsed := fmt.Sprintf("%.1f GB / %.1f GB", swapUsedGb, swapTotalGb)

	labelWidth := 12                         // adjust as needed
	valueWidth := cardWidth - labelWidth - 4 // account for padding

	prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))

	labelStyle := lipgloss.NewStyle().Width(labelWidth).AlignHorizontal(lipgloss.Left)
	valueStyle := styles.ContentTextStyle.Width(valueWidth).AlignHorizontal(lipgloss.Right)

	ramRow := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("RAM:"),
		prog.ViewAs(memory.Ram.UsedPercent/100))
	totalRow := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("Total:"),
		valueStyle.Render(fmt.Sprintf("%.1fGB", totalRamGb)))
	usedRow := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("Used:"),
		valueStyle.Render(fmt.Sprintf("%.1fGB", usedRamGb)))
	freeRow := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("Free:"),
		valueStyle.Render(fmt.Sprintf("%.1fGB", freeRamGb)))
	cachedRow := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("Cached:"),
		valueStyle.Render(fmt.Sprintf("%.1fGB\n", cachedRamGb)))
	swapRow := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("Swap:"),
		prog.ViewAs(memory.Swap.UsedPercent/100))
	swapUsageRow := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("Swap Used:"),
		valueStyle.Render(swapUsed))

	content := lipgloss.JoinVertical(lipgloss.Left, ramRow, totalRow, usedRow, freeRow, cachedRow, swapRow, swapUsageRow)
	return card.New("Memory Usage", content).SetWidth(cardWidth).Render()
}

func bytesToGB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
}
