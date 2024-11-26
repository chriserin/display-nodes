package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

const (
	ROW_HEADER = iota
	ROW_STATUS
	ROW_EX_NODE_1
	ROW_EX_NODE_2
	ROW_EX_NODE_3
	ROW_EX_NODE_4
	ROW_BLANK
	ROW_DETAILS_TITLE
)

func TestStdinNoBuffers(t *testing.T) {

	var buf bytes.Buffer

	dat, err := os.ReadFile("./testdata/analyze_no_buffers.json")
	if err != nil {
		t.Fatal(err)
	}
	source := Source{sourceType: SOURCE_STDIN, input: string(dat)}
	p := RunProgram(source, tea.WithOutput(&buf))

	cmds := []tea.Cmd{
		func() tea.Msg { return tea.WindowSizeMsg{Width: 80, Height: 24} },
		tea.Quit,
	}
	seq := tea.Sequence(cmds...)
	go p.Send(seq())

	if _, err := p.Run(); err != nil {
		t.Fatal(err)
	}

	rendered := strings.Split(buf.String(), "\n")

	assert.Contains(t, rendered[ROW_HEADER], "STDIN")
	assert.Contains(t, rendered[ROW_STATUS], "Time: 69.662ms")
	assert.Contains(t, rendered[ROW_EX_NODE_1], "Finalize Aggregate")
	assert.Contains(t, rendered[ROW_DETAILS_TITLE], "Details")
}
