package processes

import (
	"goputer/internal/components"
	"sort"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/v4/process"
)

type ProcessesModel struct {
	table  table.Model
	width  int
	height int
	card   components.Card
}

func MakeProcessesModel(width, height int) *ProcessesModel {
	card := components.NewCard("Processes", "")
	totalWidth := width - card.CardStyle.GetHorizontalBorderSize() - card.CardStyle.GetHorizontalPadding()

	t := table.New(
		table.WithColumns(getTableColumns(totalWidth)),
		table.WithFocused(true),
		table.WithWidth(totalWidth),
		table.WithHeight(9),
	)

	model := ProcessesModel{
		table:  t,
		width:  width,
		height: height,
		card:   card,
	}

	return &model
}

func (m *ProcessesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		targetWidth := m.width - m.card.CardStyle.GetHorizontalPadding() - m.card.CardStyle.GetHorizontalBorderSize()
		m.card = m.card.SetWidth(m.width).SetHeight(m.height)
		m.table.SetWidth(targetWidth)
		m.table.SetColumns(getTableColumns(targetWidth))
		return m, nil
	case ProcessMsg:
		var rows []table.Row
		for _, val := range []processInfo(msg) {
			rows = append(rows, table.Row{
				strconv.Itoa(int(val.PID)),
				val.Username,
				val.Name,
				strconv.FormatFloat(val.CPUPercent/10, 'f', -1, 32),
				strconv.FormatFloat(float64(val.MemPercent), 'f', -1, 32),
				val.Command,
			})
		}
		m.table.SetRows(rows)
		return m, checkProcesses()
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *ProcessesModel) View() string {
	card := m.card.SetContent(m.table.View())
	return card.Render()
}

func (m *ProcessesModel) Init() tea.Cmd {
	return tea.Batch(
		checkProcesses(),
	)
}

func (m *ProcessesModel) ToggleActive() {
	m.card = m.card.ToggleActive()
}

func getTableColumns(width int) []table.Column {
	autoWidth := (width / 6) - table.DefaultStyles().Cell.GetHorizontalPadding()
	columns := []table.Column{
		{Title: "PID", Width: autoWidth / 2},
		{Title: "User", Width: autoWidth / 2},
		{Title: "Process", Width: autoWidth},
		{Title: "CPU%", Width: autoWidth},
		{Title: "MEM%", Width: autoWidth},
		{Title: "Command", Width: autoWidth * 2},
	}

	return columns
}

type processInfo struct {
	PID        int32
	Name       string
	Username   string
	CPUPercent float64
	MemPercent float32
	Command    string
}

func checkProcesses() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		processes, _ := getTopProcesses()

		return ProcessMsg(processes)
	})
}

func getTopProcesses() ([]processInfo, error) {
	var processes []processInfo

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

		processes = append(processes, processInfo{
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

type ProcessMsg []processInfo
