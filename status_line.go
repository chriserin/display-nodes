package main

import (
	"fmt"
	"strings"
)

type StatusLine struct {
	ExecutionTime float64
	TotalBuffers  int
	TotalRows     int
}

func NewStatusLine(explainPlan ExplainPlan) StatusLine {
	return StatusLine{
		ExecutionTime: explainPlan.executionTime,
		TotalBuffers:  explainPlan.TotalBuffers(),
		TotalRows:     explainPlan.TotalRows(),
	}
}

func (s StatusLine) View(ctx ProgramContext) string {
	styles := ctx.StatusStyles

	if ctx.Width < 30 {
		return ""
	}

	var buf strings.Builder

	buf.WriteString(styles.AltNormal.Render("  "))
	buf.WriteString(styles.Normal.Render(""))
	buf.WriteString(styles.Normal.Render(" Time:"))
	buf.WriteString(styles.Value.Render(" %.3fms "))
	buf.WriteString(styles.AltNormal.Render(""))
	buf.WriteString(styles.Normal.Render(""))
	buf.WriteString(styles.Normal.Render(" Buffers:"))
	buf.WriteString(styles.Value.Render(" %s "))
	buf.WriteString(styles.AltNormal.Render(""))
	buf.WriteString(styles.Normal.Render(""))
	buf.WriteString(styles.Normal.Render(" Rows:"))
	buf.WriteString(styles.Value.Render(" %s "))
	buf.WriteString(styles.AltNormal.Render(" "))

	result := fmt.Sprintf(buf.String(),
		s.ExecutionTime,
		formatUnderscores(s.TotalBuffers),
		formatUnderscores(s.TotalRows),
	)

	return result
}
