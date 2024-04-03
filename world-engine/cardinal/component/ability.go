package component

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
)

type Ability interface {
	GetAbilityID() int
	Resolve(cardinal.WorldContext, *Position, Direction) (bool, error)
}

var AbilityMap = map[int]Ability{
	1: &Ability_1{},
	2: &Ability_2{},
}

func damageAtPostion(world cardinal.WorldContext, pos *Position, includePlayer bool) (damageDelt bool, err error) {
	id, err := pos.getEntityIDByPosition(world)
	if err != nil {
		return false, err
	}
	colType, err := cardinal.GetComponent[CollisionType](world, id)
	switch colType.Type {
	case "monster":
		fmt.Println("damage delt at ", pos)
		return true, decrementHealth(world, id)
	case "player":
		fmt.Println("damage delt at ", pos)
		if includePlayer {
			return true, decrementHealth(world, id)
		}
	}

	return false, err
}
