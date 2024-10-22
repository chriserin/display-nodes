package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

type StatusLine struct {
	ExecutionTime float64
	TotalBuffers  int
	TotalRows     int
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
	buf.WriteString(styles.AltNormal.Render("      %s"))

	result := fmt.Sprintf(buf.String(),
		s.ExecutionTime,
		formatUnderscores(s.TotalBuffers),
		formatUnderscores(s.TotalRows),
		ctx.StatDisplay.String(),
	)

	var finalBuf strings.Builder

	finalBuf.WriteString(result)

	needed := ctx.Width - ansi.StringWidth(result)
	space := styles.AltNormal.Render(strings.Repeat(" ", needed))

	finalBuf.WriteString(space)
	finalBuf.WriteString("\n")

	return finalBuf.String()
}
