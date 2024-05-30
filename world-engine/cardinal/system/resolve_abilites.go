package system

import (
	comp "cinco-paus/component"
	"cinco-paus/seismic/client"
	"errors"
	"log"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// resolveAbilities takes information about a spell and determines game events it causes
// if updateChainState is true, it applies the changes to the world
// potentialAbilities is a list of booleans that indicate
// if the ability at that index should be considered for activation
// spellEvents are recorded in eventLogList
// loop executes in 3 steps:
// 1. record abilities that could activate a current square
// 2. get next spell position
// 3. if wall entity at spellPos, stop
func resolveAbilities(
	world cardinal.WorldContext,
	gameID types.EntityID,
	spell *comp.Spell,
	spellPos *comp.Position,
	potentialAbilities *[client.TotalAbilities]bool,
	updateChainState bool,
	eventLogList *[]comp.GameEventLog,
) error {
	for !spell.Expired {
		// log.Println()
		// log.Println("resolveAbilities. start Spell position: ", spellPos.X, spellPos.Y)

		// log SpellBeam position
		*eventLogList = append(*eventLogList, comp.GameEventLog{X: spellPos.X, Y: spellPos.Y, Event: comp.GameEventSpellBeam})
		// record abilities that could activate a current square
		err := resolveAbilitiesAtPosition(
			world,
			gameID,
			spellPos,
			spell.Direction,
			potentialAbilities,
			updateChainState,
			eventLogList,
		)
		if err != nil {
			log.Println("resolveAbilitiesAtPosition err: ", err)
			return err
		}
		// log.Printf("potentialAbilities: %v", *potentialAbilities)

		// get next spell position
		spellPos, err = spellPos.GetUpdateFromDirection(spell.Direction)
		if err != nil {
			spell.Expired = true
		}
		if spellPos == nil {
			spell.Expired = true
			break
		}
		// log.Println("resolveAbilities. end Spell position: , err; ", spellPos.X, spellPos.Y, err)
	}
	return nil
}

// resolveAbilitiesAtPosition takes information about a spell determines game events it causes at a single position
// if updateChainState is true, it applies the changes to the world
// potentialAbilities is a list of booleans that indicate if
// the ability at that index should be considered for activation
// spellEvents are recorded in eventLogList
func resolveAbilitiesAtPosition(
	world cardinal.WorldContext,
	gameID types.EntityID,
	spellPos *comp.Position,
	direction comp.Direction,
	potentialAbilities *[client.TotalAbilities]bool,
	updateChainState bool,
	eventLogList *[]comp.GameEventLog,
) error {
	for i := 0; i < len(*potentialAbilities); i++ {
		if potentialAbilities[i] || !updateChainState { // if ability should be activated/checked
			a := comp.AbilityMap[i+1]
			if a == nil {
				return errors.New("unknown ability called")
			}
			activated, err := a.Resolve(world, gameID, spellPos, direction, updateChainState, eventLogList)
			if err != nil {
				return err
			}
			// if activated {
			// 	log.Printf("resolveAbilitiesAtPosition() activated ability %d (%s) \n", i, a.GetAbilityName())
			// }

			// only overwrite if ability activated
			potentialAbilities[i] = activated || potentialAbilities[i]

			// if spellPos.Y == 10 && i == 6 {
			// 	log.Println("6+1 = 7 activated", activated)
			// 	log.Println(potentialAbilities)
			// 	log.Println()
			// }
		}
	}

	return nil
}
