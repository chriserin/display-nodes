package main

import (
	"fmt"
	"strings"
)

type PlanNode struct {
	NodeType    string
	Plans       []PlanNode
	PlanRows    int
	ActualRows  int
	PartialMode string
}

func (node PlanNode) View(level int) string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("%d ", level))
	buf.WriteString(node.name())
	buf.WriteString(" ")
	buf.WriteString(node.rows())
	buf.WriteString("\n")

	for _, childNode := range node.Plans {
		buf.WriteString(childNode.View(level + 1))
	}

	return buf.String()
}

func (node PlanNode) name() string {
	return strings.Trim(fmt.Sprintf("%s %s", node.PartialMode, node.NodeType), " ")
}

func (node PlanNode) rows() string {
	return fmt.Sprintf("(Rows planned=%d actual=%d)", node.PlanRows, node.ActualRows)
}
