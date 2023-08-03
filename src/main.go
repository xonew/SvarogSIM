package main

//import (
//	. "SvarogSIM/src/classes"
//)

func main() {
	//var hook Actor = MakeHook()
	//order := MakeActionOrder(
	//	[]*Actor{&hook},
	//	[]*Actor{MakeHook()})

	thinger := here{value: 5}
	thingest := thing{value: 0, functo: thinger.increase}
	thingest.functo()
	println(thinger.value)
}

// TODO: Actors do not need to be pointers because they already are pointers, but attacks must be pointers
type here struct {
	value int
}
type thing struct {
	value  int
	functo func()
}

func (t *here) increase() {
	t.value += 100
}
