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
	GameEventMonsterPolymorph0                    // 15
	GameEventMonsterPolymorph1                    // 16
	GameEventMonsterPolymorph2                    // 17
	GameEventMonsterPolymorph3                    // 18
)

type Ability interface {
	GetAbilityID() int
	GetAbilityName() string
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
	colType, err := cardinal.GetComponent[Collidable](world, id)
	if err != nil {
		return false, err
	}
	if colType.Type == WallCollide {
		return false, nil
	}

	return true, DamageEntity(world, gameID, id, executeUpdates, includePlayer)
}

// DamageEntity decrements the health of an entity and handles death if health reaches 0
func DamageEntity(
	world cardinal.WorldContext,
	gameID types.EntityID,
	id types.EntityID,
	executeUpdates bool,
	includePlayer bool,
) error {
	colType, err := cardinal.GetComponent[Collidable](world, id)
	if err != nil {
		return err
	}
	switch colType.Type {
	case MonsterCollide:
		err = decrementHealthIfNeeded(world, id, executeUpdates)
		if err != nil {
			return err
		}
		return handleEntityDeath(world, gameID, id)
	case PlayerCollide:
		if includePlayer {
			err = decrementHealthIfNeeded(world, id, executeUpdates)
			if err != nil {
				return err
			}
			return handleEntityDeath(world, gameID, id)
		}
		return nil
	case WallCollide:
		return nil
	case ItemCollide:
		return nil
	default:
		return nil
	}
}

// Checks
func decrementHealthIfNeeded(
	world cardinal.WorldContext,
	id types.EntityID,
	executeUpdates bool,
) error {
	if executeUpdates {
		err := DecrementHealth(world, id)
		if err != nil {
			return err
		}
	}
	return nil
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

func handleEntityDeath(world cardinal.WorldContext, gameID types.EntityID, entityID types.EntityID) (err error) {
	health, err := cardinal.GetComponent[Health](world, entityID)
	if err != nil {
		return err
	}
	if health.CurrHealth > 0 {
		// if not dead, nothing to do
		return nil
	}
	// entity is dead, so remove it
	// Add to score if it's a monster
	colType, err := getCollisionType(world, entityID)
	if err != nil {
		return err
	}
	if colType == MonsterCollide {
		game, err := cardinal.GetComponent[Game](world, gameID)
		if err != nil {
			return err
		}
		monsterType, err := cardinal.GetComponent[Monster](world, entityID)
		if err != nil {
			return err
		}
		game.Score += 10 * (int(monsterType.Type) + 1)
		cardinal.SetComponent(world, gameID, game)
	}
	err = cardinal.Remove(world, entityID)
	if err != nil {
		return err
	}
	return nil
}
