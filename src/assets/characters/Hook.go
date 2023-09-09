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
	// TODO: call a universal entity init function here
	h.ActionValue = int(15000 / h.Spd.GetStat())
	h.CurrHp = int(h.Hp.GetStat())
	h.Heapify = heapify
	h.Left = left
	h.Right = right
}

func (h *Hook) Act() {
	pointGain := 0
	if h.Battle.SkillPoints.Get() > 0 {
		pointGain += h.skill(Target(h.Battle.Enemies))
	} else {
		pointGain += h.basicAttack(Target(h.Battle.Enemies))
	}
	h.Battle.SkillPoints.Add(pointGain)
	h.ModifyActionValue(h.GetBaseActionValue())
}

func (h *Hook) basicAttack(target Enemy) int {
	h.Event("basicStart")
	hit := h.MakeAttack("Basic Attack", target.GetName(), "fire", "skill", map[Stat]float64{
		h.Atk: 0.80,
	})
	h.HitEvent("outStart", hit)
	target.TakeDamage(hit)
	h.HitEvent("outEnd", hit)
	h.Event("basicEnd")
	h.RegenEnergy(20)
	return 1
}

func (h *Hook) skill(target Enemy) int {
	h.Event("skillStart")
	var hit *Attack

	if h.Enhanced {
		hit = h.MakeAttack("Enhanced Skill", target.GetName(), "fire", "skill", map[Stat]float64{
			h.Atk: 2.80,
		})
		hit.DamageBonus += 0.2
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
		if target.GetRight() != nil {
			blast := h.MakeAttack("Enhanced Skill Blast", target.GetName(), "fire", "skill", map[Stat]float64{
				h.Atk: 0.80,
			})
			target.GetRight().TakeDamage(blast)
		}

		h.Enhanced = false
	}

	h.talent(target)

	//apply burn
	h.Event("skillEnd")
	h.RegenEnergy(30)
	return -1
}

func (h *Hook) talent(target Enemy) {
	if target.HasDebuff("burn") {
		hit := h.MakeAttack("Talent", target.GetName(), "fire", "additional", map[Stat]float64{
			h.Atk: 1.0,
		})
		target.TakeDamage(hit)
	}
	h.RestoreHp(int(h.Hp.GetStat() * 0.05))
}

// burn has a creature source and a creature target
// every time proc() is called, the source will make a burn attack to the target
