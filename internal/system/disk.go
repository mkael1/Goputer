package system

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/v4/disk"
)

type Disk struct {
	UsageStats []disk.UsageStat
}

func CheckDisk() tea.Cmd {
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

		return DiskMsg(Disk{
			UsageStats: usageStats,
		})
	})
}

type DiskMsg Disk
