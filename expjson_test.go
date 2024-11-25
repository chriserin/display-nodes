package main

import (
	"os"
	"testing"
)

func TestConvertNoBuffers(t *testing.T) {
	data, err := os.ReadFile("./testdata/analyze_no_buffers.json")
	if err != nil {
		t.Fatal(err)
	}
	plan := Convert(string(data))
	if plan.executionTime != 69.662 {
		t.Fatalf("GOT %f WANT %f", plan.executionTime, 69.662)
	}
}
