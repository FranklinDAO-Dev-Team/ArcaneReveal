package component

import (
	"pkg.world.dev/world-engine/cardinal"
)

const Ability1ID = 1

type Ability1 struct{}

var _ Ability = &Ability1{}

func (Ability1) GetAbilityID() int {
	return Ability1ID
}

// Resolves effects of the ability
// i.e. checks if it activates and if so updates the world
// return if the ability should be revealed and if there was an error
func (Ability1) Resolve(
	world cardinal.WorldContext,
	spellPosition *Position,
	_ Direction,
	executeUpdates bool,
) (reveal bool, err error) {
	damageDealt, err := damageAtPostion(world, spellPosition, executeUpdates, false)
	if err != nil {
		return false, err
	}
	return damageDealt, nil
}
