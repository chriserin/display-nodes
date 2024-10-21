package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	AltOn          key.Binding
	AltOff         key.Binding
	IndentToggle   key.Binding
	Up             key.Binding
	Down           key.Binding
	Help           key.Binding
	Quit           key.Binding
	JoinView       key.Binding
	ToggleRows     key.Binding
	ToggleBuffers  key.Binding
	ToggleParallel key.Binding
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
		{k.Up, k.Down, k.AltOn, k.AltOff, k.IndentToggle, k.ToggleRows, k.ToggleBuffers, k.ToggleParallel}, // first column
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
	IndentToggle: key.NewBinding(
		key.WithKeys("I"),
		key.WithHelp("I", "indent toggle"),
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
	JoinView: key.NewBinding(
		key.WithKeys("J"),
		key.WithHelp("J", "Join"),
	),
	ToggleRows: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "Toggle Rows"),
	),
	ToggleBuffers: key.NewBinding(
		key.WithKeys("B"),
		key.WithHelp("B", "Toggle Buffers"),
	),
	ToggleParallel: key.NewBinding(
		key.WithKeys("P"),
		key.WithHelp("P", "Toggle Parallel"),
	),
}

type Model struct {
	keys         keyMap
	help         help.Model
	nodes        []PlanNode
	ctx          ProgramContext
	DisplayNodes []PlanNode
	StatusLine   StatusLine
}

func runProgram(nodes []PlanNode, executionTime float64, ctx ProgramContext) {
	p := tea.NewProgram(
		Model{
			nodes:        nodes,
			ctx:          ctx,
			keys:         keys,
			help:         help.New(),
			DisplayNodes: nodes,
			StatusLine: StatusLine{
				ExecutionTime: executionTime,
				TotalBuffers:  nodes[0].SharedBuffersHit + nodes[0].SharedBuffersRead,
				TotalRows:     nodes[0].ActualRows,
			},
		})

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
		case key.Matches(msg, m.keys.IndentToggle):
			m.ctx.Indent = !m.ctx.Indent
		case key.Matches(msg, m.keys.Up):
			if m.ctx.Cursor-1 >= 0 {
				m.ctx.Cursor = m.ctx.Cursor - 1
				m.ctx.SelectedNode = m.DisplayNodes[m.ctx.Cursor]
			}
		case key.Matches(msg, m.keys.Down):
			if m.ctx.Cursor+1 < len(m.DisplayNodes) {
				m.ctx.Cursor = m.ctx.Cursor + 1
				m.ctx.SelectedNode = m.DisplayNodes[m.ctx.Cursor]
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.JoinView):
			m.ctx.JoinView = !m.ctx.JoinView
			m.DisplayNodes = displayedNodes(m.nodes, m.ctx)
			m.ctx.Cursor = 0
			m.ctx.SelectedNode = m.DisplayNodes[m.ctx.Cursor]
		case key.Matches(msg, m.keys.ToggleRows):
			if m.ctx.StatDisplay == DisplayRows {
				m.ctx.StatDisplay = DisplayNothing
			} else {
				m.ctx.StatDisplay = DisplayRows
			}
		case key.Matches(msg, m.keys.ToggleBuffers):
			if m.ctx.StatDisplay == DisplayBuffers {
				m.ctx.StatDisplay = DisplayNothing
			} else {
				m.ctx.StatDisplay = DisplayBuffers
			}
		case key.Matches(msg, m.keys.ToggleParallel):
			m.ctx.DisplayParallel = !m.ctx.DisplayParallel
		default:
			return m, tea.Println(msg)
		}

	case tea.WindowSizeMsg:
		m.ctx.Width = msg.Width
	}

	return m, nil
}

func displayedNodes(nodes []PlanNode, ctx ProgramContext) []PlanNode {
	resultNodes := make([]PlanNode, 0, 1)

	for _, node := range nodes {
		if node.Display(ctx) {
			resultNodes = append(resultNodes, node)
		}
	}

	return resultNodes
}

func (m Model) View() string {
	var buf strings.Builder

	buf.WriteString(m.StatusLine.View(m.ctx.Width))

	for i, node := range m.DisplayNodes {
		buf.WriteString(node.View(i, m.ctx))
	}

	return buf.String() + "\n" + m.help.View(m.keys)
}
