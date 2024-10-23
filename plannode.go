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
	StartupCost       float64
	TotalCost         float64
	StartupTime       float64
	TotalTime         float64
	IndexName         string
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

	if ctx.DisplayParallel {
		if viewPosition.BelowGather {
			buf.WriteString(styles.Gutter.Render("┃┃ "))
		} else if node.Workers > 0 {
			buf.WriteString(styles.Workers.Render(fmt.Sprintf("%.2d ", node.Workers)))
		} else {
			buf.WriteString("   ")
		}
	}

	if ctx.Indent {
		buf.WriteString(styles.Everything.Render(strings.Repeat("  ", viewPosition.Level-1)))
	}

	if ctx.JoinView && node.RelationName != "" {
		buf.WriteString(styles.Relation.Render(node.RelationName))
	} else {
		buf.WriteString(styles.NodeName.Render(node.name()))
	}

	result := buf.String()

	needed := ctx.Width - ansi.StringWidth(result)

	if ctx.StatDisplay == DisplayRows {
		buf.WriteString(node.rows(styles, needed))
	} else if ctx.StatDisplay == DisplayBuffers {
		buf.WriteString(node.buffers(styles, needed))
	} else if ctx.StatDisplay == DisplayCost {
		buf.WriteString(node.costs(styles, needed))
	} else if ctx.StatDisplay == DisplayTime {
		buf.WriteString(node.times(styles, needed))
	} else if ctx.StatDisplay == DisplayNothing {
		buf.WriteString(styles.Everything.Render(fmt.Sprintf("%*s", needed, "")))
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

func (node PlanNode) buffers(styles Styles, space int) string {
	var buf strings.Builder
	totalBuffers := formatUnderscores(node.SharedBuffersRead + node.SharedBuffersHit)
	readBuffers := formatUnderscores(node.SharedBuffersRead)

	if false {
		buf.WriteString(styles.Bracket.Render(" ["))
		buf.WriteString(styles.Everything.Render("total="))
		buf.WriteString(totalBuffers)
		buf.WriteString(styles.Everything.Render(" read="))
		buf.WriteString(readBuffers)
		buf.WriteString(styles.Bracket.Render("]"))
	} else {
		columns := fmt.Sprintf("%15s%15s", totalBuffers, readBuffers)
		buf.WriteString(styles.Value.Render(fmt.Sprintf("%*s", space, columns)))
	}

	return buf.String()
}

func (node PlanNode) costs(styles Styles, space int) string {
	startupCost := formatUnderscoresFloat(node.StartupCost)
	totalCost := formatUnderscoresFloat(node.TotalCost)

	var buf strings.Builder
	if false {
		buf.WriteString(styles.Bracket.Render(" ["))
		buf.WriteString(styles.Everything.Render("startup="))
		buf.WriteString(styles.Value.Render(startupCost))
		buf.WriteString(styles.Everything.Render(" total="))
		buf.WriteString(styles.Value.Render(totalCost))
		buf.WriteString(styles.Bracket.Render("]"))
	} else {
		columns := fmt.Sprintf("%15s%15s", startupCost, totalCost)
		buf.WriteString(styles.Value.Render(fmt.Sprintf("%*s", space, columns)))
	}

	return buf.String()
}

func (node PlanNode) times(styles Styles, space int) string {
	startupTime := formatUnderscoresFloat(node.StartupTime)
	totalTime := formatUnderscoresFloat(node.TotalTime)

	var buf strings.Builder
	if false {
		buf.WriteString(styles.Bracket.Render(" ["))
		buf.WriteString(styles.Everything.Render("startup="))
		buf.WriteString(styles.Value.Render(startupTime))
		buf.WriteString(styles.Everything.Render(" total="))
		buf.WriteString(styles.Value.Render(totalTime))
		buf.WriteString(styles.Bracket.Render("]"))
	} else {
		columns := fmt.Sprintf("%15s%15s", startupTime, totalTime)
		buf.WriteString(styles.Value.Render(fmt.Sprintf("%*s", space, columns)))
	}

	return buf.String()
}

func (node PlanNode) rows(styles Styles, space int) string {

	separatedPlanRows := formatUnderscores(node.PlanRows)
	separatedActualRows := formatUnderscores(node.ActualRows)

	percentOfActual := float32(node.PlanRows) / float32(node.ActualRows) * 100

	rowStatus := getRowStatus(percentOfActual, styles)

	var buf strings.Builder
	if false {
		buf.WriteString(styles.Bracket.Render(" ["))
		buf.WriteString(styles.Everything.Render("p="))
		buf.WriteString(styles.Value.Render(separatedPlanRows))
		buf.WriteString(styles.Everything.Render(" "))
		buf.WriteString(styles.Everything.Render("a="))
		buf.WriteString(styles.Value.Render(separatedActualRows))
		buf.WriteString(rowStatus)
		buf.WriteString(styles.Bracket.Render("]"))
	} else {
		columns := fmt.Sprintf("%15s%15s", separatedPlanRows, separatedActualRows)
		buf.WriteString(styles.Value.Render(fmt.Sprintf("%*s", space, columns)))
	}

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
	return strings.Replace(printer.Sprintf("%d", value), ",", "_", -1)
}

func formatUnderscoresFloat(value float64) string {
	return strings.Replace(printer.Sprintf("%.2f", value), ",", "_", -1)
}

func (node PlanNode) Content(ctx ProgramContext) string {
	var buf strings.Builder

	buf.WriteString(ctx.NormalStyle.NodeName.Render(node.name()))
	buf.WriteString("\n")
	buf.WriteString(strings.Repeat("-", ctx.Width))
	buf.WriteString("\n")
	buf.WriteString(ctx.NormalStyle.Relation.Render(node.RelationName))
	buf.WriteString("\n")
	buf.WriteString(ctx.NormalStyle.Everything.Render(node.IndexName))
	buf.WriteString("\n")

	return buf.String()
}
