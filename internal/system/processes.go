package system

import (
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/v4/process"
)

type ProcessInfo struct {
	PID        int32
	Name       string
	Username   string
	CPUPercent float64
	MemPercent float32
	Command    string
}

func CheckProcesses() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		processes, _ := getTopProcesses()

		return ProcessMsg(processes)
	})
}

func getTopProcesses() ([]ProcessInfo, error) {
	var processes []ProcessInfo

	// Get all processes
	pids, err := process.Pids()
	if err != nil {
		return nil, err
	}

	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err != nil {
			continue // Process might have died
		}

		name, _ := p.Name()
		username, _ := p.Username()
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()
		cmdline, _ := p.Cmdline()

		processes = append(processes, ProcessInfo{
			PID:        pid,
			Name:       name,
			Username:   username,
			CPUPercent: cpuPercent,
			MemPercent: memPercent,
			Command:    cmdline,
		})
	}

	sort.Slice(processes, func(i, j int) bool {
		return processes[i].CPUPercent > processes[j].CPUPercent
	})

	return processes, nil
}

type ProcessMsg []ProcessInfo
