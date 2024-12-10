package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlanNodeName(t *testing.T) {
	node := PlanNode{NodeType: "Something"}
	assert.Equal(t, "Something", node.Name())
}

func TestName(t *testing.T) {
	testCases := []struct {
		desc   string
		node   PlanNode
		result string
	}{
		{
			desc:   "Partial Mode",
			node:   PlanNode{NodeType: "Aggregate", PartialMode: "Partial"},
			result: "Partial Aggregate",
		},
		{
			desc:   "Modifying a table",
			node:   PlanNode{NodeType: "ModifyTable", Operation: "Merge"},
			result: "Merge",
		},
		{
			desc:   "Join Type with Hash Join",
			node:   PlanNode{NodeType: "Hash Join", JoinType: "Right"},
			result: "Hash Right Join",
		},
		{
			desc:   "Setop with strategy and command",
			node:   PlanNode{NodeType: "SetOp", Strategy: "Hashed", Command: "Except"},
			result: "HashSetOp Except",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Equal(t, tC.result, tC.node.Name())
		})
	}
}
