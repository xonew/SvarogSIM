package classes

type SkillPoint struct {
	pool int
	max  int
}

type SkillPoints interface {
	Get() int
	Add(int)
	Reduce(int)
}

func (s *SkillPoint) Get() int {
	return s.pool
}

func (s *SkillPoint) Add(i int) {
	if s.pool+i <= 5 {
		s.pool += i
	}
}

func (s *SkillPoint) Reduce(i int) {
	if s.pool-i >= 0 {
		s.pool -= i
	}
}
