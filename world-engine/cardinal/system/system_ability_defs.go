package system

import (
	comp "cinco-paus/component"

	"pkg.world.dev/world-engine/cardinal"
)

type Ability_1 struct{}

func (Ability_1) GetAbilityID() int {
	return 1
}

// Resolves effects of the ability
// i.e. checks if it activates and if so updates the world
// return if the ability should be revealed and if there was an error
func Resolve(world cardinal.WorldContext, spell *spellhead) (reveal bool, err error) {
	reveal = false
	err = nil

	overlappedEntity, err := getEntityIDByPosition(world , *spell.pos)
	if err != nil {
		return false, err
	}
	entity, err := cardinal.GetComponent[](world, id)
	// check if it is a monser 
	// if so damage enemy and set reveal = true
	

	return
}
