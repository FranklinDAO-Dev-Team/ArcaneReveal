package component

type MonsterType int

const NumMonsterTypes = 2

const (
	LIGHT MonsterType = iota
	MEDIUM
	// HEAVY
	// XL
)

type Monster struct {
	Type MonsterType
}

func (Monster) Name() string {
	return "Monster"
}
