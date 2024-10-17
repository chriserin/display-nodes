package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type StatusLine struct {
	ExecutionTime float64
	TotalBuffers  int
	TotalRows     int
}

func (s StatusLine) View(width int) string {
	style := lipgloss.NewStyle().Background(lipgloss.Color("#9999bb")).Foreground(lipgloss.Color("#000000"))

	if width < 30 {
		return ""
	}

	line := "   Execution Time: %.3fms ◆ Total Buffers: %s ◆ Total Rows: %s "

	result := fmt.Sprintf(line,
		s.ExecutionTime,
		formatUnderscores(s.TotalBuffers),
		formatUnderscores(s.TotalRows),
	)

	return style.Render(runewidth.FillRight(result, width)) + "\n"
}
