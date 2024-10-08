package main

import (
	"encoding/json"
	"os"
)

func main() {

	decoded := decodeJson(os.Stdin)
	data := extractPlanNodes(decoded)
	runProgram(data)

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

func extractPlanNodes(plan map[string]interface{}) PlanNode {

	nodeType := plan["Node Type"].(string)
	planRows := plan["Plan Rows"].(float64)
	actualRows := plan["Actual Rows"].(float64)
	partialMode, ok := plan["Partial Mode"].(string)

	if !ok {
		partialMode = ""
	}

	plans := plan["Plans"]

	planNodes := make([]PlanNode, 0, 1)

	if plans != nil {
		for _, plan := range plans.([]interface{}) {
			if plan != nil {
				planNode := extractPlanNodes(plan.(map[string]interface{}))
				planNodes = append(planNodes, planNode)
			}
		}
	}

	return PlanNode{
		NodeType:    nodeType,
		Plans:       planNodes,
		PlanRows:    int(planRows),
		ActualRows:  int(actualRows),
		PartialMode: partialMode,
	}
}
