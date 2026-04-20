package algorithm

import (
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/luis/Tubes2_AyamCarbonara/backend/src/model"
	"github.com/luis/Tubes2_AyamCarbonara/backend/src/selector"
)

func DFS(root *model.DOMNode, filter selector.Selector) []*model.DOMNode {
	stack := []*model.DOMNode{root}
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

type Deque struct {
	mu    sync.Mutex
	items []*model.DOMNode
	head  int
}

func (this *Deque) pushBottom(n *model.DOMNode) {
	this.mu.Lock()
	this.items = append(this.items, n)
	this.mu.Unlock()
}

func (this *Deque) stealTop() *model.DOMNode {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.head >= len(this.items) {
		return nil
	}

	n := this.items[this.head]
	this.items[this.head] = nil
	this.head++

	if this.head > 1024 && this.head > len(this.items)/2 {
		copy(this.items, this.items[this.head:])
		this.items = this.items[:len(this.items)-this.head]
		this.head = 0
	}

	return n
}

func (this *Deque) popBottom() *model.DOMNode {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.head >= len(this.items) {
		return nil
	}

	last := len(this.items) - 1
	n := this.items[last]
	this.items[last] = nil
	this.items = this.items[:last]

	return n
}

func FastDFS(root *model.DOMNode, filter func(*model.DOMNode) bool) []*model.DOMNode {
	numWorkers := runtime.NumCPU()

	deques := make([]*Deque, numWorkers)
	for i := range deques {
		deques[i] = &Deque{}
	}

	var taskCount int64
	var wg sync.WaitGroup

	results := make([][]*model.DOMNode, numWorkers)

	rngs := make([]*rand.Rand, numWorkers)
	for i := 0; i < numWorkers; i++ {
		rngs[i] = rand.New(rand.NewSource(time.Now().UnixNano() + int64(i)))
	}

	atomic.AddInt64(&taskCount, 1)
	wg.Add(1)
	deques[0].pushBottom(root)

	worker := func(id int) {
		local := deques[id]
		rng := rngs[id]

		for {
			var root *model.DOMNode

			root = local.popBottom()

			if root == nil {
				start := rng.Intn(numWorkers)
				for i := 0; i < numWorkers; i++ {
					victim := (start + i) % numWorkers
					if victim == id {
						continue
					}
					root = deques[victim].stealTop()
					if root != nil {
						break
					}
				}
			}

			if root == nil {
				if atomic.LoadInt64(&taskCount) == 0 {
					return
				}
				runtime.Gosched()
				continue
			}

			if filter(root) {
				results[id] = append(results[id], root)
			}

			children := root.Children
			for i := len(children) - 1; i >= 0; i-- {
				atomic.AddInt64(&taskCount, 1)
				local.pushBottom(children[i])
				wg.Add(1)
			}

			wg.Done()
			atomic.AddInt64(&taskCount, -1)
		}
	}

	for i := 0; i < numWorkers; i++ {
		go worker(i)
	}

	wg.Wait()

	var result []*model.DOMNode
	for i := 0; i < numWorkers; i++ {
		result = append(result, results[i]...)
	}

	return result
}