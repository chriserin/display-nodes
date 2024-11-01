package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	AltOn           key.Binding
	AltOff          key.Binding
	IndentToggle    key.Binding
	Up              key.Binding
	Down            key.Binding
	Help            key.Binding
	Quit            key.Binding
	JoinView        key.Binding
	ToggleRows      key.Binding
	ToggleBuffers   key.Binding
	ToggleCost      key.Binding
	ToggleTimes     key.Binding
	NextStatDisplay key.Binding
	PrevStatDisplay key.Binding
	ToggleParallel  key.Binding
	ReExecute       key.Binding
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
		{k.Up, k.Down, k.AltOn, k.AltOff, k.IndentToggle, k.ToggleRows, k.ToggleBuffers, k.ToggleCost, k.ToggleTimes, k.NextStatDisplay, k.PrevStatDisplay, k.ToggleParallel}, // first column
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
	ToggleCost: key.NewBinding(
		key.WithKeys("C"),
		key.WithHelp("C", "Toggle Costs"),
	),
	ToggleTimes: key.NewBinding(
		key.WithKeys("T"),
		key.WithHelp("T", "Toggle Times"),
	),
	NextStatDisplay: key.NewBinding(
		key.WithKeys("]"),
		key.WithHelp("]", "Next Stat Display"),
	),
	PrevStatDisplay: key.NewBinding(
		key.WithKeys("["),
		key.WithHelp("[", "Previous Stat Display"),
	),
	ToggleParallel: key.NewBinding(
		key.WithKeys("P"),
		key.WithHelp("P", "Toggle Parallel"),
	),
	ReExecute: key.NewBinding(
		key.WithKeys("X"),
		key.WithHelp("X", "ReExecute Query"),
	),
}

type Model struct {
	keys            keyMap
	help            help.Model
	nodes           []PlanNode
	ctx             ProgramContext
	DisplayNodes    []PlanNode
	StatusLine      StatusLine
	detailsViewport viewport.Model
	source          Source
	queryRun        QueryRun
	spinner         spinner.Model
	sqlChannel      chan QueryRun
	loading         bool
}

func InitModel(explainPlan ExplainPlan, source Source, queryRun QueryRun) Model {
	return Model{
		nodes:           explainPlan.nodes,
		ctx:             InitProgramContext(explainPlan.nodes[0], explainPlan.analyzed),
		keys:            keys,
		help:            help.New(),
		DisplayNodes:    explainPlan.nodes,
		StatusLine:      NewStatusLine(explainPlan),
		detailsViewport: NewViewPort(),
		source:          source,
		queryRun:        queryRun,
		spinner:         initialSpinner(),
	}
}

func initialSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return s
}

func (m *Model) UpdateModel(explainPlan ExplainPlan) {
	m.nodes = explainPlan.nodes
	m.DisplayNodes = explainPlan.nodes
	m.StatusLine = NewStatusLine(explainPlan)
}

func NewViewPort() viewport.Model {
	vp := viewport.New(80, 10)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2).PaddingLeft(2)
	return vp
}

type Source struct {
	sourceType SourceType
	fileName   string
}

type SourceType int

const (
	SOURCE_STDIN SourceType = iota
	SOURCE_FILE
)

