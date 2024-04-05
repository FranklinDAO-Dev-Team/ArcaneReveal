package component

import (
	"pkg.world.dev/world-engine/cardinal"
)

type GameEvent int

const (
	GameEventSpellBeam           GameEvent = iota // 0
	GameEventSpellDamage                          // 1
	GameEventSpellDisappate                       // 2
	GameEventSpellWallActivation                  // 3
	GameEventMonsterAttack                        // 4
	GameEventMonsterUp                            // 5
	GameEventMonsterRight                         // 6
	GameEventMonsterDown                          // 7
	GameEventMonsterLeft                          // 8
	GameEventPlayerAttack                         // 9
	GameEventPlayerUp                             // 10
	GameEventPlayerRight                          // 11
	GameEventPlayerDown                           // 12
	GameEventPlayerLeft                           // 13
)

type GameEventLog struct {
	X     int
	Y     int
	Event GameEvent
}

type Ability interface {
	GetAbilityID() int
	Resolve(
		world cardinal.WorldContext,
		spellPosition *Position,
		direction Direction,
		executeUpdates bool,
		eventLogList *[]GameEventLog,
	) (reveal bool, err error)
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
			if executeUpdates {
				err := DecrementHealth(world, id)
				if err != nil {
					return false, err
				}
				return true, nil
			}
		case PlayerCollide:
			if includePlayer {
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
