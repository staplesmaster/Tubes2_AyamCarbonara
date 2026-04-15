package algorithm

import (
	"github.com/luis/Tubes2_AyamCarbonara/backend/src/model"
	"github.com/luis/Tubes2_AyamCarbonara/backend/src/selector"
)

func BFS(node *model.DOMNode, filter selector.Selector) []*model.DOMNode {
	var queue []*model.DOMNode
	queue = append(queue, node)

	var result []*model.DOMNode
	var current *model.DOMNode

	for len(queue) > 0 {
		current = queue[0]
		queue = append(queue, current.Children...)
		queue = queue[1:]

		if (filter(current)) {
			result = append(result, current)
		}
	}

	return result
}