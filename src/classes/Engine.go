package classes

import (
	"sort"
)

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []Actor

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].GetActionValue() < pq[j].GetActionValue()
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	item := x.(Actor)
	*pq = append(*pq, item)
	pq.Sort()
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	item := old[0]
	old[0] = nil // avoid memory leak
	*pq = old[0:]
	return item
}

func (pq *PriorityQueue) Sort() {
	sort.Sort(pq)
}

// startTurn starts the next turn and returns the number of action value elapsed
func (a Combat) startTurn() int {
	actionValue := (*a.Pq)[0].GetActionValue()
	for i := 0; i < a.Pq.Len(); i++ {
		(*a.Pq)[i].ModifyActionValue(-actionValue)
	}
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

		pq.Push(allies[i])
		heapify := combat.Pq.Sort
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
		pq.Push(enemies[i])
		heapify := combat.Pq.Sort
		enemies[i].Init(left, right, heapify)
		enemies[i].InitBattle(&combat)
	}
	return combat
}

func (c *Combat) Run() map[string]map[string]int {
	totalActionValue := 0
	for totalActionValue < 850 {
		c.Pq.Sort()
		totalActionValue += c.startTurn()
		(*c.Pq)[0].Act()
	}
	cumulativeDamage := make(map[string]map[string]int)
	for _, attacker := range c.Allies {
		cumulativeDamage[attacker.GetName()] = make(map[string]int)
		for attackType, attacks := range attacker.GetDamageOutLog() {
			cumulative := 0
			for _, attack := range attacks {
				cumulative += attack.PostMitDamage
			}
			cumulativeDamage[attacker.GetName()][attackType] = cumulative
		}
	}

	return cumulativeDamage
}
