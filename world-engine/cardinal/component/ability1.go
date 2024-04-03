package component

import (
	"pkg.world.dev/world-engine/cardinal"
)

type Ability_1 struct{}

var _ Ability = &Ability_1{}

func (Ability_1) GetAbilityID() int {
	return 1
}

// Resolves effects of the ability
// i.e. checks if it activates and if so updates the world
// return if the ability should be revealed and if there was an error
func (Ability_1) Resolve(world cardinal.WorldContext, spellPosition *Position, direction Direction) (reveal bool, err error) {
	damage_delt, err := damageAtPostion(world, spellPosition, false)
	if err != nil {
		return false, err
	}
	return damage_delt, nil
}
