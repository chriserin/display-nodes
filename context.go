package main

import "github.com/charmbracelet/lipgloss"

type ProgramContext struct {
	Indent           bool
	Cursor           int
	JoinView         bool
	DisplayRows      bool
	NormalStyle      Styles
	CursorStyle      Styles
	ChildCursorStyle Styles
	SelectedNode     PlanNode
}

type Styles struct {
	Gutter     lipgloss.Style
	NodeName   lipgloss.Style
	Everything lipgloss.Style
	Relation   lipgloss.Style
	Bracket    lipgloss.Style
	Value      lipgloss.Style
	Warning    lipgloss.Style
	Caution    lipgloss.Style
}

func InitProgramContext(selectedNode PlanNode) ProgramContext {
	normal := NormalStyles()

	return ProgramContext{
		Cursor:           0,
		Indent:           true,
		JoinView:         false,
		NormalStyle:      normal,
		CursorStyle:      CursorStyle(normal),
		ChildCursorStyle: ChildCursorStyle(normal),
		SelectedNode:     selectedNode,
	}
}

func NormalStyles() Styles {
	gutterStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#777"))
	nodeNameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#65bcff"))
	everythingStyle := lipgloss.NewStyle()
	relationStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#c099ff"))
	bracketStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#828bb8"))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#c3e88d"))
	warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#c53b53"))
	cautionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffc777"))

	return Styles{
		Gutter:     gutterStyle,
		NodeName:   nodeNameStyle,
		Everything: everythingStyle,
		Relation:   relationStyle,
		Bracket:    bracketStyle,
		Value:      valueStyle,
		Warning:    warningStyle,
		Caution:    cautionStyle,
	}
}

func CursorStyle(style Styles) Styles {
	background := lipgloss.Color("#222277")

	return Styles{
		Gutter:     style.Gutter.Foreground(lipgloss.Color("#ff966c")),
		NodeName:   style.NodeName.Background(background),
		Everything: style.Everything.Background(background),
		Relation:   style.Relation.Background(background),
		Bracket:    style.Bracket.Background(background),
		Value:      style.Value.Background(background),
		Warning:    style.Warning.Background(background),
		Caution:    style.Caution.Background(background),
	}
}

func ChildCursorStyle(style Styles) Styles {
	background := lipgloss.Color("#2f334d")

	return Styles{
		Gutter:     style.Gutter,
		NodeName:   style.NodeName.Background(background),
		Everything: style.Everything.Background(background),
		Relation:   style.Relation.Background(background),
		Bracket:    style.Bracket.Background(background),
		Value:      style.Value.Background(background),
		Warning:    style.Warning.Background(background),
		Caution:    style.Caution.Background(background),
	}
}
