package component

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

const Ability6ID = 6

type Ability6 struct{}

var _ Ability = &Ability6{}

func (Ability6) GetAbilityID() int {
	return Ability6ID
}

func (a Ability6) Resolve(
	world cardinal.WorldContext,
	gameID types.EntityID,
	spellPosition *Position,
	_ Direction,
	executeUpdates bool,
	eventLogList *[]GameEventLog,
) (reveal bool, err error) {
	// if not on right side of the map, don't do anything
	if spellPosition.X != MaxX {
		return false, nil
	}

	return ResolveWallHeal(world, gameID, spellPosition, executeUpdates, eventLogList)
}
