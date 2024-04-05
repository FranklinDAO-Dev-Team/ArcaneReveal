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
	// if not on right side of the map, don't do anything
	if spellPosition.X != 10 {
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
