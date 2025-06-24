package system

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/v4/cpu"
)

type Cpu struct {
	Info cpu.InfoStat
}

func CheckCpu() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		info, err := cpu.Info()
		if err != nil {
			panic(err)
		}
		cpu := Cpu{
			Info: info[0],
		}

		return CpuMsg(cpu)
	})
}

type CpuMsg Cpu
