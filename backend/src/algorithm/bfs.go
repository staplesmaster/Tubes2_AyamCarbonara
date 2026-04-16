package algorithm

import (
	"runtime"
	"sync"

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