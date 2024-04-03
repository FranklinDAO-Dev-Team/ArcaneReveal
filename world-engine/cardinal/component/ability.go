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
			err := DecrementHealth(world, id)
			if err != nil {
				return false, err
			}
			return true, nil
		case PlayerCollide:
			if includePlayer {
				fmt.Println("damage delt at ", pos)
				err := DecrementHealth(world, id)
				if err != nil {
					return false, err
				}
				return true, nil
			} else {
				return false, nil
			}
		}
	}

	return false, err
}
