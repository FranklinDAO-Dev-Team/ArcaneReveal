package component

import "pkg.world.dev/world-engine/cardinal"

const Ability6ID = 6

type Ability6 struct{}

var _ Ability = &Ability6{}

func (Ability6) GetAbilityID() int {
	return Ability6ID
}

func (a Ability6) Resolve(
	world cardinal.WorldContext,
	spellPosition *Position,
	direction Direction,
	executeUpdates bool,
	eventLogList *[]GameEventLog,
) (reveal bool, err error) {
	panic("implement me")
	return false, nil
}
