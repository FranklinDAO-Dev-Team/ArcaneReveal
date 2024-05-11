package component

import (
	"log"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

const Ability2ID = 2

type Ability2 struct{}

var _ Ability = &Ability2{}

func (Ability2) GetAbilityID() int {
	return Ability2ID
}

func (a Ability2) Resolve(
	world cardinal.WorldContext,
	gameID types.EntityID,
	spellPosition *Position,
	direction Direction,
	executeUpdates bool,
	eventLogList *[]GameEventLog,
) (reveal bool, err error) {
	perpDirOne := direction.rotateClockwise()
	perpDirTwo := direction.rotateClockwise().rotateClockwise().rotateClockwise()

	damageDealtOne, err := resolveOneA2Check(world, gameID, spellPosition, perpDirOne, executeUpdates, eventLogList)
	if err != nil {
		log.Println("Ability2.Resolve err 1: ", err)
		return false, err
	}
	damageDealtTwo, err := resolveOneA2Check(world, gameID, spellPosition, perpDirTwo, executeUpdates, eventLogList)
	if err != nil {
		log.Println("Ability2.Resolve err 2: ", err)
		return false, err
	}

	reveal = damageDealtOne || damageDealtTwo

	return reveal, nil
}

func resolveOneA2Check(
	world cardinal.WorldContext,
	gameID types.EntityID,
	spellPosition *Position,
	perpDir Direction,
	executeUpdates bool,
	eventLogList *[]GameEventLog,
) (reveal bool, err error) {
	adjPos, err := spellPosition.GetUpdateFromDirection(perpDir)
	if err != nil {
		return false, err
	}
	hitWall, err := IsCollisonThere(world, gameID, *adjPos)
	if err != nil {
		return false, err
	}
	if hitWall { // this spell cannot hit through walls
		return false, nil
	}
	adjPlayablePos, err := adjPos.GetUpdateFromDirection(perpDir)
	// log.Println("adjPlayablePos", adjPlayablePos)
	if err == nil {
		damageDealt, err := damageAtPosition(world, gameID, adjPlayablePos, executeUpdates, false)
		if err != nil {
			return false, err
		}
		if damageDealt {
			gameEvent := GameEventLog{X: adjPlayablePos.X, Y: adjPlayablePos.Y, Event: GameEventSpellDamage}
			*eventLogList = append(*eventLogList, gameEvent)
			return true, nil
		}
	}
	return false, nil
}
