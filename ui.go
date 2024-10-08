package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	AltOn     key.Binding
	AltOff    key.Binding
	IndentOn  key.Binding
	IndentOff key.Binding
	Up        key.Binding
	Down      key.Binding
	Help      key.Binding
	Quit      key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.AltOn, k.AltOff, k.IndentOn, k.IndentOff}, // first column
		{k.Help, k.Quit}, // second column
	}
}

var keys = keyMap{
	AltOn: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "alt screen on"),
	),
	AltOff: key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "alt screen off"),
	),
	IndentOn: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "indent on"),
	),
	IndentOff: key.NewBinding(
		key.WithKeys("I"),
		key.WithHelp("I", "indent off"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type Model struct {
	keys    keyMap
	help    help.Model
	topNode PlanNode
	ctx     ProgramContext
}

func runProgram(topNode PlanNode, ctx ProgramContext) {
	p := tea.NewProgram(Model{topNode: topNode, ctx: ctx, keys: keys, help: help.New()})

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
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.AltOn):
			return m, tea.EnterAltScreen
		case key.Matches(msg, m.keys.AltOff):
			return m, tea.ExitAltScreen
		case key.Matches(msg, m.keys.IndentOff):
			m.ctx.Indent = false
			return m, nil
		case key.Matches(msg, m.keys.IndentOn):
			m.ctx.Indent = true
			return m, nil
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		default:
			return m, tea.Println(msg)
		}

	case tea.WindowSizeMsg:
		// Handle window size changes if needed
	}

	return m, nil
}

func (m Model) View() string {
	return m.topNode.View(1, m.ctx) + "\n" + m.help.View(m.keys)
}
