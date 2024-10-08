package main

import "strings"

type PlanNode struct {
	NodeType   string
	Plans      []PlanNode
	PlanRows   int
	ActualRows int
}

func (node PlanNode) View() string {
	var buf strings.Builder

	buf.WriteString(node.NodeType + "\n")

	for _, childNode := range node.Plans {
		buf.WriteString(childNode.View())
	}

	return buf.String()
}
