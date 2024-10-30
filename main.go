package main

import (
	"io"
	"os"
)

func main() {
	input, _ := io.ReadAll(os.Stdin)
	explainPlan := Convert(string(input))
	runProgram(explainPlan)
}
