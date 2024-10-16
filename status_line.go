package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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
	result := fmt.Sprintf("   Execution Time: %.3fms ◆ Total Buffers: %s ◆ Total Rows: %s ",
		s.ExecutionTime,
		formatUnderscores(s.TotalBuffers),
		formatUnderscores(s.TotalRows),
	)
	additionalWidthNeeded := width - len(result)
	result += strings.Repeat(" ", additionalWidthNeeded)
	return style.Render(result) + "\n"
}
