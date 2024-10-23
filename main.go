package main

import (
	"encoding/json"
	"os"
	"slices"
	"strings"
)

func main() {
	decoded, executionTime := decodeJson(os.Stdin)
	nodes := make([]PlanNode, 0, 1)
	id := 0
	extractPlanNodes(decoded,
		position{Id: 0, Level: 0, Parent: 0},
		position{Id: 0, Level: 0, Parent: 0},
		ParseContext{Id: &id, Nodes: &nodes},
	)
	runProgram(nodes, executionTime, InitProgramContext(nodes[0]))
}

func decodeJson(data *os.File) (map[string]interface{}, float64) {
	var decoded any

	err := json.NewDecoder(os.Stdin).Decode(&decoded)

	if err != nil {
		panic("panic!")
	}

	planObject := decoded.([]interface{})[0].(map[string]interface{})
	plan := planObject["Plan"].(map[string]interface{})
	executionTime := planObject["Execution Time"].(float64)

	return plan, executionTime
}

type position struct {
	Id          int
	Level       int
	Parent      int
	Display     bool
	BelowGather bool
}

type ParseContext struct {
	Id          *int
	Nodes       *[]PlanNode
	BelowGather bool
}

func extractPlanNodes(plan map[string]interface{}, parentPosition position, parentJoinPosition position, parseContext ParseContext) PlanNode {
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

	indexName, ok := plan["Index Name"].(string)
	if !ok {
		indexName = ""
	}

	indexCond, ok := plan["Index Cond"].(string)
	if !ok {
		indexCond = ""
	}

	filter, ok := plan["Filter"].(string)
	if !ok {
		filter = ""
	}

	sharedReadBlocks := plan["Shared Read Blocks"].(float64)
	sharedHitBlocks := plan["Shared Hit Blocks"].(float64)
	startupCost := plan["Startup Cost"].(float64)
	totalCost := plan["Total Cost"].(float64)
	startupTime := plan["Actual Startup Time"].(float64)
	totalTime := plan["Actual Total Time"].(float64)

	workersLaunched, ok := plan["Workers Launched"].(float64)

	plans := plan["Plans"]

	id := parseContext.Id
	*id = *id + 1

	isGather := strings.Contains(nodeType, "Gather")

	var workers int
	if isGather {
		workers = int(workersLaunched) + 1
	} else {
		workers = 0
	}

	newPosition := position{
		Id:          *id,
		Level:       parentPosition.Level + 1,
		Parent:      parentPosition.Id,
		Display:     true,
		BelowGather: parseContext.BelowGather,
	}

	var joinViewPosition position
	if isJoinType(nodeType) || relationName != "" || isGather {
		joinViewPosition = position{
			Id:          *id,
			Level:       parentJoinPosition.Level + 1,
			Parent:      parentJoinPosition.Id,
			Display:     true,
			BelowGather: parseContext.BelowGather,
		}
	} else {
		joinViewPosition = parentJoinPosition
		joinViewPosition.Display = false
	}

	extractedNode := PlanNode{
		NodeType:          nodeType,
		PlanRows:          int(planRows),
		ActualRows:        int(actualRows),
		PartialMode:       partialMode,
		Position:          newPosition,
		JoinViewPosition:  joinViewPosition,
		RelationName:      relationName,
		SharedBuffersHit:  int(sharedHitBlocks),
		SharedBuffersRead: int(sharedReadBlocks),
		IsGather:          isGather,
		Workers:           workers,
		StartupCost:       startupCost,
		TotalCost:         totalCost,
		StartupTime:       startupTime,
		TotalTime:         totalTime,
		IndexName:         indexName,
		IndexCond:         indexCond,
		Filter:            filter,
	}

	nodes := parseContext.Nodes

	*nodes = append(*nodes, extractedNode)

	if plans != nil {
		for _, plan := range plans.([]interface{}) {
			if plan != nil {
				extractPlanNodes(plan.(map[string]interface{}), newPosition, joinViewPosition, ParseContext{Id: id, Nodes: nodes, BelowGather: isGather || parseContext.BelowGather})
			}
		}
	}

	return extractedNode
}

func isJoinType(nodeType string) bool {
	return slices.Contains([]string{"Nested Loop", "Hash Join", "Merge Join"}, nodeType)
}
