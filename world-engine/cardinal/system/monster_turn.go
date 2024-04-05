package system

import (
	comp "cinco-paus/component"
	"fmt"

	"math/rand"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

func MonsterTurnSystem(world cardinal.WorldContext, eventLogList *[]comp.GameEventLog) error {
	fmt.Println("MonsterTurnSystem")
	var turnErr error
	playerPos, err := cardinal.GetComponent[comp.Position](world, 0)
	if err != nil {
		return err
	}
	searchErr := cardinal.NewSearch(
		world,
		filter.Contains(comp.Monster{}),
	).Each(func(id types.EntityID) bool {
		// get original monster position
		fmt.Printf("Monster id: %d\n", id)
		origMonsterPos, err := cardinal.GetComponent[comp.Position](world, id)
		if err != nil {
			fmt.Println("MonsterTurn err 1")
			turnErr = err
			return false // if error, break out of search
		}
		// get manhatten distance between original monster position and player position
		manDist := origMonsterPos.ManhattenDistance(playerPos)

		if manDist == 2 {
			// attack since player is in range
			playerID, err := queryPlayerID(world)
			if err != nil {
				fmt.Println("MonsterTurn err 2")
				turnErr = err
				return false
			}
			// Decrement health of player
			comp.DecrementHealth(world, playerID)
			// add event to event log
			*eventLogList = append(*eventLogList, comp.GameEventLog{X: origMonsterPos.X, Y: origMonsterPos.Y, Event: comp.GameEventMonsterAttack})
		} else {

			// get move options (places that are legal, not moving into a wall, etc),
			direction, err := decideMonsterMovementDirection(world, origMonsterPos, playerPos)
			if err != nil {
				fmt.Println("MonsterTurn err 3")
				turnErr = err
				return false // if error, break out of search
			}
			fmt.Printf("monsterID: %d, origMonsterPos: %v, direction: %d\n", id, origMonsterPos, direction)

			// calculate new monster position
			newMonsterPos, err := origMonsterPos.GetUpdateFromDirection(direction)
			if err != nil {
				fmt.Println("MonsterTurn err 4")
				turnErr = err
				return false // if error, break out of search
			}
			// calculate new monster position
			newNewMonsterPos, err := newMonsterPos.GetUpdateFromDirection(direction)
			if err != nil {
				fmt.Println("MonsterTurn err 4")
				turnErr = err
				return false // if error, break out of search
			}
			fmt.Printf("newMonsterPos: %v, newNewMonsterPos: %v\n", newMonsterPos, newNewMonsterPos)

			// update monster position onchain
			err = cardinal.SetComponent[comp.Position](world, id, newNewMonsterPos)
			if err != nil {
				fmt.Println("MonsterTurn err 5")
				turnErr = err
				return false // if error, break out of search
			}

			// add event to event log
			*eventLogList = append(*eventLogList, comp.GameEventLog{X: origMonsterPos.X, Y: origMonsterPos.Y, Event: directionToMonsterAttack(direction)})
		}
		return true // always return true to move on to the next monster
	})

	if searchErr != nil {
		return searchErr
	}
	if turnErr != nil {
		fmt.Println("MonsterTurnSystem. Error: ", turnErr)
		return turnErr
	}
	return nil
}

func decideMonsterMovementDirection(world cardinal.WorldContext, monsterPos *comp.Position, playerPos *comp.Position) (comp.Direction, error) {
	bestDirection := comp.Direction(-1)
	bestManDist := 100 // set to a large number to start
	directions := []comp.Direction{comp.LEFT, comp.RIGHT, comp.UP, comp.DOWN}
	// shuffle directions so picked is more random
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})
	for _, direction := range directions {
		valid, manDist, err := CheckMonsterMovementUtility(world, playerPos, monsterPos, direction)
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
	playerPos *comp.Position,
	monsterPos *comp.Position,
	direction comp.Direction,
) (valid bool, manDist int, err error) {
	// check first step for walls
	newMonsterPos, err := monsterPos.GetUpdateFromDirection(direction)
	if err != nil {
		return false, 0, err
	}
	valid, err = isCollisonThere(world, *newMonsterPos)
	if err != nil {
		return false, 0, err
	} else if valid {
		return false, 0, nil // invalid postion, but don't return error, just check next direction
	}

	// check second step for other collisons, ex monster
	newNewMonsterPos, err := newMonsterPos.GetUpdateFromDirection(direction)
	if err != nil {
		return false, 0, err
	}
	valid, err = isCollisonThere(world, *newMonsterPos)
	if err != nil {
		return false, 0, err
	} else if valid {
		return false, 0, nil // invalid postion, but don't return error, just check next direction
	}

	// direction is valid option, return manhatten distance
	manDist = newNewMonsterPos.ManhattenDistance(playerPos)
	return true, manDist, nil
}

func directionToMonsterAttack(direction comp.Direction) comp.GameEvent {
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

func isCollisonThere(world cardinal.WorldContext, pos comp.Position) (bool, error) {
	found, id, err := pos.GetEntityIDByPosition(world)
	if err != nil {
		return false, err
	}
	if found {
		colType, err := cardinal.GetComponent[comp.Collidable](world, id)
		if err != nil {
			return false, err
		}
		switch colType.Type {
		case comp.ItemCollide:
			// ok to overlap items
			return false, nil
		default:
			// not ok to overlap other types of collidable
			return true, nil
		}
	}
	// no entity found, so it's not a wall
	return false, nil
}
