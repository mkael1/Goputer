package disk

import (
	"fmt"
	"goputer/internal/components"
	"goputer/internal/storage"
	"goputer/internal/styles"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v4/disk"
)

type DiskModel struct {
	disks  []disk.UsageStat
	width  int
	height int
	card   components.Card
}

func MakeDiskModel(width, height int) *DiskModel {
	card := components.NewCard("Disk Usage", "")
	model := DiskModel{
		width:  width,
		height: height,
		card:   card,
	}
	return &model
}

func (m *DiskModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case DiskMsg:
		m.disks = []disk.UsageStat(msg)
		return m, checkDisk()
	case tea.WindowSizeMsg:
		m.width = msg.Width / 2
		m.height = msg.Height
		m.card = m.card.SetWidth(m.width)
		return m, nil
	}
	var cmd tea.Cmd
	return m, cmd
}

func (m *DiskModel) View() string {
	var content string
	for _, usageStat := range m.disks {
		prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"), progress.WithWidth(20)).ViewAs(usageStat.UsedPercent / 100)
		diskBlock := lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Left, styles.LabelStyle.Render(usageStat.Path), prog),
			styles.LabelStyle.Render("Size:")+fmt.Sprintf("%.1f GB / %.1f GB", storage.BytesToGb(usageStat.Free), storage.BytesToGb(usageStat.Total)),
		)

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			diskBlock,
		)
	}

	return m.card.SetContent(content).Render()
}

func (m *DiskModel) Init() tea.Cmd {
	return tea.Batch(
		checkDisk(),
	)
}

func (m *DiskModel) ToggleActive() {
	m.card = m.card.ToggleActive()
}

func checkDisk() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		partitions, _ := disk.Partitions(false)
		var usageStats []disk.UsageStat
		for _, partition := range partitions {
			diskStat, err := disk.Usage(partition.Mountpoint)
			if err != nil {
				panic(err)
			}
			usageStats = append(usageStats, *diskStat)
		}

		return DiskMsg(usageStats)
	})
}

type DiskMsg []disk.UsageStat
