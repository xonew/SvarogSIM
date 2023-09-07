package classes

type Ally interface {
	Creature
	Init(left Ally, right Ally, heapify func())
	GetAggro() float64
	GetCharacter() *Character
	GetLeft() Ally
	GetRight() Ally
}

type Character struct {
	Entity
	Element string
	Path    string
	Eidolon int

	CurrEnergy  float64
	MaxEnergy   float64
	EnergyRegen float64

	Aggro Stat // Aggro is standard for characters, is -2 for elites, -3 for adds, and is
}

func MakeCharacter(name string, element string, path string,
	eidolon int, level int, maxEnergy float64,
	baseHp float64, baseAtk float64, baseDef float64, baseSpd float64) Character {
	aggro := 0
	switch path {
	case "hunt":
		aggro = 75
	case "erudition":
		aggro = 75
	case "harmony":
		aggro = 100
	case "abundance":
		aggro = 100
	case "nihility":
		aggro = 100
	case "destruction":
		aggro = 125
	case "preservation":
		aggro = 150
	}
	return Character{
		Entity: Entity{
			Name:               name,
			Level:              level,
			CurrHp:             0,
			Hp:                 Stat{Base: baseHp, Percent: 0, Flat: 0},
			Atk:                Stat{Base: baseAtk, Percent: 0, Flat: 0},
			Def:                Stat{Base: baseDef, Percent: 0, Flat: 0},
			Spd:                Stat{Base: baseSpd, Percent: 0, Flat: 0},
			EffectHitRate:      0,
			EffectResist:       0,
			CrowdControlResist: 0,
			DmgBonus:           make(map[string]float64),
			Buffs:              make(map[string]map[string]Effect),
			Debuffs:            make(map[string]map[string]Effect),
			ResPen:             make(map[string]float64),
			Res:                make(map[string]float64),
			DefPen:             0,
			ActionValue:        0,
			CritDmg:            0.5,
			CritRate:           0.05,

			DamageOutLog: make(map[string][]*Attack),
			DamageInLog:  make(map[string][]*Attack),
		},
		Element: element,
		Path:    path,
		Eidolon: eidolon,

		CurrEnergy:  0,
		MaxEnergy:   maxEnergy,
		EnergyRegen: 0,

		Aggro: Stat{float64(aggro), 0, 0},
	}
}

// RegenEnergy regenerates energy for the character
func (c *Character) RegenEnergy(energy float64) {
	c.CurrEnergy += energy * (1 + c.EnergyRegen)
	if c.CurrEnergy > c.MaxEnergy {
		c.CurrEnergy = c.MaxEnergy
	}
}

func (c *Character) GetAggro() float64 {
	return c.Aggro.GetStat()
}

func Target(enemies []Enemy) Enemy {
	for i := range enemies {
		if enemies[i].GetAggro() == "boss" {
			return enemies[i]
		}
	}
	for i := range enemies {
		if enemies[i].GetAggro() == "elite" {
			return enemies[i]
		}
	}
	return enemies[0]
}

func (c *Character) GetCharacter() *Character {
	return c
} // CHECK IF THIS ACTUALLY WORKS
