package system

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/v4/mem"
)

type Memory struct {
	Ram  mem.VirtualMemoryStat
	Swap mem.SwapMemoryStat
}

func CheckMemory() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		ram, err := mem.VirtualMemory()
		if err != nil {
			return err
		}

		swap, err := mem.SwapMemory()
		if err != nil {
			return err
		}

		m := Memory{
			Ram:  *ram,
			Swap: *swap,
		}

		return MemoryMsg(m)
	})
}

type MemoryMsg Memory
