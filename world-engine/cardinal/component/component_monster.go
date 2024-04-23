package component

type MonsterType int

const NumMonsterTypes = 3

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
