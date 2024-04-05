package component

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
)

const Ability2ID = 2

type Ability2 struct{}

var _ Ability = &Ability2{}

func (Ability2) GetAbilityID() int {
	return Ability2ID
}

func (a Ability2) Resolve(
	world cardinal.WorldContext,
	spellPosition *Position,
	direction Direction,
	executeUpdates bool,
	eventLogList *[]GameEventLog,
) (reveal bool, err error) {
	perpDirOne := (direction + 1) % 4
	damageDealtOne := false
	adjOne, err := spellPosition.GetUpdateFromDirection(perpDirOne)
	// adjPlayableOne, err := adjOne.GetUpdateFromDirection(perpDirOne)
	fmt.Println("adjOne", adjOne)
	if err == nil {
		damageDealtOne, err = damageAtPostion(world, adjOne, executeUpdates, false)
		if err != nil {
			return false, err
		}
		if damageDealtOne {
			*eventLogList = append(*eventLogList, GameEventLog{X: adjOne.X, Y: adjOne.Y, Event: GameEventSpellDamage})
		}
	}

	perpDirTwo := (direction + 3) % 4
	damageDealtTwo := false
	adjTwo, err := spellPosition.GetUpdateFromDirection(perpDirTwo)
	fmt.Println("adjTwo", adjTwo)
	if err == nil {
		damageDealtTwo, err = damageAtPostion(world, adjTwo, executeUpdates, false)
		if err != nil {
			return false, err
		}
		if damageDealtTwo {
			*eventLogList = append(*eventLogList, GameEventLog{X: adjTwo.X, Y: adjTwo.Y, Event: GameEventSpellDamage})
		}
	}

	reveal = damageDealtOne || damageDealtTwo

	return reveal, nil
}
