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
	if spellPosition.X != 10 {
		return false, nil
	}

	if !executeUpdates {
		return true, nil
	}

	playerID, err := QueryPlayerID(world)
	if err != nil {
		return false, err
	}

	inc, err := IncrementHealth(world, playerID)
	if err != nil {
		return false, err
	}
	if inc {
		*eventLogList = append(*eventLogList, GameEventLog{X: spellPosition.X, Y: spellPosition.Y, Event: GameEventSpellWallActivation})
	} else {
		*eventLogList = append(*eventLogList, GameEventLog{X: spellPosition.X, Y: spellPosition.Y, Event: GameEventSpellDisappate})
	}

	return inc, nil
}
