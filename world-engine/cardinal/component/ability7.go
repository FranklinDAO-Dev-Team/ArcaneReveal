package component

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
)

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
	if spellPosition.Y == 10 {
		fmt.Printf("Ability7 pos %v entered\n", spellPosition)
	}
	// if not on bottom side of the map, don't do anything
	if spellPosition.Y != 10 {
		// fmt.Printf("Ability7 pos %v exit 1\n", spellPosition)
		return false, nil
	}

	playerID, err := QueryPlayerID(world)
	if err != nil {
		fmt.Printf("Ability7 pos %v exit 2\n", spellPosition)
		return false, err
	}

	playerHealth, err := cardinal.GetComponent[Health](world, playerID)
	if err != nil {
		fmt.Printf("Ability7 pos %v exit 3\n", spellPosition)
		return false, err
	}
	if playerHealth.CurrHealth == playerHealth.MaxHealth {
		fmt.Printf("Ability7 pos %v exit 4\n", spellPosition)
		return false, err // ability cannot activate if player is at max health
	}

	if executeUpdates {
		err := IncrementHealth(world, playerID)
		if err != nil {
			fmt.Printf("Ability7 pos %v exit 5\n", spellPosition)
			return false, err
		}
	}

	*eventLogList = append(*eventLogList, GameEventLog{X: spellPosition.X, Y: spellPosition.Y, Event: GameEventSpellWallActivation})
	fmt.Println("Ability7 activated")
	fmt.Printf("Ability7 pos %v exit 6\n", spellPosition)
	return true, nil
}
