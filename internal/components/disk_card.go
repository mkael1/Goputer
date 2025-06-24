package components

import (
	"fmt"
	"goputer/internal/card"
	"goputer/internal/system"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

func RenderDiskCard(disk system.Disk, cardWidth int) string {

	overallBar := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C")).ViewAs(disk.UsageStat.UsedPercent / 100)
	rootBarRow := lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render(disk.UsageStat.Path), overallBar)
	rootDisk := lipgloss.JoinVertical(
		lipgloss.Left,
		rootBarRow,
		fmt.Sprintf("Size: %.1f GB / %.1f GB", bytesToGB(disk.UsageStat.Free), bytesToGB(disk.UsageStat.Total)),
	)

	return card.New("Disk Usage", rootDisk).SetWidth(cardWidth).Render()
}
