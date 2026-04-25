package algorithm

import "github.com/luis/Tubes2_AyamCarbonara/backend/src/model"

func FindLCA(root *model.DOMNode, id1, id2 int) *model.DOMNode {
	search1 := FastBFS(root, func(node *model.DOMNode) bool {return node.Id == id1})
	if (search1 == nil) {
		return  nil
	}
	search2 := FastBFS(root, func(node *model.DOMNode) bool {return node.Id == id2})
	if (search2 == nil) {
		return  nil
	}
	return FindLCAByNode(search1[0], search2[0])
}

func FindLCAByNode(node1, node2 *model.DOMNode) *model.DOMNode {
	for node1 != nil && node2 != nil {
		if (node1.Depth > node2.Depth) {
			node1 = node1.Parent
		} else if node1.Depth < node2.Depth {
			node2 = node2.Parent;
		} else {
			if node1 == node2 {
				return  node1;
			}
			node1 = node1.Parent
			node2 = node2.Parent
		}
	}
	return nil;
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
