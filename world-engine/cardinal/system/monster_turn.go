package system

import (
	comp "cinco-paus/component"
	"fmt"
	"log"

	"math/rand"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

const playerAttackDistance = 2

func MonsterTurnSystem(
	world cardinal.WorldContext,
	gameID types.EntityID,
	eventLogList *[]comp.GameEventLog,
) error {
	log.Println("MonsterTurnSystem")
	var turnErr error

	playerID, err := comp.QueryPlayerID(world, gameID)
	if err != nil {
		return err
	}
	playerPos, err := cardinal.GetComponent[comp.Position](world, playerID)
	if err != nil {
		return err
	}
	searchErr := cardinal.NewSearch(
		world,
		filter.Contains(comp.Monster{}),
	).Each(func(id types.EntityID) bool {
		// get original monster position
		origMonsterPos, err := cardinal.GetComponent[comp.Position](world, id)
		if err != nil {
			// log.Println("MonsterTurn err 1")
			turnErr = err
			return false // if error, break out of search
		}
		// get manhatten distance between original monster position and player position
		manDist := origMonsterPos.ManhattenDistance(playerPos)

		if manDist == playerAttackDistance {
			turnErr = executeMonsterAttack(world, gameID, eventLogList, origMonsterPos)
			if turnErr != nil {
				return false // if error, break out of search
			}
		} else {
			turnErr = executeMonsterMove(world, gameID, eventLogList, origMonsterPos, id, playerPos)
			if turnErr != nil {
				return false // if error, break out of search
			}
		}
		return true // always return true to move on to the next monster
	})

	if searchErr != nil {
		return searchErr
	}
	if turnErr != nil {
		log.Println("MonsterTurnSystem. Error: ", turnErr)
		return turnErr
	}
	return nil
}

func executeMonsterAttack(
	world cardinal.WorldContext,
	gameID types.EntityID,
	eventLogList *[]comp.GameEventLog,
	origMonsterPos *comp.Position,
) error {
	// attack since player is in range
	playerID, err := comp.QueryPlayerID(world, gameID)
	if err != nil {
		return err
	}
	// Decrement health of player
	err = comp.DecrementHealth(world, playerID)
	if err != nil {
		return err
	}

	// add event to event log
	monsterAttackEvent := comp.GameEventLog{X: origMonsterPos.X, Y: origMonsterPos.Y, Event: comp.GameEventMonsterAttack}
	*eventLogList = append(*eventLogList, monsterAttackEvent)

	return nil
}

func executeMonsterMove(
	world cardinal.WorldContext,
	gameID types.EntityID,
	eventLogList *[]comp.GameEventLog,
	origMonsterPos *comp.Position,
	monsterID types.EntityID,
	playerPos *comp.Position,
) error {
	// get move options (places that are legal, not moving into a wall, etc),
	direction, err := decideMonsterMovementDirection(world, gameID, origMonsterPos, playerPos)
	if err != nil {
		return err
	}
	moveEventType := directionToMonsterMoveEvent(direction)

	monsterMoveEvent1 := comp.GameEventLog{X: origMonsterPos.X, Y: origMonsterPos.Y, Event: moveEventType}
	*eventLogList = append(*eventLogList, monsterMoveEvent1)

	// calculate new monster position
	newMonsterPos, err := origMonsterPos.GetUpdateFromDirection(direction)
	if err != nil {
		return err
	}
	monsterMoveEvent2 := comp.GameEventLog{X: newMonsterPos.X, Y: newMonsterPos.Y, Event: moveEventType}
	*eventLogList = append(*eventLogList, monsterMoveEvent2)

	// calculate new monster position
	newNewMonsterPos, err := newMonsterPos.GetUpdateFromDirection(direction)
	if err != nil {
		return err
	}

	// double check that newNewMonsterPos is not colliding with anything
	found, eID, err := newNewMonsterPos.GetEntityIDByPosition(world, gameID)
	if err != nil {
		return err
	}
	if found {
		return fmt.Errorf("executeMonsterMove() newNewMonsterPos collides with entity %d", eID)
	}

	// update monster position onchain
	err = cardinal.SetComponent[comp.Position](world, monsterID, newNewMonsterPos)
	if err != nil {
		return err
	}
	return nil
}

func decideMonsterMovementDirection(
	world cardinal.WorldContext,
	gameID types.EntityID,
	monsterPos *comp.Position,
	playerPos *comp.Position,
) (comp.Direction, error) {
	bestDirection := comp.Direction(-1)
	bestManDist := 100 // set to a large number to start
	directions := []comp.Direction{comp.LEFT, comp.RIGHT, comp.UP, comp.DOWN}
	// shuffle directions so picked is more random
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})
	for _, direction := range directions {
		valid, manDist, err := CheckMonsterMovementUtility(world, gameID, playerPos, monsterPos, direction)
		if err != nil {
			return -1, err
		}
		if valid && manDist < bestManDist {
			bestManDist = manDist
			bestDirection = direction
		}
	}
	return bestDirection, nil
}

func CheckMonsterMovementUtility(
	world cardinal.WorldContext,
	gameID types.EntityID,
	playerPos *comp.Position,
	monsterPos *comp.Position,
	direction comp.Direction,
) (valid bool, manDist int, err error) {
	// check first step for walls
	newMonsterPos, err := monsterPos.GetUpdateFromDirection(direction)
	if err != nil {
		return false, 0, err
	}
	invalid, err := comp.IsCollisonThere(world, gameID, *newMonsterPos)
	if err != nil {
		return false, 0, err
	} else if invalid {
		return false, 0, nil // invalid position, but don't return error, just check next direction
	}

	// check second step for other collisons, ex monster
	newNewMonsterPos, err := newMonsterPos.GetUpdateFromDirection(direction)
	if err != nil {
		return false, 0, err
	}
	invalid, err = comp.IsCollisonThere(world, gameID, *newNewMonsterPos)
	if err != nil {
		return false, 0, err
	} else if invalid {
		return false, 0, nil // invalid position, but don't return error, just check next direction
	}

	// direction is valid option, return manhatten distance
	manDist = newNewMonsterPos.ManhattenDistance(playerPos)
	return true, manDist, nil
}

func directionToMonsterMoveEvent(direction comp.Direction) comp.GameEvent {
	switch direction {
	case comp.LEFT:
		return comp.GameEventMonsterLeft
	case comp.RIGHT:
		return comp.GameEventMonsterRight
	case comp.UP:
		return comp.GameEventMonsterUp
	case comp.DOWN:
		return comp.GameEventMonsterDown
	default:
		panic("invalid direction")
	}
}
