package component

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

const Ability7ID = 7

type Ability7 struct{}

var _ Ability = &Ability7{}

func (Ability7) GetAbilityID() int {
	return Ability7ID
}

func (a Ability7) Resolve(
	world cardinal.WorldContext,
	gameID types.EntityID,
	spellPosition *Position,
	_ Direction,
	executeUpdates bool,
	eventLogList *[]GameEventLog,
) (reveal bool, err error) {
	// if not on bottom side of the map, don't do anything
	if spellPosition.Y != MaxY {
		// log.Printf("Ability7 pos %v exit 1\n", spellPosition)
		return false, nil
	}

	return ResolveWallHeal(world, gameID, spellPosition, executeUpdates, eventLogList)
}
