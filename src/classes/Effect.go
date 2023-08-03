package classes

type Observer interface {
	Update()
}

type Effect interface {
	Apply()
	Proc()
	Observer
	GetId() string
	GetSource() string
}

type Buff struct {
	Id         string
	Duration   int
	Strength   float64
	Value      float64 // optional
	SourceName string
}

type Debuff struct {
	Buff
	IsControlEffect bool
	Removable       bool
}

type Dot struct {
	Debuff
	Stacks  int
	Source  *Creature
	Holder  *Creature
	IsBreak bool
}

func (b *Buff) GetId() string {
	return b.Id
}

func (b *Buff) GetSource() string {
	return b.SourceName
}

func MakeBurn(source *Creature, Holder *Creature, strength float64, duration int) *Dot {
	return &Dot{
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

func (d *Dot) Apply() {
	if d.Holder.ApplyDebuff(d) {
		d.Holder.AddListener(d.Update, "turnEnd", d.Id)
		d.Holder.AddListener(d.Proc, "turnStart", d.Id)
	}
}

func (d *Dot) Proc() {
	attack := &Attack{
		Name:          d.Source.Name,
		Attacker:      d.Source.Name,
		Target:        d.Holder.Name,
		Element:       "fire", // todo: make this dynamic
		AttackerLevel: d.Source.Level,
		PreMitDamage:  d.Source.Atk.GetStat() * d.Strength,
		DefPen:        0,
	}
	d.Holder.TakeDamage(attack)
	d.Source.LogDamageOut(attack)
}

func (d *Dot) Update() {
	d.Duration--
}
