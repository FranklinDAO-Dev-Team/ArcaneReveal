package component

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
)

type Ability_1 struct{}

func (Ability_1) GetAbilityID() int {
	return 1
}

var _ Ability = &Ability_1{}

// Resolves effects of the ability
// i.e. checks if it activates and if so updates the world
// return if the ability should be revealed and if there was an error
func (Ability_1) Resolve(world cardinal.WorldContext, spellPosition Position) (reveal bool, err error) {
	reveal = false
	err = nil

	id, err := spellPosition.getEntityIDByPosition(world)
	if err != nil {
		return false, err
	}
	colType, err := cardinal.GetComponent[CollisionType](world, id)
	switch colType.Type {
	case "monster":
		monster_health, err := cardinal.GetComponent[Health](world, id)
		if err != nil {
			return false, err
		}
		monster_health.CurrHealth--
		err = cardinal.SetComponent[Health](world, id, monster_health)
		if err != nil {
			return false, err
		}
	// case "player":
	// 	fmt.Println("Handling collision with a player")
	// case "wall":
	// 	fmt.Println("Handling collision with a wall")
	// case "item":
	// 	fmt.Println("Handling collision with an item")
	default:
		fmt.Println("Unknown collision type")
	}

	// check if it is a monser
	// if so damage enemy and set reveal = true

	return
}
