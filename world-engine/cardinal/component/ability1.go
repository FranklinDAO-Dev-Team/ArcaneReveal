package component

import (
	"log"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// damage attack
const Ability1ID = 1

type Ability1 struct{}

var _ Ability = &Ability1{}

func (Ability1) GetAbilityID() int {
	return Ability1ID
}

func (Ability1) GetAbilityName() string {
	return "damage attack"
}

// Resolves effects of the ability
// i.e. checks if it activates and if so updates the world
// return if the ability should be revealed and if there was an error
func (Ability1) Resolve(
	world cardinal.WorldContext,
	gameID types.EntityID,
	spellPosition *Position,
	_ Direction,
	executeUpdates bool,
	eventLogList *[]GameEventLog,
) (reveal bool, err error) {
	damageDealt, err := damageAtPosition(world, gameID, spellPosition, executeUpdates, false)
	if err != nil {
		log.Println("Ability1.Resolve err: ", err)
		return false, err
	}
	if damageDealt {
		gameEvent := GameEventLog{X: spellPosition.X, Y: spellPosition.Y, Event: GameEventSpellDamage}
		*eventLogList = append(*eventLogList, gameEvent)
	}
	return damageDealt, nil
}
