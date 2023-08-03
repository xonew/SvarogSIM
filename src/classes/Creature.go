package classes

import "math/rand"

type Stat struct {
	Base    float64
	Percent float64
	Flat    float64
}

func (s *Stat) GetStat() float64 {
	return s.Base*(1+s.Percent) + s.Flat
}

func (s *Stat) GetReducedStat(reduction float64) float64 {
	return s.Base*(1+s.Percent-reduction) + s.Flat
}

type Actor interface {
	Act(allies []Ally, enemies []Enemy, skillPoints int) int
	ApplyBuff(buff Effect) bool
	ApplyDebuff(debuff Effect) bool
	//Dispel() bool
	//Cleanse() bool
	GetName() string
	GetLeft() Actor
	GetRight() Actor
	TakeDamage(attack *Attack)
	GetActionValue() int // returns action value
	//RestoreHp(amount int) int
	AddToRight(actor Actor)
	AddToLeft(actor Actor)
}

type Creature struct {
	Name  string
	Level int

	CurrHp int
	Hp     Stat
	Atk    Stat
	Def    Stat
	Spd    Stat

	EffectHitRate      float64
	EffectResist       float64
	CrowdControlResist float64

	DmgBonus map[string]float64
	Buffs    map[string]map[string]Effect // map[Id]map[Source]Effect
	Debuffs  map[string]map[string]Effect

	Res         map[string]float64
	ResPen      map[string]float64
	DefPen      float64
	ActionValue int

	CritRate float64
	CritDmg  float64

	Listeners    map[string]map[string]func() //map[Event]map[Id]func()
	HitListeners map[string]map[string]func(*Attack)
	DamageOutLog map[string][]*Attack
	DamageInLog  map[string][]*Attack
	Left         Actor
	Right        Actor
	Heapify      func()
}

type Attack struct {
	Name          string
	Attacker      string
	Target        string
	Element       string
	AttackType    string
	AttackerLevel int
	PreMitDamage  float64
	DefPen        float64
	ResPen        float64
	PostMitDamage int
}

func (c *Creature) GetName() string {
	return c.Name
}

/*
Outgoing DMG = Base DMG * DMG% Multiplier * DEF Multiplier * RES Multiplier * DMG Taken Multiplier * Universal DMG Reduction Multiplier * Weaken Multiplier

Base DMG = (Skill Multiplier + Extra Multiplier) * Scaling Attribute + Extra DMG
DMG% Multiplier = 100% + Elemental DMG% + All Type DMG% + DoT DMG% + Other DMG%
DEF Multiplier = 100% - [DEF / (DEF + 200 + 10 * Attacker Level)]
	DEF = Base DEF * (100% + DEF% - (DEF Reduction + DEF Ignore)) + Flat DEF
RES Multiplier = 100% - (RES% - RES PEN%)
DMG Taken Multiplier = 100% + Elemental DMG Taken% + All Type DMG Taken%
Universal DMG Reduction Multiplier = 100% * (1 - DMG Reduction_1) * (1 - DMG Reduction_2) * ...
	When an enemy has Toughness, they have 10% Universal DMG Reduction, which is reduced to 0% when broken. Note this multiplier stacks multiplicative with other sources.
*/

func (c *Creature) GetActionValue() int {
	return c.ActionValue
}

// ApplyBuff applies a buff to the creature, overwriting buffs of the same type and wielder
func (c *Creature) ApplyBuff(buff Effect) {
	c.Buffs[buff.GetId()][buff.GetSource()] = buff
	buff.Apply()
	//TODO: debuffs, cleanse, dispel, stacking debuffs
}

// ApplyDebuff applies a debuff to the creature, overwriting debuffs of the same type and wielder
func (c *Creature) ApplyDebuff(debuff Effect) bool {
	c.Debuffs[debuff.GetId()][debuff.GetSource()] = debuff
	debuff.Apply()
}

// TakeDamage receives an attack, and returns the final damage taken
func (c *Creature) TakeDamage(attack *Attack) {
	defMultiplier := 0.01 * (100 - (c.Def.GetStat() / (c.Def.GetReducedStat(attack.DefPen) + 200 + 10*float64(attack.AttackerLevel))))
	resMultiplier := 0.01 * (100 - (c.Res[attack.Element] - attack.ResPen))
	dmgTakenMultiplier := 1.0              //TODO: Implement elemental dmg taken
	universalDmgReductionMultiplier := 0.9 //TODO: Implement universal dmg reduction
	postMitDamage := int(attack.PreMitDamage * defMultiplier * resMultiplier * dmgTakenMultiplier * universalDmgReductionMultiplier)
	if c.CurrHp-postMitDamage <= 0 {
		c.CurrHp = 0
		c.Event("death")
	} else {
		c.CurrHp -= postMitDamage
	}
	attack.PostMitDamage = postMitDamage
	c.LogDamageIn(attack)
}

func (c *Creature) MakeAttack(name string, target string,
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
		PreMitDamage: preMitDamage *
			(1 + c.DmgBonus[attackType] + c.DmgBonus["all"] + c.DmgBonus[element]) *
			c.RollCrit(),
		DefPen:        c.DefPen,
		ResPen:        c.ResPen[element] + c.ResPen["all"] + c.ResPen[attackType],
		PostMitDamage: 0,
	}
	c.LogDamageOut(attack)
	return attack
}

// RollCrit returns the crit multiplier if the attack crits, or 1 if it doesn't
func (c *Creature) RollCrit() float64 {
	if rand.Float64() <= c.CritRate {
		return 1 + c.CritDmg
	} else {
		return 1
	}
}

// LogDamageOut logs an attack to the creature's damage out log
func (c *Creature) LogDamageOut(attack *Attack) {
	c.DamageOutLog[attack.Target] = append(c.DamageOutLog[attack.Target], attack)
}

// LogDamageIn logs an attack to the creature's damage in log
func (c *Creature) LogDamageIn(attack *Attack) {
	c.DamageInLog[attack.Attacker] = append(c.DamageInLog[attack.Attacker], attack)
}

func (c *Creature) AddToRight(actor Actor) {
	c.Right = actor
}

func (c *Creature) AddToLeft(actor Actor) {
	c.Left = actor
}

func (c *Creature) GetRight() Actor {
	return c.Right
}

func (c *Creature) GetLeft() Actor {
	return c.Left
}

func (c *Creature) AddListener(function func(), event string, id string) {
	c.Listeners[event][id] = function
}

func (c *Creature) AddHitListener(function func(*Attack), event string, id string) {
	c.HitListeners[event][id] = function
}

func (c *Creature) Event(event string) {
	for _, function := range c.Listeners[event] {
		function()
	}
}
func (c *Creature) HitEvent(event string, attack *Attack) {
	for _, function := range c.HitListeners[event] {
		function(attack)
	}
}

func (c *Creature) RemoveListener(event string, id string) {
	delete(c.Listeners[event], id)
}

func (c *Creature) RemoveHitListener(event string, id string) {
	delete(c.HitListeners[event], id)
}
