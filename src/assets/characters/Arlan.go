package characters

import . "SvarogSIM/src/classes"

type Arlan struct {
	Character
}

func MakeArlan() *Arlan {
	return &Arlan{
		Character: MakeCharacter("arlan", "fire", "destruction",
			0, 80, 110,
			1199, 599, 330, 102),
	}
}

func (a *Arlan) Init(left Ally, right Ally, heapify func()) {
	a.EffectResist += 0.04 + 0.08 + 0.08
	a.Hp.Percent += 0.04 + 0.06
	a.Atk.Percent += 0.04 + 0.04 + 0.06 + 0.06 + 0.08
	a.ActionValue = int(15000 / a.Spd.GetStat())
	a.CurrHp = int(a.Hp.GetStat())
	a.Heapify = heapify
	a.Left = left
	a.Right = right
}

func (a *Arlan) Act() {
	target := Target(a.Battle.Enemies)
	target.TakeDamage(a.MakeAttack("skill", target.GetName(), "lightning", "skill", map[Stat]float64{
		a.Atk: 2.64,
	}))
}
