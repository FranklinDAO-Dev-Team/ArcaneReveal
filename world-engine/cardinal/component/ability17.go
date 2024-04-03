package component

import (
	"pkg.world.dev/world-engine/cardinal"
)

type Ability_17 struct{}

var _ Ability = &Ability_17{}

func (Ability_17) GetAbilityID() int {
	return 2
}

// todo: code heal_bottom ability
func (Ability_17) Resolve(world cardinal.WorldContext, spellPosition *Position, direction Direction) (reveal bool, err error) {
	return false, nil
}
