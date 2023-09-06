package classes

type Logger interface {
	Log(Actor string, Action string, Target string, Value float64)
}

type BattleLog []Entry
type Entry struct {
	Actor  string
	Action string
	Target string
	Value  float64
}

func (b *BattleLog) Log(Actor string, Action string, Target string, Value float64) {
	*b = append(*b, Entry{Actor, Action, Target, Value})
}
