package classes

import "math/rand"

type Enemy interface {
	Actor
	Init(left Enemy, right Enemy, heapify *func())
	IsWeakTo(element string) bool
	GetAggro() string
}

type Mob struct {
	Entity
	EnemyType       string
	MaxToughness    float64
	CurrToughness   float64
	Weaknesses      map[string]bool
	ToughnessBroken bool
}

func MakeMob(name string, level int, enemyType string, toughness float64,
	weaknesses map[string]bool, baseHp float64, baseAtk float64, baseDef float64, baseSpd float64) *Mob {
	return &Mob{
		Entity: Entity{
			Name:               name,
			Level:              level,
			CurrHp:             int(baseHp),
			Hp:                 Stat{Base: baseHp, Percent: 0, Flat: 0},
			Atk:                Stat{Base: baseAtk, Percent: 0, Flat: 0},
			Def:                Stat{Base: baseDef, Percent: 0, Flat: 0},
			Spd:                Stat{Base: baseSpd, Percent: 0, Flat: 0},
			EffectHitRate:      0,
			EffectResist:       0,
			CrowdControlResist: 0,
			CritRate:           0,
			CritDmg:            0,
			DmgBonus:           make(map[string]float64),
			Buffs:              make(map[string]map[string]Effect),
			Debuffs:            make(map[string]map[string]Effect),
			ResPen:             make(map[string]float64),
			Res:                make(map[string]float64),
		},
		EnemyType:       enemyType,
		MaxToughness:    toughness,
		CurrToughness:   toughness,
		Weaknesses:      weaknesses,
		ToughnessBroken: false,
	}
}

func (m *Mob) MakeAttack(ally Ally) *Attack {
	//TODO: do actual attacks
	return &Attack{
		Name:          m.Name,
		Attacker:      m.Name,
		Target:        ally.GetName(),
		Element:       "none",
		AttackerLevel: m.Level,
		PreMitDamage:  0,
		DefPen:        0,
		ResPen:        0,
		PostMitDamage: 0,
	}
}

func (m *Mob) Act(allies []Ally, enemies []Enemy, skillPoints int) int {
	if m.CurrToughness <= 0 {
		m.ToughnessBroken = false
	}
	target := ChooseTarget(allies)
	target.TakeDamage(m.MakeAttack(target))
	return skillPoints
}

func ChooseTarget(allies []Ally) Ally {
	var totalAggro float64
	for _, ally := range allies {
		totalAggro += ally.GetAggro()
	}
	aggro := rand.Float64() * totalAggro
	for _, ally := range allies {
		aggro -= ally.GetAggro()
		if aggro <= 0 {
			return ally
		}
	}
	return allies[0]
}

func (m *Mob) Init(left Enemy, right Enemy, heapify func()) {
	m.CurrHp = int(m.Hp.GetStat())
	m.CurrToughness = m.MaxToughness
	m.ActionValue = int(15000 / m.Spd.GetStat())
	m.Heapify = heapify
	m.Left = left
	m.Right = right
}
func (m *Mob) IsWeakTo(element string) bool {
	return m.Weaknesses[element]
}

func MakeWeakCocolia() *Mob {
	cocolia := MakeMob("Weak Cocolia", 100, "boss", 360, map[string]bool{
		"physical":  true,
		"fire":      true,
		"ice":       true,
		"lightning": true,
		"wind":      true,
		"quantum":   true,
		"imaginary": true,
	}, 592150, 773, 1200, 158)
	cocolia.EffectHitRate = 0.4
	cocolia.EffectResist = 0.4
	cocolia.CrowdControlResist = 0.5
	return cocolia
}

func (m *Mob) GetAggro() string {
	return m.EnemyType
}
