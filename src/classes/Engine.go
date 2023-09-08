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

// nextTurn does the next turn and returns the number of action value elapsed
func (a Combat) nextTurn() int {
	actionValue := (*a.Pq)[0].priority
	for i := 0; i < a.Pq.Len(); i++ {
		(*a.Pq)[i].priority -= actionValue
	}

	a.Pq.update((*a.Pq)[0], (*a.Pq)[0].actor.GetActionValue())
	return actionValue
}

type Combat struct {
	Pq          *PriorityQueue
	Enemies     []Enemy
	Allies      []Ally
	SkillPoints SkillPoints
}

func MakeBattle(allies []Ally, enemies []Enemy) Combat {
	pq := make(PriorityQueue, 0)
	combat := Combat{
		Pq:          &pq,
		Enemies:     enemies,
		Allies:      allies,
		SkillPoints: &SkillPoint{pool: 0, max: 5},
	}
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
		allies[i].InitBattle(&combat)
	}

	for i := 0; i < len(enemies); i++ {
		var left Enemy
		var right Enemy
		if i-1 < 0 {
			left = nil
		} else {
			left = enemies[i-1]
		}
		if i >= len(enemies)-1 {
			right = nil
		} else {
			right = enemies[i+1]
		}

		item := &Item{
			actor:    enemies[i],
			priority: enemies[i].GetActionValue(),
		}
		pq.Push(item)
		heapify := func() {
			item.priority = item.actor.GetActionValue()
			pq.update(item, item.priority)
		}
		enemies[i].Init(left, right, heapify)
		enemies[i].InitBattle(&combat)
	}
	return combat
}

func (c *Combat) Run() map[string]int {
	totalActionValue := 0
	for totalActionValue < 850 {
		(*c.Pq)[0].actor.Act()
		totalActionValue += c.nextTurn()
	}
	cumulativeDamage := make(map[string]int)
	for _, attacker := range c.Allies {
		for attackType, attacks := range attacker.GetDamageOutLog() {
			cumulative := 0
			for _, attack := range attacks {
				cumulative += attack.PostMitDamage
			}
			cumulativeDamage[attackType] = cumulative
		}
	}

	return cumulativeDamage
}
