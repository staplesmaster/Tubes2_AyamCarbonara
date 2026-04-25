package algorithm

import (
	"sync"
	"time"

	"github.com/luis/Tubes2_AyamCarbonara/backend/src/model"
	"github.com/luis/Tubes2_AyamCarbonara/backend/src/selector"
)

func BFS(root *model.DOMNode, filter selector.Selector) []*model.DOMNode {
	queue := []*model.DOMNode{root}
	result := make([]*model.DOMNode, 0, 32)

	head := 0
	for head < len(queue) {
		current := queue[head]
		head++

		queue = append(queue, current.Children...)

		if filter(current) {
			result = append(result, current)
		}

		if head > 1024 {
			copy(queue, queue[head:])
			queue = queue[:len(queue)-head]
			head = 0
		}
	}

	return result
}

func BFSWithSteps(root *model.DOMNode, sel selector.Selector) ([]model.TraversalStep, []int, model.TraversalStats) {
	start := time.Now()

	var steps []model.TraversalStep
	var matchedIDs []int
	visited, matched, maxDepth, step := 0, 0, 0, 0

	queue := []*model.DOMNode{root}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		visitStep(cur, sel, &step, &visited, &matched, &maxDepth, &steps, &matchedIDs)

		queue = append(queue, cur.Children...)
	}

	stats := model.TraversalStats{
		Visited:  visited,
		Matched:  matched,
		MaxDepth: maxDepth,
		Elapsed:  float64(time.Since(start).Microseconds()) / 1000,
	}

	return steps, matchedIDs, stats
}

func FastBFS(root *model.DOMNode, filter selector.Selector) []*model.DOMNode {
	if root == nil {
		return nil
	}

	currentLevel := []*model.DOMNode{root}
	result := make([]*model.DOMNode, 0, 32)

	for len(currentLevel) > 0 {
		nextLevelChunks := make([][]*model.DOMNode, len(currentLevel))
		matches := make([]*model.DOMNode, len(currentLevel))

		var wg sync.WaitGroup
		for i, node := range currentLevel {
			wg.Add(1)
			go func(i int, node *model.DOMNode) {
				defer wg.Done()

				if filter(node) {
					matches[i] = node
				}
				nextLevelChunks[i] = node.Children
			}(i, node)
		}
		wg.Wait()

		for _, match := range matches {
			if match != nil {
				result = append(result, match)
			}
		}

		nextLevelSize := 0
		for _, children := range nextLevelChunks {
			nextLevelSize += len(children)
		}

		nextLevel := make([]*model.DOMNode, 0, nextLevelSize)
		for _, children := range nextLevelChunks {
			nextLevel = append(nextLevel, children...)
		}
		currentLevel = nextLevel
	}

	return result
}