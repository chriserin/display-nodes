package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type StatusLine struct {
	ExecutionTime float64
	TotalBuffers  int
	TotalRows     int
}

func (s StatusLine) View(ctx ProgramContext) string {
	color_a := lipgloss.Color("#9999bb")
	color_b := lipgloss.Color("#452297")
	color_c := lipgloss.Color("#000000")
	style_a := lipgloss.NewStyle().Background(color_a).Foreground(color_c)
	style_b := lipgloss.NewStyle().Background(color_c).Foreground(color_a)
	style_c := lipgloss.NewStyle().Background(color_a).Foreground(color_b)

	if ctx.Width < 30 {
		return ""
	}

	line :=
		style_b.Render("  ") +
			style_a.Render("") +
			style_a.Render(" Time:") +
			style_c.Render(" %.3fms ") +
			style_b.Render("") +
			style_a.Render("") +
			style_a.Render(" Buffers:") +
			style_c.Render(" %s ") +
			style_b.Render("") +
			style_a.Render("") +
			style_a.Render(" Rows:") +
			style_c.Render(" %s ") +
			style_b.Render("      %s")

	result := fmt.Sprintf(line,
		s.ExecutionTime,
		formatUnderscores(s.TotalBuffers),
		formatUnderscores(s.TotalRows),
		ctx.StatDisplay.String(),
	)

	needed := ctx.Width - ansi.StringWidth(result)
	space := style_b.Render(strings.Repeat(" ", needed))

	return result + space + "\n"
}
