package main

import (
	"encoding/json"
	"os"
)

func main() {
	decoded := decodeJson(os.Stdin)
	nodes := make([]PlanNode, 0, 1)
	lineNumber := 0
	extractPlanNodes(decoded, 1, &lineNumber, &nodes)
	runProgram(nodes, ProgramContext{})
}

func decodeJson(data *os.File) map[string]interface{} {
	var decoded any

	err := json.NewDecoder(os.Stdin).Decode(&decoded)

	if err != nil {
		panic("panic!")
	}

	plan := decoded.([]interface{})[0].(map[string]interface{})["Plan"].(map[string]interface{})

	return plan
}

func extractPlanNodes(plan map[string]interface{}, level int, lineNumber *int, nodes *[]PlanNode) PlanNode {
	nodeType := plan["Node Type"].(string)
	planRows := plan["Plan Rows"].(float64)
	actualRows := plan["Actual Rows"].(float64)
	partialMode, ok := plan["Partial Mode"].(string)

	if !ok {
		partialMode = ""
	}

	plans := plan["Plans"]

	*lineNumber = *lineNumber + 1

	extractedNode := PlanNode{
		NodeType:    nodeType,
		PlanRows:    int(planRows),
		ActualRows:  int(actualRows),
		PartialMode: partialMode,
		LineNumber:  *lineNumber,
		Level:       level,
	}

	*nodes = append(*nodes, extractedNode)

	if plans != nil {
		for _, plan := range plans.([]interface{}) {
			if plan != nil {
				extractPlanNodes(plan.(map[string]interface{}), level+1, lineNumber, nodes)
			}
		}
	}

	return extractedNode
}
