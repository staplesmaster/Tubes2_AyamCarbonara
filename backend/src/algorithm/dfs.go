package algorithm

import (
	"github.com/luis/Tubes2_AyamCarbonara/backend/src/model"
	"github.com/luis/Tubes2_AyamCarbonara/backend/src/selector"
)

func DFS(node *model.DOMNode, filter selector.Selector) []*model.DOMNode {
	var result []*model.DOMNode
	if (filter(node)) {
		result = append(result, node)
	}
	for _, child := range node.Children {
		childResult := DFS(child, filter)
		result = append(result, childResult...)
	}
	return result
}