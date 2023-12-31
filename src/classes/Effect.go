package classes

type Observer interface {
	Update()
}

type Effect interface {
	Apply()
	Observer
	GetId() string
	GetSource() string
	GetEffectiveHitRate() float64
	IsStackable() bool
}

type Dot interface {
	Effect
	ProcDot()
}

type Buff struct {
	Id         string
	Duration   int
	Strength   float64
	Value      float64 // optional
	SourceName string
	Stackable  bool
}

type Debuff struct {
	Buff
	IsControlEffect bool
	Removable       bool
	BaseHitRate     float64
	EffectHitRate   float64
}

type DotDebuff struct {
	Debuff
	Stacks  int
	Source  *Entity
	Holder  *Entity
	IsBreak bool
	Scaling map[*Stat]float64
}

func (b *Buff) IsStackable() bool {
	return b.Stackable
}

func (b *Buff) GetId() string {
	return b.Id
}

func (b *Buff) GetSource() string {
	return b.SourceName
}

func MakeBurn(source *Entity, Holder *Entity, strength float64, duration int) *DotDebuff {
	return &DotDebuff{
		Debuff: Debuff{
			Buff: Buff{
				Id:       "burn",
				Duration: duration,
				Strength: strength,
				Value:    0,
			},
			IsControlEffect: false,
			Removable:       true,
		},
		Stacks: 1,
		Source: source,
		Holder: Holder,
	}
}

func (d *DotDebuff) Apply() {
	if d.Holder.ApplyDebuff(d) {
		d.Holder.AddListener(d.Update, "turnStart", d.Id)
	}
}

func (d *DotDebuff) ProcDot() {
	scaling := make(map[Stat]float64)
	for keyPointer, value := range d.Scaling {
		// Dereference the pointer to "snapshot" the stat
		key := *keyPointer
		scaling[key] = value
	}
	attack := &Attack{
		Name:          d.Source.Name,
		Attacker:      d.Source.Name,
		Target:        d.Holder.Name,
		Element:       d.Id,
		AttackerLevel: d.Source.Level,
		Scaling:       scaling,
		DefPen:        0,
	}
	d.Holder.TakeDamage(attack)
	d.Source.LogDamageOut(attack)
}

func (d *DotDebuff) Update() {
	d.ProcDot()
	d.Duration--
}

func (d *Debuff) GetEffectiveHitRate() float64 {
	return d.BaseHitRate * (1 + d.EffectHitRate)
}
