package classes

import (
	"container/heap"
)

// An Item is something we manage in a priority queue.
type Item struct {
	actor    Actor // The actor of the item; arbitrary.
	priority int   // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and actor of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, priority int) {
	item.priority = priority
	heap.Fix(pq, item.index)
}

func (a ActionOrder) nextTurn() {
	actionValue := a.pq[0].priority
	for i := 0; i < a.pq.Len(); i++ {
		a.pq[i].priority -= actionValue
	}

	a.pq.update(a.pq[0], a.pq[0].actor.GetActionValue())
}

type ActionOrder struct {
	pq          PriorityQueue
	enemies     []Enemy
	allies      []Ally
	skillPoints SkillPoints
}

func MakeActionOrder(allies []Ally, enemies []Enemy) ActionOrder {
	pq := make(PriorityQueue, 0)
	for i := 0; i < len(allies); i++ {
		var left Ally
		var right Ally
		if i-1 < 0 {
			left = nil
		} else {
			left = allies[i-1]
		}
		if i >= len(allies)-1 {
			right = nil
		} else {
			right = allies[i+1]
		}

		item := &Item{
			actor:    allies[i],
			priority: allies[i].GetActionValue(),
		}
		pq.Push(item)
		heapify := func() {
			item.priority = item.actor.GetActionValue()
			pq.update(item, item.priority)
		}
		allies[i].Init(left, right, heapify)
	}
	pq.Push()
	return ActionOrder{
		pq:          pq,
		enemies:     enemies,
		allies:      allies,
		skillPoints: &SkillPoint{pool: 0, max: 5},
	}
}
