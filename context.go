package main

import "github.com/charmbracelet/lipgloss"

type ProgramContext struct {
	Indent           bool
	Cursor           int
	JoinView         bool
	NormalStyle      Styles
	CursorStyle      Styles
	ChildCursorStyle Styles
}

type Styles struct {
	Gutter     lipgloss.Style
	NodeName   lipgloss.Style
	Everything lipgloss.Style
	Relation   lipgloss.Style
}

func InitProgramContext() ProgramContext {
	normal := NormalStyles()
	return ProgramContext{
		Cursor:           1,
		Indent:           true,
		JoinView:         true,
		NormalStyle:      normal,
		CursorStyle:      CursorStyle(normal),
		ChildCursorStyle: ChildCursorStyle(normal),
	}
}

func NormalStyles() Styles {
	gutterStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#777"))
	nodeNameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#65bcff"))
	everythingStyle := lipgloss.NewStyle()
	relationStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#c099ff"))

	return Styles{
		Gutter:     gutterStyle,
		NodeName:   nodeNameStyle,
		Everything: everythingStyle,
		Relation:   relationStyle,
	}
}

func CursorStyle(style Styles) Styles {
	background := lipgloss.Color("#2f334d")

	return Styles{
		Gutter:     style.Gutter.Foreground(lipgloss.Color("#ff336c")),
		NodeName:   style.NodeName.Background(background),
		Everything: style.Everything.Background(background),
		Relation:   style.Relation.Background(background),
	}
}

func ChildCursorStyle(style Styles) Styles {
	background := lipgloss.Color("#4f445e")

	return Styles{
		Gutter:     style.Gutter,
		NodeName:   style.NodeName.Background(background),
		Everything: style.Everything.Background(background),
		Relation:   style.Relation.Background(background),
	}
}
