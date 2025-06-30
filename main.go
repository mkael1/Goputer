package main

import (
	"fmt"
	"goputer/internal/components"
	"goputer/internal/keys"
	"goputer/internal/panel"
	"os"
	"os/user"
	"runtime"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

type model struct {
	OS            string
	User          string
	width, height int
	panelManager  *panel.PanelManager
}

func initialModel() model {
	user, _ := user.Current()
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	var panels []panel.Panel
	panels = append(
		panels,
		components.MakeCpuModel(width/2, height),
		components.MakeMemoryModel(width/2, height),
		components.MakeDiskModel(width/2, height),
		components.MakeProcessesModel(width/2, height),
		components.MakeKeybindingsModel(width, height),
	)

	model := model{
		OS:           runtime.GOOS,
		User:         user.Username,
		width:        width,
		height:       height,
		panelManager: panel.NewPanelManager(width, height),
	}

	for _, p := range panels {
		model.panelManager.AddPanel(p)
	}

	return model
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Global.Quit):
			return m, tea.Quit
		default:
			return m.panelManager.Update(msg)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		_, cmd := m.panelManager.Update(msg)
		cmds = append(cmds, cmd)
	}

	if _, isKeyMsg := msg.(tea.KeyMsg); !isKeyMsg {
		_, cmd := m.panelManager.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	// rightSideHeader := strings.ToUpper(m.OS) + " | " + time.Now().Format(time.Kitchen)

	// var headerStyle = lipgloss.NewStyle().
	// 	BorderStyle(lipgloss.RoundedBorder()).
	// 	BorderForeground(lipgloss.Color("#ffffff")).
	// 	BorderBottom(true).
	// 	Width(m.width).
	// 	Height(1)

	// headerString := fmt.Sprintf("System Monitor Resources - %s", m.User)

	// // Create the full header with left and right content
	// header := lipgloss.JoinHorizontal(
	// 	lipgloss.Top,
	// 	headerString,
	// 	strings.Repeat(" ", m.width-lipgloss.Width(headerString)-lipgloss.Width(rightSideHeader)),
	// 	rightSideHeader,
	// )

	// header = headerStyle.Render(header)
	// footer := headerStyle.BorderTop(true).BorderBottom(false).Render("? toggle help * q quit")

	return m.panelManager.View()
}

func (m model) Init() tea.Cmd {
	var batch []tea.Cmd

	batch = append(batch, m.panelManager.Init())

	return tea.Batch(
		batch...,
	)
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Ahhh... there's been an error: %v", err)
		os.Exit(1)
	}
}
