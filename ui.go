package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	topNode PlanNode
	ctx     ProgramContext
}

func runProgram(topNode PlanNode, ctx ProgramContext) {
	p := tea.NewProgram(Model{topNode: topNode, ctx: ctx})

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		case "a":
			return m, tea.EnterAltScreen
		case "A":
			return m, tea.ExitAltScreen
		case "i":
			m.ctx.Indent = false
			return m, nil
		case "I":
			m.ctx.Indent = true
			return m, nil
		default:
			return m, tea.Println(msg)
		}

	case tea.WindowSizeMsg:
		// Handle window size changes if needed
	}

	return m, nil
}

func (m Model) View() string {
	return m.topNode.View(1, m.ctx)
}
