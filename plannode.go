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

func (node PlanNode) View(level int, ctx ProgramContext) string {
	levelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#777"))
	nodeNameStyle := lipgloss.NewStyle().Bold(true)

	var buf strings.Builder

	buf.WriteString(levelStyle.Render(fmt.Sprintf("%d ", level)))
	if ctx.Indent {
		buf.WriteString(strings.Repeat(" ", level-1))
	}
	buf.WriteString(nodeNameStyle.Render(node.name()))
	buf.WriteString(" ")
	buf.WriteString(node.rows())
	buf.WriteString("\n")

	for _, childNode := range node.Plans {
		buf.WriteString(childNode.View(level+1, ctx))
	}

	return buf.String()
}

func (node PlanNode) name() string {
	return strings.Trim(fmt.Sprintf("%s %s", node.PartialMode, node.NodeType), " ")
}

func (node PlanNode) rows() string {
	return fmt.Sprintf("(Rows planned=%d actual=%d)", node.PlanRows, node.ActualRows)
}
