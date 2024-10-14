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
	nodeNameStyle := lipgloss.NewStyle().Bold(true)
	everythingStyle := lipgloss.NewStyle()

	return Styles{
		Gutter:     gutterStyle,
		NodeName:   nodeNameStyle,
		Everything: everythingStyle,
	}
}

func CursorStyle(style Styles) Styles {
	background := lipgloss.Color("#2f334d")

	return Styles{
		Gutter:     style.Gutter.Foreground(lipgloss.Color("#ff336c")),
		NodeName:   style.NodeName.Background(background),
		Everything: style.Everything.Background(background),
	}
}

func ChildCursorStyle(style Styles) Styles {
	background := lipgloss.Color("#4f445e")

	return Styles{
		Gutter:     style.Gutter,
		NodeName:   style.NodeName.Background(background),
		Everything: style.Everything.Background(background),
	}
}
