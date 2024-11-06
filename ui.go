package main

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	IndentToggle     key.Binding
	Up               key.Binding
	Down             key.Binding
	Help             key.Binding
	Quit             key.Binding
	JoinView         key.Binding
	ToggleRows       key.Binding
	ToggleBuffers    key.Binding
	ToggleCost       key.Binding
	ToggleTimes      key.Binding
	ToggleDisplaySql key.Binding
	NextStatDisplay  key.Binding
	PrevStatDisplay  key.Binding
	ToggleParallel   key.Binding
	ReExecute        key.Binding
	PrevQueryRun     key.Binding
	NextQueryRun     key.Binding
	SqlUp            key.Binding
	SqlDown          key.Binding
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
		{k.Up, k.Down, k.IndentToggle, k.ToggleRows, k.ToggleBuffers, k.ToggleCost, k.ToggleTimes, k.NextStatDisplay, k.PrevStatDisplay, k.ToggleParallel, k.ToggleDisplaySql}, // first column
		{k.Help, k.Quit}, // second column
	}
}

var keys = keyMap{
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
	ToggleDisplaySql: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "Toggle Display SQL"),
	),
	ReExecute: key.NewBinding(
		key.WithKeys("X"),
		key.WithHelp("X", "ReExecute Query"),
	),
	PrevQueryRun: key.NewBinding(
		key.WithKeys("{"),
		key.WithHelp("{", "Prev Query Run"),
	),
	NextQueryRun: key.NewBinding(
		key.WithKeys("}"),
		key.WithHelp("}", "Next Query Run"),
	),
	SqlUp: key.NewBinding(
		key.WithKeys("ctrl+p"),
		key.WithHelp("ctrl+p", "SQL Up"),
	),
	SqlDown: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "SQL Down"),
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
	sqlViewport     viewport.Model
	source          Source
	queryRun        QueryRun
	spinner         spinner.Model
	sqlChannel      chan QueryRun
	loading         bool
	pgexPointer     string
}

