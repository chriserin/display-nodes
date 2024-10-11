package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type PlanNode struct {
	NodeType     string
	Plans        []PlanNode
	PlanRows     int
	ActualRows   int
	PartialMode  string
	Position     position
	RelationName string
}

func (node PlanNode) View(ctx ProgramContext) string {

	var background lipgloss.Color

	if ctx.Cursor == node.Position.LineNumber {
		background = lipgloss.Color("#f33")
	} else if ctx.Cursor == node.Position.Parent {
		background = lipgloss.Color("#a33")
	} else {
		background = lipgloss.Color("#000")
	}

	levelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#777")).Background(background)
	nodeNameStyle := lipgloss.NewStyle().Bold(true).Background(background)
	everythingStyle := lipgloss.NewStyle().Background(background)

	var buf strings.Builder

	buf.WriteString(levelStyle.Render(fmt.Sprintf("%2d ", node.Position.LineNumber)))
	buf.WriteString(levelStyle.Render(fmt.Sprintf("%2d ", node.Position.Level)))

	if ctx.Indent {
		buf.WriteString(everythingStyle.Render(strings.Repeat("  ", node.Position.Level-1)))
	}

	buf.WriteString(nodeNameStyle.Render(node.name() + " " + node.rows()))
	buf.WriteString("\n")

	return buf.String()
}

func (node PlanNode) name() string {
	return strings.Trim(fmt.Sprintf("%s %s", node.PartialMode, node.NodeType), " ")
}

func (node PlanNode) rows() string {
	return fmt.Sprintf("(Rows planned=%d actual=%d)", node.PlanRows, node.ActualRows)
}
