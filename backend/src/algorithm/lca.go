package algorithm

import "github.com/luis/Tubes2_AyamCarbonara/backend/src/model"

func FindLCA(root *model.DOMNode, id1, id2 int) *model.DOMNode {
	search1 := BFS(root, func(node *model.DOMNode) bool {return node.Id == id1})
	if len(search1) == 0 {
		return  nil
	}
	search2 := BFS(root, func(node *model.DOMNode) bool { return node.Id == id2 })
	if len(search2) == 0 {
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