package component

import (
	"pkg.world.dev/world-engine/cardinal"
)

type Ability_2 struct{}

var _ Ability = &Ability_2{}

func (Ability_2) GetAbilityID() int {
	return 2
}

func (Ability_2) Resolve(world cardinal.WorldContext, spellPosition *Position, direction Direction) (reveal bool, err error) {
	perpDirOne := (direction + 1) % 4
	damageDeltOne := false
	AdjOne, err := spellPosition.GetUpdateFromDirection(perpDirOne)
	if err == nil {
		damageDeltOne, err = damageAtPostion(world, AdjOne, false)
		if err != nil {
			return false, err
		}
	}

	perpDirTwo := (direction + 3) % 4
	damageDeltTwo := false
	AdjTwo, err := spellPosition.GetUpdateFromDirection(perpDirTwo)
	if err == nil {
		damageDeltTwo, err = damageAtPostion(world, AdjTwo, false)
		if err != nil {
			return false, err
		}
	}

	reveal = damageDeltOne || damageDeltTwo

	return
}
