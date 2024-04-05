package component

type MonsterType int

const (
	LIGHT MonsterType = iota
	MEDIUM
	HEAVY
)

type Monster struct {
	Type MonsterType
}

func (Monster) Name() string {
	return "Monster"
}
