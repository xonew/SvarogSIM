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

func (c *Entity) MakeAttack(name string, target string,
	element string, attackType string,
	scaling map[Stat]float64) *Attack {
	preMitDamage := 0.0
	for stat, mod := range scaling {
		preMitDamage += mod * stat.GetStat()
	}
	attack := &Attack{
		Name:          name,
		Attacker:      c.Name,
		Target:        target,
		Element:       element,
		AttackerLevel: c.Level,
		Scaling:       scaling,
		DamageBonus:   1 + c.DmgBonus[attackType] + c.DmgBonus["all"] + c.DmgBonus[element],
		CritRate:      c.CritRate,
		CritDmg:       c.CritDmg,
		DefPen:        c.DefPen,
		ResPen:        c.ResPen[element] + c.ResPen["all"] + c.ResPen[attackType],
		PostMitDamage: 0,
	}
	c.LogDamageOut(attack)
	return attack
}

// TakeDamage receives an attack, and returns the final damage taken
func (c *Entity) TakeDamage(attack *Attack) {
	c.HitEvent("inStart", attack)
	defMultiplier := 0.01 * (100 - (c.Def.GetStat() / (c.Def.GetReducedStat(attack.DefPen) + 200 + 10*float64(attack.AttackerLevel))))
	resMultiplier := 0.01 * (100 - (c.Res[attack.Element] - attack.ResPen))
	dmgTakenMultiplier := 1.0              //TODO: Implement elemental dmg taken
	universalDmgReductionMultiplier := 0.9 //TODO: Implement universal dmg reduction
	preMitDamage := attack.FlatDamage
	for stat, mod := range attack.Scaling {
		preMitDamage += mod * stat.GetStat()
	}
	postMitDamage := int(preMitDamage * defMultiplier * resMultiplier * dmgTakenMultiplier * universalDmgReductionMultiplier)
	if c.CurrHp-postMitDamage <= 0 {
		c.CurrHp = 0
		c.Event("death")
	} else {
		c.CurrHp -= postMitDamage
	}
	attack.PostMitDamage = postMitDamage
	c.LogDamageIn(attack)
	c.HitEvent("inEnd", attack)
}
