package system

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
)

type Cpu struct {
	Info    cpu.InfoStat
	Cores   int
	Threads int
	Usage   []float64
	Uptime  uint64
}

func CheckCpu() tea.Cmd {
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

		cpu := Cpu{
			Info:    info[0],
			Usage:   cpuUsage,
			Threads: threads,
			Cores:   cores,
			Uptime:  uptime,
		}

		return CpuMsg(cpu)
	})
}

type CpuMsg Cpu
