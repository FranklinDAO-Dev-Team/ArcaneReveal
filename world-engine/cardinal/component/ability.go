package component

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

var AbilityMap = map[int]Ability{
	1:  &Ability1{},  // damage attack
	2:  &Ability2{},  // side damage attack
	3:  &Ability3{},  // wall damage attack
	4:  &Ability4{},  // explosion
	5:  &Ability5{},  // up heal
	6:  &Ability6{},  // right heal
	7:  &Ability7{},  // down heal
	8:  &Ability8{},  // left heal
	9:  &Ability9{},  // heal monster
	10: &Ability10{}, // polymorph
}

type GameEvent int

type GameEventLog struct {
	X     int
	Y     int
	Event GameEvent
}

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
	GameEventMonsterHeal                          // 14
	GameEventMonsterPolymorph                     // 15
)

type Ability interface {
	GetAbilityID() int
	Resolve(
		world cardinal.WorldContext,
		gameID types.EntityID,
		spellPosition *Position,
		direction Direction,
		executeUpdates bool,
		eventLogList *[]GameEventLog,
	) (reveal bool, err error)
}

func damageAtPosition(
	world cardinal.WorldContext,
	gameID types.EntityID,
	pos *Position,
	executeUpdates bool,
	includePlayer bool,
) (damageDelt bool, err error) {
	// Lookup if entity exists
	found, id, err := pos.GetEntityIDByPosition(world, gameID)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}

	return damageEntity(world, id, executeUpdates, includePlayer)
}

func damageEntity(
	world cardinal.WorldContext,
	id types.EntityID,
	executeUpdates bool,
	includePlayer bool,
) (bool, error) {
	colType, err := cardinal.GetComponent[Collidable](world, id)
	if err != nil {
		return false, err
	}

	switch colType.Type {
	case MonsterCollide:
		return decrementHealthIfNeeded(world, id, executeUpdates)
	case PlayerCollide:
		if includePlayer {
			return decrementHealthIfNeeded(world, id, executeUpdates)
		}
		return false, nil
	case WallCollide:
		return false, nil
	case ItemCollide:
		return false, nil
	default:
		return false, nil
	}
}

// Checks
func decrementHealthIfNeeded(
	world cardinal.WorldContext,
	id types.EntityID,
	executeUpdates bool,
) (bool, error) {
	if executeUpdates {
		err := DecrementHealth(world, id)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func incrementHealthIfNeeded(
	world cardinal.WorldContext,
	id types.EntityID,
	executeUpdates bool,
) (bool, error) {
	if executeUpdates {
		err := IncrementHealth(world, id)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
