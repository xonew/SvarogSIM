package classes

import (
	"math"
	"math/rand"
)

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
	GetActionValue() int // returns action value

}
type Creature interface {
	Actor
	ApplyBuff(buff Effect) bool
	ApplyDebuff(debuff Effect) bool
	//Dispel() bool // gets rid of buff
	//Cleanse() bool // gets rid of debuff
	GetName() string
	GetLeft() Actor
	GetRight() Actor
	TakeDamage(attack *Attack)
	RestoreHp(amount int) int
	AddToRight(actor Actor)
	AddToLeft(actor Actor)
	HasDebuff(s string) bool
}

type Entity struct {
	Name  string
	Level int

	CurrHp int
	Hp     Stat
	Atk    Stat
	Def    Stat
	Spd    Stat

	IncomingHealingBonus float64
	OutgoingHealingBonus float64

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

	Listeners map[string]map[string]func() //map[Event]map[Id]func()  why do i need id?
	//Events: turnStart, turnEnd, death, basicStart, basicEnd, skillStart, skillEnd, outStart, outEnd, inStart, inEnd
	HitListeners map[string]map[string]func(*Attack)
	//Events: outStart, outEnd, inStart, inEnd

	DamageOutLog map[string][]*Attack
	DamageInLog  map[string][]*Attack
	Left         Actor
	Right        Actor
	Heapify      func()
}

func (c *Entity) GetName() string {
	return c.Name
}

func (c *Entity) RestoreHp(amount int) int {
	postBonusHealing := int(float64(amount) * (1 + c.IncomingHealingBonus))
	if c.CurrHp+postBonusHealing > int(c.Hp.GetStat()) {
		postBonusHealing = int(c.Hp.GetStat()) - c.CurrHp
	}
	c.CurrHp += postBonusHealing
	return postBonusHealing
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

func (c *Entity) GetActionValue() int {
	return c.ActionValue
}

func (c *Entity) GetDebuffs() *map[string]map[string]Effect {
	return &c.Debuffs
}

func (c *Entity) HasDebuff(s string) bool {
	for _, ids := range c.Debuffs {
		for _, source := range ids {
			if source.GetId() == s {
				return true
			}
		}
	}
	return false
}

func (c *Entity) GetBaseActionValue() int {
	return int(10000 / c.Spd.GetStat())
}

func (c *Entity) ActionAdvance(value float64) {
	c.ActionValue = int(math.Max(0, float64(c.ActionValue)*(1-value)))
}

// ApplyBuff applies a buff to the creature, overwriting buffs of the same type and wielder
func (c *Entity) ApplyBuff(buff Effect) {
	c.Buffs[buff.GetId()][buff.GetSource()] = buff
	buff.Apply()
	//TODO: debuffs, cleanse, dispel, stacking debuffs
}

// ApplyDebuff applies a debuff to the creature, overwriting debuffs of the same type and wielder
func (c *Entity) ApplyDebuff(debuff Effect) bool {
	if debuff.GetEffectiveHitRate()*(1-c.EffectResist) < rand.Float64() { //TODO: make this into a helper function
		c.Debuffs[debuff.GetId()][debuff.GetSource()] = debuff
		debuff.Apply()
		c.Event("debuffApplied")
		return true
	} else {
		c.Event("debuffResisted")
		return false
	}
}

// RollCrit returns the crit multiplier if the attack crits, or 1 if it doesn't
func (c *Entity) RollCrit() float64 {
	if rand.Float64() <= c.CritRate {
		return 1 + c.CritDmg
	} else {
		return 1
	}
}

// LogDamageOut logs an attack to the creature's damage out log
func (c *Entity) LogDamageOut(attack *Attack) {
	c.DamageOutLog[attack.Target] = append(c.DamageOutLog[attack.Target], attack)
}

// LogDamageIn logs an attack to the creature's damage in log
func (c *Entity) LogDamageIn(attack *Attack) {
	c.DamageInLog[attack.Attacker] = append(c.DamageInLog[attack.Attacker], attack)
}

func (c *Entity) AddToRight(actor Actor) {
	c.Right = actor
}

func (c *Entity) AddToLeft(actor Actor) {
	c.Left = actor
}

func (c *Entity) GetRight() Actor {
	return c.Right
}

func (c *Entity) GetLeft() Actor {
	return c.Left
}

func (c *Entity) AddListener(function func(), event string, id string) {
	c.Listeners[event][id] = function
}

func (c *Entity) AddHitListener(function func(*Attack), event string, id string) {
	c.HitListeners[event][id] = function
}

func (c *Entity) Event(event string) {
	for _, function := range c.Listeners[event] {
		function()
	}
}
func (c *Entity) HitEvent(event string, attack *Attack) {
	for _, function := range c.HitListeners[event] {
		function(attack)
	}
}

func (c *Entity) RemoveListener(event string, id string) {
	delete(c.Listeners[event], id)
}

func (c *Entity) RemoveHitListener(event string, id string) {
	delete(c.HitListeners[event], id)
}