func InitModel(source Source) Model {
	return Model{
		ctx:             InitProgramContext(),
		keys:            keys,
		help:            help.New(),
		detailsViewport: NewViewPort(),
		sqlViewport:     NewViewPort(),
		source:          source,
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
	m.SetDisplayNodes(displayedNodes(explainPlan.nodes, m.ctx))
	m.StatusLine = NewStatusLine(explainPlan)
}

func (m *Model) SetDisplayNodes(nodes []PlanNode) {
	m.DisplayNodes = displayedNodes(nodes, m.ctx)
	m.setSqlViewHeight()
}

func (m *Model) setSqlViewHeight() {
	m.sqlViewport.Height = m.ctx.Height - len(m.DisplayNodes) - 6
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
	input      string
}

func (s Source) DisplayName() string {
	_, file := path.Split(s.fileName)
	return file
}

func (s Source) View(ctx ProgramContext) string {
	if s.sourceType == SOURCE_FILE {
		return ctx.StatusStyles.AltNormal.Render(fmt.Sprintf("FILE - %s", s.DisplayName()))
	} else {
		return "STDIN"
	}
}

type SourceType int

const (
	SOURCE_STDIN SourceType = iota
	SOURCE_FILE
)

func RunProgram(source Source) {
	model := InitModel(source)

	if source.sourceType == SOURCE_STDIN {
		explainPlan := Convert(source.input)
		model.UpdateModel(explainPlan)
		model.ctx.ResetContext(explainPlan)
	}

	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)

	if _, err := program.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

type previousQueryRunMsg struct{}

func PreviousQueryRunCmd() tea.Msg {
	return previousQueryRunMsg{}
}

type nextQueryRunMsg struct{}

func NextQueryRunCmd() tea.Msg {
	return nextQueryRunMsg{}
}

type executeQueryMsg struct{}

func ExecuteQueryCmd() tea.Msg {
	return executeQueryMsg{}
}

type executeExplainQueryMsg struct{}

func ExecuteExplainQueryCmd() tea.Msg {
	return executeExplainQueryMsg{}
}

func (m Model) Init() tea.Cmd {
	if m.source.sourceType == SOURCE_STDIN {
		return nil
	} else {
		return ExecuteExplainQueryCmd
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
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
			m.SetDisplayNodes(displayedNodes(m.nodes, m.ctx))
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
		case key.Matches(msg, m.keys.ToggleDisplaySql):
			m.ctx.DisplaySql = !m.ctx.DisplaySql
		case key.Matches(msg, m.keys.ReExecute):
			return m, ExecuteQueryCmd
		case key.Matches(msg, m.keys.PrevQueryRun):
			return m, PreviousQueryRunCmd
		case key.Matches(msg, m.keys.NextQueryRun):
			return m, NextQueryRunCmd
		case key.Matches(msg, m.keys.SqlUp):
			m.sqlViewport.LineUp(1)
		case key.Matches(msg, m.keys.SqlDown):
			m.sqlViewport.LineDown(1)
		default:
			return m, tea.Println(msg)
		}
	case executeExplainQueryMsg:
		m.loading = true
		m.sqlChannel = make(chan QueryRun)
		go func() {
			queryRun := NewQueryRun(m.source.fileName)
			queryWithExplain := queryRun.WithExplain()
			result := ExecuteExplain(queryWithExplain)
			queryRun.SetResult(result)
			m.sqlChannel <- queryRun
		}()
		return m, m.spinner.Tick
	case executeQueryMsg:
		m.loading = true
		m.sqlChannel = make(chan QueryRun)
		go func() {
			queryRun := NewQueryRun(m.source.fileName)
			queryWithExplain := queryRun.WithExplainAnalyze()
			result := ExecuteExplain(queryWithExplain)
			queryRun.SetResult(result)
			pgexDir := CreatePgexDir()
			queryRun.WritePgexFile(pgexDir)
			m.sqlChannel <- queryRun
		}()
		return m, m.spinner.Tick
	case previousQueryRunMsg:
		newQueryRun := m.queryRun.previousQueryRun()
		if newQueryRun != m.queryRun {
			UpdateModel(&m, newQueryRun)
		}
		return m, nil
	case nextQueryRunMsg:
		newQueryRun := m.queryRun.nextQueryRun()
		if newQueryRun != m.queryRun {
			UpdateModel(&m, newQueryRun)
		}
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		select {
		case queryRun := <-m.sqlChannel:
			UpdateModel(&m, queryRun)
			m.loading = false
			close(m.sqlChannel)
			if !m.ctx.Analyzed {
				return m, ExecuteQueryCmd
			}
		default:
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.ctx.Width = msg.Width
		m.ctx.Height = msg.Height
		m.setSqlViewHeight()
	}

	return m, nil
}

func UpdateModel(m *Model, queryRun QueryRun) {
	m.queryRun = queryRun
	explainPlan := Convert(queryRun.result)
	m.UpdateModel(explainPlan)
	m.ctx.ResetContext(explainPlan)
	m.sqlViewport.SetContent(ansi.Wordwrap(queryRun.query, m.ctx.Width-6, ""))
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

	var spinnerView string
	if m.loading {
		spinnerView = m.spinner.View()
	} else {
		spinnerView = "  "
	}
	buf.WriteString(spinnerView)
	sourceView := m.source.View(m.ctx)
	buf.WriteString(sourceView)

	spaceAvailable := m.ctx.Width - ansi.StringWidth(sourceView)

	buf.WriteString(fmt.Sprintf("%*s%*s\n", spaceAvailable-10, m.ctx.StatDisplay.String(), 10, ""))

	statusLine := m.StatusLine.View(m.ctx)
	buf.WriteString(statusLine)
	buf.WriteString(HeadersView(m.ctx, m.ctx.Width-ansi.StringWidth(statusLine)))
	buf.WriteString("\n")

	for i, node := range m.DisplayNodes {
		buf.WriteString(node.View(i, m.ctx))
	}

	buf.WriteString("\n")

	if !m.help.ShowAll {
		if m.ctx.DisplaySql {
			buf.WriteString(m.sqlViewport.View())
		} else {
			m.detailsViewport.Width = m.ctx.Width - 3
			if m.ctx.SelectedNode.NodeType == "" {
				m.detailsViewport.Height = 3
			} else {
				m.detailsViewport.Height = 10
			}
			m.detailsViewport.SetContent(m.ctx.SelectedNode.Content(m.ctx))
			buf.WriteString(m.detailsViewport.View())
		}
	}

	buf.WriteString("\n")
	buf.WriteString(m.help.View(m.keys))

	return buf.String()
}

func HeadersView(ctx ProgramContext, spaceAvailable int) string {
	var headers string
	if ctx.StatDisplay == DisplayTime {
		headers = fmt.Sprintf("%10s%15s ", "Startup", "Total")
	} else if ctx.StatDisplay == DisplayCost {
		headers = fmt.Sprintf("%10s%15s ", "Startup", "Total")
	} else if ctx.StatDisplay == DisplayBuffers {
		headers = fmt.Sprintf("%10s%15s ", "Total", "Read")
	} else if ctx.StatDisplay == DisplayRows {
		headers = fmt.Sprintf("%10s%15s ", "Planned", "Actual")
	} else if ctx.StatDisplay == DisplayNothing {
		headers = ""
	}
	return fmt.Sprintf("%*s", spaceAvailable, headers)
}
