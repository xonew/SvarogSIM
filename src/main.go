package main

import (
	. "SvarogSIM/src/assets/characters"
	. "SvarogSIM/src/classes"
	"fmt"
)

func main() {
	order := MakeBattle(
		[]Ally{MakeHook(),
			MakeArlan()},
		[]Enemy{MakeWeakCocolia()})
	for ally, mp := range order.Run() {
		fmt.Printf("%s:", ally)
		for key, value := range mp {
			fmt.Printf("%s: %d\n", key, value)
		}
	}
}
