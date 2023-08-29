package classes

type SkillPoint struct {
	pool int
	max  int
}

type SkillPoints interface {
	get() int
	add(int)
	reduce(int)
}

func (s *SkillPoint) get() int {
	return s.pool
}

func (s *SkillPoint) add(i int) {
	if s.pool+i <= 5 {
		s.pool += i
	}
}

func (s *SkillPoint) reduce(i int) {
	if s.pool-i >= 0 {
		s.pool -= i
	}
}
