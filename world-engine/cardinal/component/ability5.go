package component

import "pkg.world.dev/world-engine/cardinal"

const Ability5ID = 5

type Ability5 struct{}

var _ Ability = &Ability5{}

func (Ability5) GetAbilityID() int {
	return Ability5ID
}

func (a Ability5) Resolve(
	world cardinal.WorldContext,
	spellPosition *Position,
	_ Direction,
	executeUpdates bool,
	eventLogList *[]GameEventLog,
) (reveal bool, err error) {
	// if not on top side of the map, don't do anything
	if spellPosition.Y != 0 {
		return false, nil
	}

	playerID, err := QueryPlayerID(world)
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

	*eventLogList = append(*eventLogList, GameEventLog{X: spellPosition.X, Y: spellPosition.Y, Event: GameEventSpellWallActivation})
	return true, nil
}
