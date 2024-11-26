package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertNoBuffers(t *testing.T) {
	data, err := os.ReadFile("./testdata/analyze_no_buffers.json")
	if err != nil {
		t.Fatal(err)
	}
	plan := Convert(string(data))
	assert.Equal(t, plan.executionTime, 69.662)
}

func TestConvertNoWriteBuffers(t *testing.T) {
	data, err := os.ReadFile("./testdata/analyze_buffers.json")
	if err != nil {
		t.Fatal(err)
	}
	plan := Convert(string(data))
	assert.Equal(t, plan.nodes[0].Analyzed.TempWriteBlocks, 0)
}
