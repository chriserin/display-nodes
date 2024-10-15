package main

import (
	"fmt"
	"strings"
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

	var viewPosition position
	if ctx.JoinView {
		viewPosition = node.JoinViewPosition
	} else {
		viewPosition = node.Position
	}

	var styles Styles
	if ctx.Cursor == viewPosition.LineNumber {
		styles = ctx.CursorStyle
	} else if ctx.Cursor == viewPosition.Parent {
		styles = ctx.ChildCursorStyle
	} else {
		styles = ctx.NormalStyle
	}

	var buf strings.Builder
	buf.WriteString(styles.Gutter.Render(fmt.Sprintf("%2d ", viewPosition.LineNumber)))

	if ctx.Indent {
		buf.WriteString(styles.Everything.Render(strings.Repeat("  ", viewPosition.Level-1)))
	}

	if ctx.JoinView && node.RelationName != "" {
		buf.WriteString(styles.Relation.Render(node.RelationName))
	} else {
		buf.WriteString(styles.NodeName.Render(node.name() + " "))
		buf.WriteString(styles.Everything.Render(node.rows()))
	}

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
