package algorithm

import (
	"runtime"
	"sync"

	"github.com/luis/Tubes2_AyamCarbonara/backend/src/model"
	"github.com/luis/Tubes2_AyamCarbonara/backend/src/selector"
)

func BFS(node *model.DOMNode, filter selector.Selector) []*model.DOMNode {
	queue := []*model.DOMNode{node}
	result := make([]*model.DOMNode, 0, 32)

	head := 0
	for head < len(queue) {
		current := queue[head]
		head++

		queue = append(queue, current.Children...)

		if (filter(current)) {
			result = append(result, current)
		}
	}

	return result
}

func FastBFS(node *model.DOMNode, filter selector.Selector) []*model.DOMNode {
	queue := make(chan *model.DOMNode, 100)
    resultChan := make(chan *model.DOMNode, 100)

    var wg sync.WaitGroup

    worker := func() {
        for n := range queue {
            if filter(n) {
                resultChan <- n
            }

            for _, child := range n.Children {
                wg.Add(1)
                queue <- child
            }
            wg.Done()
        }
    }

    cpu := runtime.NumCPU()
    for i := 0; i < cpu; i++ {
        go worker()
    }

    wg.Add(1)
    queue <- node

    go func() {
        wg.Wait()
        close(queue)
        close(resultChan)
    }()

    var result []*model.DOMNode
    for r := range resultChan {
        result = append(result, r)
    }

    return result
}