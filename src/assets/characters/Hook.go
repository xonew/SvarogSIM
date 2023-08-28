package characters

import . "SvarogSIM/src/classes"

type Hook struct {
	Character
	Enhanced bool
}

func MakeHook() *Hook {
	return &Hook{
		Character: MakeCharacter("hook", "fire", "destruction",
			0, 80, 120,
			1340, 617, 352, 94),
		Enhanced: false,
	}
}

func (h *Hook) Init(left Ally, right Ally, heapify func()) {
	h.CritDmg += 0.053 + 0.08
	h.Hp.Percent += 0.04 + 0.08 + 0.06
	h.Atk.Percent += 0.04 + 0.04 + 0.06 + 0.06 + 0.08
	h.ActionValue = int(15000 / h.Spd.GetStat())
	h.CurrHp = int(h.Hp.GetStat())
	h.Heapify = heapify
	h.Left = left
	h.Right = right
}

func (h *Hook) Act(allies []Ally, enemies []Enemy, skillPoints int) int {
	pointGain := 0
	if skillPoints > 0 {
		pointGain += h.skill(Target(enemies))
	} else {
		pointGain += h.basicAttack(Target(enemies))
	}
	return skillPoints + pointGain
}

func (h *Hook) basicAttack(target Actor) int {
	h.Event("basicStart")
	hit := h.MakeAttack("Basic Attack", target.GetName(), "fire", "skill", map[Stat]float64{
		h.Atk: 0.80,
	})
	target.TakeDamage(hit)
	h.HitEvent("outStart", hit)
	h.Event("basicEnd")
	h.RegenEnergy(20)
	return 1
}

func (h *Hook) skill(target Actor) int {
	var hit *Attack

	if h.Enhanced {
		h.DmgBonus["skill"] += 0.20
		hit = h.MakeAttack("Enhanced Skill", target.GetName(), "fire", "skill", map[Stat]float64{
			h.Atk: 2.80,
		})
		h.DmgBonus["skill"] -= 0.20
	} else {
		hit = h.MakeAttack("Skill", target.GetName(), "fire", "skill", map[Stat]float64{
			h.Atk: 2.40,
		})
	}
	target.TakeDamage(hit)

	if h.Enhanced {
		if target.GetLeft() != nil {
			blast := h.MakeAttack("Enhanced Skill Blast", target.GetName(), "fire", "skill", map[Stat]float64{
				h.Atk: 0.80,
			})
			target.GetRight().TakeDamage(blast)
		}
		if target.GetRight != nil {
			blast := h.MakeAttack("Enhanced Skill Blast", target.GetName(), "fire", "skill", map[Stat]float64{
				h.Atk: 0.80,
			})
			target.GetRight().TakeDamage(blast)
		}

		h.Enhanced = false
	}

	h.talent(target)

	//apply burn

	h.RegenEnergy(30)
	return -1
}

func (h *Hook) talent(target Actor) {
	//TODO: if target.HasDebuff("burn") || target.HasDebuff("breakBurn") {
	if true {
		hit := h.MakeAttack("Talent", target.GetName(), "fire", "additional", map[Stat]float64{
			h.Atk: 1.0,
		})
		target.TakeDamage(hit)
	}
	//TODO: h.RestorePercentHp(0.05)
}

func (h *Hook) GetCharacter() *Character {
	return &h.Character
}

// burn has a creature source and a creature target
// every time proc() is called, the source will make a burn attack to the target
