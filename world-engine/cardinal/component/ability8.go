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
	_ Direction,
	executeUpdates bool,
	eventLogList *[]GameEventLog,
) (reveal bool, err error) {
	// if not on left side of the map, don't do anything
	if spellPosition.X != 0 {
		return false, nil
	}

	return ResolveWallHeal(world, spellPosition, executeUpdates, eventLogList)
}
