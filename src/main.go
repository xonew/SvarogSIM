package main

import (
	. "SvarogSIM/src/classes"
	"fmt"
)
import (
	. "SvarogSIM/src/assets/characters"
)

func main() {
	order := MakeBattle(
		[]Ally{MakeHook(),
			MakeHook()},
		[]Enemy{MakeWeakCocolia()})
	for key, value := range order.Run() {
		fmt.Printf("%s: %d\n", key, value)
	}
}
