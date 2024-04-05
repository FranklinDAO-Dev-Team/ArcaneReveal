package component

import "pkg.world.dev/world-engine/cardinal"

const Ability4ID = 4

type Ability4 struct{}

var _ Ability = &Ability4{}

func (Ability4) GetAbilityID() int {
	return Ability4ID
}

func (a Ability4) Resolve(
	world cardinal.WorldContext,
	spellPosition *Position,
	direction Direction,
	executeUpdates bool,
	eventLogList *[]GameEventLog,
) (reveal bool, err error) {
	// panic("implement me a4")
	return false, nil
}
