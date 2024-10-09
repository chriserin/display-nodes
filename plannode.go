package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type PlanNode struct {
	NodeType    string
	Plans       []PlanNode
	PlanRows    int
	ActualRows  int
	PartialMode string
	LineNumber  int
	Level       int
}

func (node PlanNode) View(ctx ProgramContext) string {

	var background lipgloss.Color

	if ctx.Cursor == node.LineNumber {
		background = lipgloss.Color("#f33")
	} else {
		background = lipgloss.Color("#000")
	}

	levelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#777")).Background(background)
	nodeNameStyle := lipgloss.NewStyle().Bold(true).Background(background)
	everythingStyle := lipgloss.NewStyle().Background(background)

	var buf strings.Builder

	buf.WriteString(levelStyle.Render(fmt.Sprintf("%d ", node.Level)))
	buf.WriteString(levelStyle.Render(fmt.Sprintf("%d ", node.LineNumber)))

	if ctx.Indent {
		buf.WriteString(everythingStyle.Render(strings.Repeat(" ", node.Level-1)))
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
