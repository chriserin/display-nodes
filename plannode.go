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
}

func (node PlanNode) View(level int, lineNumber int, ctx ProgramContext) (string, int) {

	var background lipgloss.Color

	if ctx.Cursor == lineNumber {
		background = lipgloss.Color("#f33")
	} else {
		background = lipgloss.Color("#000")
	}

	levelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#777")).Background(background)
	nodeNameStyle := lipgloss.NewStyle().Bold(true).Background(background)
	everythingStyle := lipgloss.NewStyle().Background(background)
	newLineNumber := lineNumber + 1

	var buf strings.Builder

	buf.WriteString(levelStyle.Render(fmt.Sprintf("%d ", level)))
	buf.WriteString(levelStyle.Render(fmt.Sprintf("%d ", newLineNumber)))
	if ctx.Indent {
		buf.WriteString(everythingStyle.Render(strings.Repeat(" ", level-1)))
	}
	buf.WriteString(nodeNameStyle.Render(node.name() + " " + node.rows()))
	buf.WriteString("\n")

	var renderedNode string

	for _, childNode := range node.Plans {
		renderedNode, newLineNumber = childNode.View(level+1, newLineNumber, ctx)
		buf.WriteString(renderedNode)
	}

	return buf.String(), newLineNumber
}

func (node PlanNode) name() string {
	return strings.Trim(fmt.Sprintf("%s %s", node.PartialMode, node.NodeType), " ")
}

func (node PlanNode) rows() string {
	return fmt.Sprintf("(Rows planned=%d actual=%d)", node.PlanRows, node.ActualRows)
}
