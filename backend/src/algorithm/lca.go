package algorithm

import "github.com/luis/Tubes2_AyamCarbonara/backend/src/model"

func FindLCA(root *model.DOMNode, id1, id2 int) *model.DOMNode {
	lca, _ := FindLCAWithSteps(root, id1, id2)
	return lca
}

func FindLCAWithSteps(root *model.DOMNode, id1, id2 int) (*model.DOMNode, []model.LCAStep) {
	node1 := findNodeByID(root, id1)
	if node1 == nil {
		return nil, nil
	}

	node2 := findNodeByID(root, id2)
	if node2 == nil {
		return nil, nil
	}

	return FindLCAByNodeWithSteps(node1, node2)
}

func FindLCAByNode(node1, node2 *model.DOMNode) *model.DOMNode {
	lca, _ := FindLCAByNodeWithSteps(node1, node2)
	return lca
}

func FindLCAByNodeWithSteps(node1, node2 *model.DOMNode) (*model.DOMNode, []model.LCAStep) {
	var steps []model.LCAStep
	step := 0

	for node1 != nil && node2 != nil {
		step++
		activeIDs := []int{node1.Id}
		if node2.Id != node1.Id {
			activeIDs = append(activeIDs, node2.Id)
		}

		if node1 == node2 {
			lcaID := node1.Id
			steps = append(steps, model.LCAStep{
				Step:          step,
				NodeA:         node1.Id,
				NodeB:         node2.Id,
				ActiveNodeIDs: activeIDs,
				LCANodeID:     &lcaID,
			})
			return node1, steps
		}

		steps = append(steps, model.LCAStep{
			Step:          step,
			NodeA:         node1.Id,
			NodeB:         node2.Id,
			ActiveNodeIDs: activeIDs,
		})

		if node1.Depth > node2.Depth {
			node1 = node1.Parent
		} else if node2.Depth > node1.Depth {
			node2 = node2.Parent
		} else {
			node1 = node1.Parent
			node2 = node2.Parent
		}
	}

	return nil, steps
}

func findNodeByID(root *model.DOMNode, id int) *model.DOMNode {
	matches := BFS(root, func(node *model.DOMNode) bool { return node.Id == id })
	if len(matches) == 0 {
		return nil
	}
	return matches[0]
}
