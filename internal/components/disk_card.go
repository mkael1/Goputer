package components

import (
	"fmt"
	"goputer/internal/card"
	"goputer/internal/system"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

func RenderDiskCard(disk system.Disk, cardWidth int) string {
	var content string
	for _, usageStat := range disk.UsageStats {
		prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"), progress.WithWidth(20)).ViewAs(usageStat.UsedPercent / 100)
		diskBlock := lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render(usageStat.Path), prog),
			labelStyle.Render("Size:")+fmt.Sprintf("%.1f GB / %.1f GB", bytesToGB(usageStat.Free), bytesToGB(usageStat.Total)),
		)

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			diskBlock,
		)
	}

	return card.New("Disk Usage", content).SetWidth(cardWidth).Render()
}
