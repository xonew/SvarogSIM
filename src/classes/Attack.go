package classes

type Attack struct {
	Name          string
	Attacker      string
	Target        string
	Element       string
	AttackType    string
	AttackerLevel int
	Scaling       map[Stat]float64
	FlatDamage    float64
	CritRate      float64
	CritDmg       float64
	DamageBonus   float64
	DefPen        float64
	ResPen        float64
	PostMitDamage int
}

func (e *Entity) MakeAttack(name string, target string,
	element string, attackType string,
	scaling map[Stat]float64) *Attack {
	preMitDamage := 0.0
	for stat, mod := range scaling {
		preMitDamage += mod * stat.GetStat()
	}
	attack := &Attack{
		Name:          name,
		Attacker:      e.Name,
		Target:        target,
		Element:       element,
		AttackerLevel: e.Level,
		Scaling:       scaling,
		DamageBonus:   1 + e.DmgBonus[attackType] + e.DmgBonus["all"] + e.DmgBonus[element],
		CritRate:      e.CritRate,
		CritDmg:       e.CritDmg,
		DefPen:        e.DefPen,
		ResPen:        e.ResPen[element] + e.ResPen["all"] + e.ResPen[attackType],
		PostMitDamage: 0,
	}
	e.LogDamageOut(attack)
	return attack
}

// TakeDamage receives an attack, and returns the final damage taken
func (e *Entity) TakeDamage(attack *Attack) {
	e.HitEvent("inStart", attack)
	defMultiplier := 0.01 * (100 - (e.Def.GetStat() / (e.Def.GetReducedStat(attack.DefPen) + 200 + 10*float64(attack.AttackerLevel))))
	resMultiplier := 0.01 * (100 - (e.Res[attack.Element] - attack.ResPen))
	dmgTakenMultiplier := 1.0              //TODO: Implement elemental dmg taken
	universalDmgReductionMultiplier := 0.9 //TODO: Implement universal dmg reduction
	preMitDamage := attack.FlatDamage
	for stat, mod := range attack.Scaling {
		preMitDamage += mod * stat.GetStat()
	}
	postMitDamage := int(preMitDamage * defMultiplier * resMultiplier * dmgTakenMultiplier * universalDmgReductionMultiplier)
	if e.CurrHp-postMitDamage <= 0 {
		e.CurrHp = 0
		e.Event("death")
	} else {
		e.CurrHp -= postMitDamage
	}
	attack.PostMitDamage = postMitDamage
	e.LogDamageIn(attack)
	e.HitEvent("inEnd", attack)
}
