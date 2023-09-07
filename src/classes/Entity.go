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
	Act()
	GetActionValue() int // returns action value

}
type Creature interface {
	Actor
	ApplyBuff(buff Effect) bool
	ApplyDebuff(debuff Effect) bool
	//Dispel() bool // gets rid of buff
	//Cleanse() bool // gets rid of debuff
	GetName() string
	TakeDamage(attack *Attack)
	RestoreHp(amount int) int
	AddToRight(actor Actor)
	AddToLeft(actor Actor)
	HasDebuff(s string) bool
	InitBattle(battle Combat) // adds battle to creature
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
	Battle       *Combat
}

func (e *Entity) InitBattle(battle *Combat) {
	e.Battle = battle
}

func (e *Entity) GetName() string {
	return e.Name
}

func (e *Entity) RestoreHp(amount int) int {
	postBonusHealing := int(float64(amount) * (1 + e.IncomingHealingBonus))
	if e.CurrHp+postBonusHealing > int(e.Hp.GetStat()) {
		postBonusHealing = int(e.Hp.GetStat()) - e.CurrHp
	}
	e.CurrHp += postBonusHealing
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

func (e *Entity) GetActionValue() int {
	return e.ActionValue
}

func (e *Entity) GetDebuffs() *map[string]map[string]Effect {
	return &e.Debuffs
}

func (e *Entity) HasDebuff(s string) bool {
	for _, ids := range e.Debuffs {
		for _, source := range ids {
			if source.GetId() == s {
				return true
			}
		}
	}
	return false
}

func (e *Entity) GetBaseActionValue() int {
	return int(10000 / e.Spd.GetStat())
}

func (e *Entity) ActionAdvance(value float64) {
	e.ActionValue = int(math.Max(0, float64(e.ActionValue)*(1-value)))
}

// ApplyBuff applies a buff to the creature, overwriting buffs of the same type and wielder
func (e *Entity) ApplyBuff(buff Effect) {
	e.Buffs[buff.GetId()][buff.GetSource()] = buff
	buff.Apply()
	//TODO: debuffs, cleanse, dispel, stacking debuffs
}

// ApplyDebuff applies a debuff to the creature, overwriting debuffs of the same type and wielder
func (e *Entity) ApplyDebuff(debuff Effect) bool {
	if debuff.GetEffectiveHitRate()*(1-e.EffectResist) < rand.Float64() { //TODO: make this into a helper function
		e.Debuffs[debuff.GetId()][debuff.GetSource()] = debuff
		debuff.Apply()
		e.Event("debuffApplied")
		return true
	} else {
		e.Event("debuffResisted")
		return false
	}
}

// RollCrit returns the crit multiplier if the attack crits, or 1 if it doesn't
func (e *Entity) RollCrit() float64 {
	if rand.Float64() <= e.CritRate {
		return 1 + e.CritDmg
	} else {
		return 1
	}
}

// LogDamageOut logs an attack to the creature's damage out log
func (e *Entity) LogDamageOut(attack *Attack) {
	e.DamageOutLog[attack.Target] = append(e.DamageOutLog[attack.Target], attack)
}

// LogDamageIn logs an attack to the creature's damage in log
func (e *Entity) LogDamageIn(attack *Attack) {
	e.DamageInLog[attack.Attacker] = append(e.DamageInLog[attack.Attacker], attack)
}

func (e *Entity) AddToRight(actor Actor) {
	e.Right = actor
}

func (e *Entity) AddToLeft(actor Actor) {
	e.Left = actor
}

func (e *Entity) GetRight() Actor {
	return e.Right
}

func (e *Entity) GetLeft() Actor {
	return e.Left
}

func (e *Entity) AddListener(function func(), event string, id string) {
	e.Listeners[event][id] = function
}

func (e *Entity) AddHitListener(function func(*Attack), event string, id string) {
	e.HitListeners[event][id] = function
}

func (e *Entity) Event(event string) {
	for _, function := range e.Listeners[event] {
		function()
	}
}
func (e *Entity) HitEvent(event string, attack *Attack) {
	for _, function := range e.HitListeners[event] {
		function(attack)
	}
}

func (e *Entity) RemoveListener(event string, id string) {
	delete(e.Listeners[event], id)
}

func (e *Entity) RemoveHitListener(event string, id string) {
	delete(e.HitListeners[event], id)
}
