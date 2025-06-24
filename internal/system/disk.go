package system

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/v4/disk"
)

type Disk struct {
	UsageStat disk.UsageStat
}

func CheckDisk() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		disk, err := disk.Usage("/")
		if err != nil {

		}

		return DiskMsg(Disk{
			UsageStat: *disk,
		})
	})
}

type DiskMsg Disk
