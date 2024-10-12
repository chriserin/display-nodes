package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type PlanNode struct {
	NodeType         string
	Plans            []PlanNode
	PlanRows         int
	ActualRows       int
	PartialMode      string
	Position         position
	JoinViewPosition position
	RelationName     string
}

func (node PlanNode) View(ctx ProgramContext) string {

	var background lipgloss.Color

	var viewPosition position

	if ctx.JoinView {
		viewPosition = node.JoinViewPosition
	} else {
		viewPosition = node.Position
	}

	if ctx.Cursor == viewPosition.LineNumber {
		background = lipgloss.Color("#f33")
	} else if ctx.Cursor == viewPosition.Parent {
		background = lipgloss.Color("#a33")
	} else {
		background = lipgloss.Color("#000")
	}

	levelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#777")).Background(background)
	nodeNameStyle := lipgloss.NewStyle().Bold(true).Background(background)
	everythingStyle := lipgloss.NewStyle().Background(background)

	var buf strings.Builder

	buf.WriteString(levelStyle.Render(fmt.Sprintf("%2d ", viewPosition.LineNumber)))

	if ctx.Indent {
		buf.WriteString(everythingStyle.Render(strings.Repeat("  ", viewPosition.Level-1)))
	}

	buf.WriteString(nodeNameStyle.Render(node.name() + " " + node.rows()))
	buf.WriteString("\n")

	return buf.String()
}

func (node PlanNode) Display(ctx ProgramContext) bool {
	if ctx.JoinView {
		return node.JoinViewPosition.Display
	} else {
		return node.Position.Display
	}
}

func (node PlanNode) name() string {
	return strings.Trim(fmt.Sprintf("%s %s", node.PartialMode, node.NodeType), " ")
}

func (node PlanNode) rows() string {
	return fmt.Sprintf("(Rows planned=%d actual=%d)", node.PlanRows, node.ActualRows)
}
