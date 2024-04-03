package component

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
)

type Ability interface {
	GetAbilityID() int
	Resolve(cardinal.WorldContext, *Position, Direction, bool) (bool, error)
}

var AbilityMap = map[int]Ability{
	1: &Ability1{},
	2: &Ability2{},
}

func damageAtPostion(
	world cardinal.WorldContext,
	pos *Position,
	executeUpdates bool,
	includePlayer bool,
) (damageDelt bool, err error) {
	// lookup if entity exists
	found, id, err := pos.GetEntityIDByPosition(world)
	if err != nil {
		return false, err
	}
	if found {
		// check entity type
		colType, err := cardinal.GetComponent[Collidable](world, id)
		if err != nil {
			return false, err
		}
		switch colType.Type {
		case MonsterCollide:
			fmt.Println("damage delt at ", pos)
			if executeUpdates {
				err := DecrementHealth(world, id)
				if err != nil {
					return false, err
				}
				return true, nil
			}
		case PlayerCollide:
			if includePlayer {
				fmt.Println("damage delt at ", pos)
				if executeUpdates {
					err := DecrementHealth(world, id)
					if err != nil {
						return false, err
					}
				}
				return true, nil
			} else {
				return false, nil
			}
		default:
			return false, nil
		}
	}
	return false, err
}
