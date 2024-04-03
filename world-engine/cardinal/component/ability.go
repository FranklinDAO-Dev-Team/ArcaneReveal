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
	found, id, err := pos.GetEntityIDByPosition(world)
	if err != nil {
		return false, err
	}
	if found {
		colType, err := cardinal.GetComponent[Collidable](world, id)
		if err != nil {
			return false, err
		}
		switch colType.Type {
		case MonsterCollide:
			fmt.Println("damage delt at ", pos)
			return true, DecrementHealth(world, id)
		case PlayerCollide:
			fmt.Println("damage delt at ", pos)
			if includePlayer {
				return true, DecrementHealth(world, id)
			}
		}
	}

	return false, err
}
