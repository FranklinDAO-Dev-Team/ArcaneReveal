package component

import "pkg.world.dev/world-engine/cardinal"

const Ability8ID = 8

type Ability8 struct{}

var _ Ability = &Ability8{}

func (Ability8) GetAbilityID() int {
	return Ability8ID
}

func (a Ability8) Resolve(
	world cardinal.WorldContext,
	spellPosition *Position,
	direction Direction,
	executeUpdates bool,
	eventLogList *[]GameEventLog,
) (reveal bool, err error) {
	panic("implement me")
	return false, nil
}
