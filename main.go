package main

import (
	"encoding/json"
	"os"
)

func main() {
	decoded := decodeJson(os.Stdin)
	nodes := make([]PlanNode, 0, 1)
	lineNumber := 0
	extractPlanNodes(decoded, position{LineNumber: 0, Level: 0, Parent: 0}, &lineNumber, &nodes)
	runProgram(nodes, ProgramContext{Cursor: 1})
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

type position struct {
	LineNumber int
	Level      int
	Parent     int
}

func extractPlanNodes(plan map[string]interface{}, parentPosition position, lineNumber *int, nodes *[]PlanNode) PlanNode {
	nodeType := plan["Node Type"].(string)
	planRows := plan["Plan Rows"].(float64)
	actualRows := plan["Actual Rows"].(float64)
	partialMode, ok := plan["Partial Mode"].(string)

	if !ok {
		partialMode = ""
	}

	relationName, ok := plan["Relation Name"].(string)

	if !ok {
		relationName = ""
	}

	plans := plan["Plans"]

	*lineNumber = *lineNumber + 1

	newPosition := position{
		LineNumber: *lineNumber,
		Level:      parentPosition.Level + 1,
		Parent:     parentPosition.LineNumber,
	}

	extractedNode := PlanNode{
		NodeType:     nodeType,
		PlanRows:     int(planRows),
		ActualRows:   int(actualRows),
		PartialMode:  partialMode,
		Position:     newPosition,
		RelationName: relationName,
	}

	*nodes = append(*nodes, extractedNode)

	if plans != nil {
		for _, plan := range plans.([]interface{}) {
			if plan != nil {
				extractPlanNodes(plan.(map[string]interface{}), newPosition, lineNumber, nodes)
			}
		}
	}

	return extractedNode
}
