package system

import (
	comp "cinco-paus/component"
	"errors"

	"pkg.world.dev/world-engine/cardinal"
)

// resolveAbilities takes information about a spell and determines game events it causes
// if updateChainState is true, it applies the changes to the world
// potentialAbilities is a list of booleans that indicate if the ability at that index should be considered for activation
// spellEvents are recorded in eventLogList
// loop executes in 3 steps:
// 1. record abilities that could activate a current square
// 2. get next spell position
// 3. if wall entity at spellPos, stop
func resolveAbilities(
	world cardinal.WorldContext,
	spell *comp.Spell,
	spellPos *comp.Position,
	potentialAbilities *[comp.TotalAbilities]bool,
	updateChainState bool,
	eventLogList *[]comp.GameEventLog,
) error {
	for !spell.Expired {
		// log SpellBeam position
		*eventLogList = append(*eventLogList, comp.GameEventLog{X: spellPos.X, Y: spellPos.Y, Event: comp.GameEventSpellBeam})
		// record abilities that could activate a current square
		err := resolveAbilitiesAtPosition(world, spellPos, spell.Direction, potentialAbilities, updateChainState, eventLogList)
		if err != nil {
			return err
		}

		// get next spell position
		spellPos, err = spellPos.GetUpdateFromDirection(spell.Direction)
		if err != nil {
			spell.Expired = true
		}
		if spellPos == nil {
			spell.Expired = true
			break
		}

		// if wall entity at spellPos, stop
		found, id, err := spellPos.GetEntityIDByPosition(world)
		if err != nil {
			return err
		}
		if found {
			colType, err := cardinal.GetComponent[comp.Collidable](world, id)
			if err != nil {
				return err
			}
			if colType.Type == comp.WallCollide {
				spell.Expired = true
			}
		}
	}
	return nil
}

// resolveAbilitiesAtPosition takes information about a spell determines game events it causes at a single position
// if updateChainState is true, it applies the changes to the world
// potentialAbilities is a list of booleans that indicate if the ability at that index should be considered for activation
// spellEvents are recorded in eventLogList
func resolveAbilitiesAtPosition(
	world cardinal.WorldContext,
	spellPos *comp.Position,
	direction comp.Direction,
	potentialAbilities *[comp.TotalAbilities]bool,
	updateChainState bool,
	eventLogList *[]comp.GameEventLog,
) error {
	for i := 0; i < len(*potentialAbilities); i++ {
		if (*potentialAbilities)[i] { // if ability should be activated/checked
			a := comp.AbilityMap[i+1]
			if a == nil {
				return errors.New("unknown ability called")
			}
			activated, err := a.Resolve(world, spellPos, direction, updateChainState, eventLogList)
			// only overwrite if ability activated
			(*potentialAbilities)[i] = activated || (*potentialAbilities)[i]
			if err != nil {
				return err
			}
		}
	}

	return nil
}