func RunProgram(explainPlan ExplainPlan, source Source, queryRun QueryRun) {
	program := tea.NewProgram(
		InitModel(explainPlan, source, queryRun),
	)

	if _, err := program.Run(); err != nil {
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
			if len(m.DisplayNodes) > 0 {
				m.ctx.SelectedNode = m.DisplayNodes[m.ctx.Cursor]
			} else {
				m.ctx.SelectedNode = PlanNode{}
			}
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
		case key.Matches(msg, m.keys.ToggleCost):
			if m.ctx.StatDisplay == DisplayCost {
				m.ctx.StatDisplay = DisplayNothing
			} else {
				m.ctx.StatDisplay = DisplayCost
			}
		case key.Matches(msg, m.keys.ToggleTimes):
			if m.ctx.StatDisplay == DisplayTime {
				m.ctx.StatDisplay = DisplayNothing
			} else {
				m.ctx.StatDisplay = DisplayTime
			}
		case key.Matches(msg, m.keys.NextStatDisplay):
			m.ctx.StatDisplay = nextStatDisplay(m.ctx)
		case key.Matches(msg, m.keys.PrevStatDisplay):
			m.ctx.StatDisplay = prevStatDisplay(m.ctx)
		case key.Matches(msg, m.keys.ToggleParallel):
			m.ctx.DisplayParallel = !m.ctx.DisplayParallel
		case key.Matches(msg, m.keys.ReExecute):
			m.sqlChannel = make(chan QueryRun)
			m.loading = true
			go func() {
				queryRun := NewQueryRun(m.queryRun.originalFilename)
				queryWithExplain := queryRun.WithExplainAnalyze()
				result := ExecuteExplain(queryWithExplain)
				queryRun.SetResult(result)
				pgexDir := CreatePgexDir()
				queryRun.WritePgexFile(pgexDir)
				m.sqlChannel <- queryRun
			}()
			return m, m.spinner.Tick
		default:
			return m, tea.Println(msg)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		select {
		case queryRun := <-m.sqlChannel:
			m.queryRun = queryRun
			explainPlan := Convert(queryRun.result)
			m.UpdateModel(explainPlan)
			m.loading = false
			close(m.sqlChannel)
		default:
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.ctx.Width = msg.Width
	}

	return m, nil
}

func prevStatDisplay(ctx ProgramContext) StatView {
	newStatDisplay := ctx.StatDisplay

	for true {
		if newStatDisplay == 0 {
			newStatDisplay = DisplayCost
		} else {
			newStatDisplay = (newStatDisplay - 1) % 5
		}

		if ctx.Analyzed {
			break
		} else if slices.Contains([]StatView{DisplayRows, DisplayCost, DisplayNothing}, newStatDisplay) {
			break
		}
	}
	return newStatDisplay
}

func nextStatDisplay(ctx ProgramContext) StatView {
	newStatDisplay := ctx.StatDisplay

	for true {
		newStatDisplay = (newStatDisplay + 1) % 5

		if ctx.Analyzed {
			break
		} else if slices.Contains([]StatView{DisplayRows, DisplayCost, DisplayNothing}, newStatDisplay) {
			break
		}
	}
	return newStatDisplay
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

	buf.WriteString(m.StatusLine.View(m.ctx))
	if m.loading {
		buf.WriteString(m.spinner.View())
	}
	buf.WriteString(HeadersView(m.ctx, m.source))

	for i, node := range m.DisplayNodes {
		buf.WriteString(node.View(i, m.ctx))
	}

	m.detailsViewport.Width = m.ctx.Width - 3
	m.detailsViewport.SetContent(m.ctx.SelectedNode.Content(m.ctx))

	buf.WriteString("\n")
	buf.WriteString(m.detailsViewport.View())
	buf.WriteString("\n")
	buf.WriteString(m.help.View(m.keys))

	return buf.String()
}

func HeadersView(ctx ProgramContext, source Source) string {
	var sourceOutput string

	if source.sourceType == SOURCE_FILE {
		sourceOutput = fmt.Sprintf("   %s", source.fileName)
	} else {
		sourceOutput = "   STDIN"
	}

	var headers string
	if ctx.StatDisplay == DisplayTime {
		headers = fmt.Sprintf("%15s%15s ", "Startup", "Total")
	} else if ctx.StatDisplay == DisplayCost {
		headers = fmt.Sprintf("%15s%15s ", "Startup", "Total")
	} else if ctx.StatDisplay == DisplayBuffers {
		headers = fmt.Sprintf("%15s%15s ", "Total", "Read")
	} else if ctx.StatDisplay == DisplayRows {
		headers = fmt.Sprintf("%15s%15s ", "Planned", "Actual")
	} else if ctx.StatDisplay == DisplayNothing {
		headers = ""
	}
	needed := ctx.Width - len(sourceOutput)
	return fmt.Sprintf("%s%*s\n", sourceOutput, needed, headers)
}
