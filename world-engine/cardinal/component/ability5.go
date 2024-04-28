package component

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

const Ability5ID = 5

type Ability5 struct{}

var _ Ability = &Ability5{}

func (Ability5) GetAbilityID() int {
	return Ability5ID
}

func (a Ability5) Resolve(
	world cardinal.WorldContext,
	gameID types.EntityID,
	spellPosition *Position,
	_ Direction,
	executeUpdates bool,
	eventLogList *[]GameEventLog,
) (reveal bool, err error) {
	// if not on top side of the map, don't do anything
	if spellPosition.Y != 0 {
		return false, nil
	}

	return ResolveWallHeal(world, gameID, spellPosition, executeUpdates, eventLogList)
}

func ResolveWallHeal(
	world cardinal.WorldContext,
	gameID types.EntityID,
	spellPosition *Position,
	executeUpdates bool,
	eventLogList *[]GameEventLog,
) (reveal bool, err error) {
	playerID, err := QueryPlayerID(world, gameID)
	if err != nil {
		return false, err
	}

	playerHealth, err := cardinal.GetComponent[Health](world, playerID)
	if err != nil {
		return false, err
	}
	if playerHealth.CurrHealth == playerHealth.MaxHealth {
		return false, err // ability cannot activate if player is at max health
	}

	if executeUpdates {
		err := IncrementHealth(world, playerID)
		if err != nil {
			return false, err
		}
	}
	gameEvent := GameEventLog{X: spellPosition.X, Y: spellPosition.Y, Event: GameEventSpellWallActivation}
	*eventLogList = append(*eventLogList, gameEvent)
	return true, nil
}
