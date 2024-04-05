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
	// if not on left side of the map, don't do anything
	if spellPosition.X != 0 {
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
