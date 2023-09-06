package main

import (
	. "SvarogSIM/src/classes"
)
import (
	. "SvarogSIM/src/assets/characters"
)

func main() {
	order := MakeActionOrder(
		[]Ally{MakeHook()},
		[]Enemy{MakeWeakCocolia()})

}
