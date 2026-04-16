package algorithm

import (
	"github.com/luis/Tubes2_AyamCarbonara/backend/src/model"
	"github.com/luis/Tubes2_AyamCarbonara/backend/src/selector"
)

func DFS(node *model.DOMNode, filter selector.Selector) []*model.DOMNode {
	stack := []*model.DOMNode{node}
	result := make([]*model.DOMNode, 0, 32)

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if filter(current) {
			result = append(result, current)
		}

		for i := len(current.Children) - 1; i >= 0; i-- {
			stack = append(stack, current.Children[i])
		}
	}

	return result
}