package component

import (
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
	perpDirTwo := (direction + 3) % 4

	damageDealtOne, err := resolveOneA2Check(world, spellPosition, perpDirOne, executeUpdates, eventLogList)
	if err != nil {
		return false, err
	}
	damageDealtTwo, err := resolveOneA2Check(world, spellPosition, perpDirTwo, executeUpdates, eventLogList)
	if err != nil {
		return false, err
	}

	reveal = damageDealtOne || damageDealtTwo

	return reveal, nil
}

func resolveOneA2Check(world cardinal.WorldContext, spellPosition *Position, perpDir Direction, executeUpdates bool, eventLogList *[]GameEventLog) (reveal bool, err error) {
	adjPos, err := spellPosition.GetUpdateFromDirection(perpDir)
	if err != nil {
		return false, err
	}
	hitWall, err := IsCollisonThere(world, *adjPos)
	if err != nil {
		return false, err
	}
	if !hitWall { // this spell cannot hit through walls
		adjPlayablePos, err := adjPos.GetUpdateFromDirection(perpDir)
		// fmt.Println("adjPlayablePos", adjPlayablePos)
		if err == nil {
			damageDealt, err := damageAtPostion(world, adjPlayablePos, executeUpdates, false)
			if err != nil {
				return false, err
			}
			if damageDealt {
				*eventLogList = append(*eventLogList, GameEventLog{X: adjPlayablePos.X, Y: adjPlayablePos.Y, Event: GameEventSpellDamage})
				return true, nil
			}
		}
	}
	return false, nil
}
