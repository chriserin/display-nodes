package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type PlanNode struct {
	NodeType          string
	Plans             []PlanNode
	PlanRows          int
	ActualRows        int
	PartialMode       string
	Position          position
	JoinViewPosition  position
	RelationName      string
	SharedBuffersRead int
	SharedBuffersHit  int
	IsGather          bool
	Workers           int
}

func (node PlanNode) View(i int, ctx ProgramContext) string {

	var viewPosition position
	if ctx.JoinView {
		viewPosition = node.JoinViewPosition
	} else {
		viewPosition = node.Position
	}

	var styles Styles
	if ctx.Cursor == i {
		styles = ctx.CursorStyle
	} else if ctx.SelectedNode.Position.Id == viewPosition.Parent {
		styles = ctx.ChildCursorStyle
	} else {
		styles = ctx.NormalStyle
	}

	var buf strings.Builder
	buf.WriteString(styles.Gutter.Render(fmt.Sprintf("%2d ", i+1)))
	if viewPosition.BelowGather {
		buf.WriteString(styles.Gutter.Render("┃┃"))
	} else if node.Workers > 0 {
		buf.WriteString(styles.Gutter.Render(fmt.Sprintf("%.2d", node.Workers)))
	} else {
		buf.WriteString("  ")
	}

	if ctx.Indent {
		buf.WriteString(styles.Everything.Render(strings.Repeat("  ", viewPosition.Level-1)))
	}

	if ctx.JoinView && node.RelationName != "" {
		buf.WriteString(styles.Relation.Render(node.RelationName))
	} else {
		buf.WriteString(styles.NodeName.Render(node.name()))
	}

	if ctx.DisplayRows {
		buf.WriteString(node.rows(styles))
	} else if ctx.DisplayBuffers {
		buf.WriteString(node.buffers(styles))
	}

	result := buf.String()

	needed := ctx.Width - ansi.StringWidth(result)
	if needed > 0 {
		buf.WriteString(styles.Everything.Render(strings.Repeat(" ", needed)))
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

func (node PlanNode) buffers(styles Styles) string {
	var buf strings.Builder
	buf.WriteString(styles.Bracket.Render(" Buffers["))
	buf.WriteString(styles.Everything.Render("total="))
	buf.WriteString(styles.Value.Render(formatUnderscores(node.SharedBuffersRead + node.SharedBuffersHit)))
	buf.WriteString(styles.Everything.Render(" read="))
	buf.WriteString(styles.Value.Render(formatUnderscores(node.SharedBuffersRead)))
	buf.WriteString(styles.Bracket.Render("]"))

	return buf.String()
}

func (node PlanNode) rows(styles Styles) string {

	var buf strings.Builder

	separatedPlanRows := formatUnderscores(node.PlanRows)
	separatedActualRows := formatUnderscores(node.ActualRows)

	percentOfActual := float32(node.PlanRows) / float32(node.ActualRows) * 100

	rowStatus := getRowStatus(percentOfActual, styles)

	buf.WriteString(styles.Bracket.Render(" Rows["))
	buf.WriteString(styles.Everything.Render("p="))
	buf.WriteString(styles.Value.Render(separatedPlanRows))
	buf.WriteString(styles.Everything.Render(" "))
	buf.WriteString(styles.Everything.Render("a="))
	buf.WriteString(styles.Value.Render(separatedActualRows))
	buf.WriteString(rowStatus)
	buf.WriteString(styles.Bracket.Render("]"))

	return buf.String()
}

func getRowStatus(percentOfActual float32, styles Styles) string {
	if percentOfActual < 10 {
		return styles.Warning.Render(fmt.Sprintf(" %.1f%%", percentOfActual))
	} else if percentOfActual < 50 {
		return styles.Caution.Render(fmt.Sprintf(" %.1f%%", percentOfActual))
	} else {
		return styles.Everything.Render(fmt.Sprintf(" %.1f%%", percentOfActual))
	}
}

var printer *message.Printer = message.NewPrinter(language.English)

func formatUnderscores(value int) string {
	return strings.Replace(printer.Sprintf("%d", value), ",", "_", 1)
}
