package component

import "pkg.world.dev/world-engine/cardinal"

const Ability7ID = 7

type Ability7 struct{}

var _ Ability = &Ability7{}

func (Ability7) GetAbilityID() int {
	return Ability7ID
}

func (a Ability7) Resolve(
	world cardinal.WorldContext,
	spellPosition *Position,
	direction Direction,
	executeUpdates bool,
	eventLogList *[]GameEventLog,
) (reveal bool, err error) {
	panic("implement me")
	return false, nil
}
