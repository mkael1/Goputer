package components

import (
	"fmt"
	"goputer/internal/card"
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
}

func MakeMemoryModel(width, height int) *MemoryModel {
	model := MemoryModel{
		width:  width,
		height: height,
	}
	return &model
}

func (m *MemoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case MemoryMsg:
		m.memory = Memory(msg)
		return m, checkMemory()
	case tea.WindowSizeMsg:
		m.width = msg.Width / 2
		m.height = msg.Height
		return m, nil
	}
	var cmd tea.Cmd
	return m, cmd
}

func (m *MemoryModel) View() string {
	totalRamGb := bytesToGB(m.memory.ram.Total)
	usedRamGb := bytesToGB(m.memory.ram.Used)
	freeRamGb := bytesToGB(m.memory.ram.Free)
	cachedRamGb := bytesToGB(m.memory.ram.Cached)

	swapUsedGb := bytesToGB(m.memory.swap.Used)
	swapTotalGb := bytesToGB(m.memory.swap.Total)
	swapUsed := fmt.Sprintf("%.1f GB / %.1f GB", swapUsedGb, swapTotalGb)

	labelWidth := 12                       // adjust as needed
	valueWidth := m.width - labelWidth - 4 // account for padding

	prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))

	labelStyle := lipgloss.NewStyle().Width(labelWidth).AlignHorizontal(lipgloss.Left)
	valueStyle := styles.ContentTextStyle.Width(valueWidth).AlignHorizontal(lipgloss.Right)

	ramRow := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("RAM:"),
		prog.ViewAs(m.memory.ram.UsedPercent/100))
	totalRow := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("Total:"),
		valueStyle.Render(fmt.Sprintf("%.1fGB", totalRamGb)))
	usedRow := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("Used:"),
		valueStyle.Render(fmt.Sprintf("%.1fGB", usedRamGb)))
	freeRow := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("Free:"),
		valueStyle.Render(fmt.Sprintf("%.1fGB", freeRamGb)))
	cachedRow := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("Cached:"),
		valueStyle.Render(fmt.Sprintf("%.1fGB\n", cachedRamGb)))
	swapRow := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("Swap:"),
		prog.ViewAs(m.memory.swap.UsedPercent/100))
	swapUsageRow := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("Swap Used:"),
		valueStyle.Render(swapUsed))

	content := lipgloss.JoinVertical(lipgloss.Left, ramRow, totalRow, usedRow, freeRow, cachedRow, swapRow, swapUsageRow)
	log.Println(m.width)
	return card.New("Memory Usage", content).SetWidth(m.width).Render()
}

func (m *MemoryModel) Init() tea.Cmd {
	return tea.Batch(
		checkMemory(),
	)
}

func bytesToGB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
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
