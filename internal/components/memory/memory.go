package memory

import (
	"fmt"
	"goputer/internal/components"
	"goputer/internal/storage"
	"goputer/internal/styles"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v4/mem"
)

type MemoryModel struct {
	memory Memory
	width  int
	height int
	card   components.Card
}

func MakeMemoryModel(width, height int) *MemoryModel {
	card := components.NewCard("Memory Usage", "")
	model := MemoryModel{
		width:  width,
		height: height,
		card:   card,
	}
	return &model
}

func (m *MemoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case MemoryMsg:
		m.memory = Memory(msg)
		return m, checkMemory()
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		log.Printf("memory height %d", m.height)
		m.card = m.card.SetWidth(m.width).SetHeight(m.height)
		return m, nil
	}
	var cmd tea.Cmd
	return m, cmd
}

func (m *MemoryModel) View() string {
	totalRamGb := storage.BytesToGb(m.memory.ram.Total)
	usedRamGb := storage.BytesToGb(m.memory.ram.Used)
	freeRamGb := storage.BytesToGb(m.memory.ram.Free)
	cachedRamGb := storage.BytesToGb(m.memory.ram.Cached)

	swapUsedGb := storage.BytesToGb(m.memory.swap.Used)
	swapTotalGb := storage.BytesToGb(m.memory.swap.Total)
	swapUsed := fmt.Sprintf("%.1f GB / %.1f GB", swapUsedGb, swapTotalGb)

	labelWidth := 12                       // adjust as needed
	valueWidth := m.width - labelWidth - 4 // account for padding

	prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))

	valueStyle := styles.ContentTextStyle.Width(valueWidth).AlignHorizontal(lipgloss.Right)

	ramRow := lipgloss.JoinHorizontal(lipgloss.Top,
		styles.LabelStyle.Render("RAM:"),
		prog.ViewAs(m.memory.ram.UsedPercent/100))
	totalRow := lipgloss.JoinHorizontal(lipgloss.Top,
		styles.LabelStyle.Render("Total:"),
		valueStyle.Render(fmt.Sprintf("%.1fGB", totalRamGb)))
	usedRow := lipgloss.JoinHorizontal(lipgloss.Top,
		styles.LabelStyle.Render("Used:"),
		valueStyle.Render(fmt.Sprintf("%.1fGB", usedRamGb)))
	freeRow := lipgloss.JoinHorizontal(lipgloss.Top,
		styles.LabelStyle.Render("Free:"),
		valueStyle.Render(fmt.Sprintf("%.1fGB", freeRamGb)))
	cachedRow := lipgloss.JoinHorizontal(lipgloss.Top,
		styles.LabelStyle.Render("Cached:"),
		valueStyle.Render(fmt.Sprintf("%.1fGB\n", cachedRamGb)))
	swapRow := lipgloss.JoinHorizontal(lipgloss.Top,
		styles.LabelStyle.Render("Swap:"),
		prog.ViewAs(m.memory.swap.UsedPercent/100))
	swapUsageRow := lipgloss.JoinHorizontal(lipgloss.Top,
		styles.LabelStyle.Render("Swap Used:"),
		valueStyle.Render(swapUsed))

	content := lipgloss.JoinVertical(lipgloss.Left, ramRow, totalRow, usedRow, freeRow, cachedRow, swapRow, swapUsageRow)
	return m.card.SetContent(content).Render()
}

func (m *MemoryModel) Init() tea.Cmd {
	return tea.Batch(
		checkMemory(),
	)
}

func (m *MemoryModel) ToggleActive() {
	m.card = m.card.ToggleActive()
}

type Memory struct {
	ram  mem.VirtualMemoryStat
	swap mem.SwapMemoryStat
}

func checkMemory() tea.Cmd {
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
			ram:  *ram,
			swap: *swap,
		}

		return MemoryMsg(m)
	})
}

type MemoryMsg Memory
