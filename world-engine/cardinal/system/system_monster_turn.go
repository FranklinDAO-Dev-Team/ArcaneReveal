package system

import (
	comp "cinco-paus/component"
	"fmt"

	"math/rand"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

func MonsterTurnSystem(world cardinal.WorldContext) error {
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
		fmt.Println("MonsterTurnSystem. Monster id: ", id)
		monsterPos, err := cardinal.GetComponent[comp.Position](world, id)
		if err != nil {
			turnErr = err
			return false // if error, break out of search
		}
		manDist := monsterPos.ManhattenDistance(playerPos)
		if manDist == 1 {
			// attack since player is in range
			playerID, err := queryPlayerID(world)
			if err != nil {
				turnErr = err
				return false
			}
			comp.DecrementHealth(world, playerID)
			fmt.Println("successful monster action: attack")
		} else {
			// get move options (places that are legal, not moving into a wall, etc),
			direction, err := decideMonsterMovementDirection(world, monsterPos, playerPos)
			if err != nil {
				turnErr = err
				return false // if error, break out of search
			}
			newMonsterPos, err := monsterPos.GetUpdateFromDirection(direction)
			if err != nil {
				turnErr = err
				return false // if error, break out of search
			}
			fmt.Println("MonsterTurnSystem. New monster position: ", newMonsterPos)
			err = cardinal.SetComponent[comp.Position](world, id, newMonsterPos)
			if err != nil {
				turnErr = err
				return false // if error, break out of search
			}
			fmt.Println("successful monster action: Move")
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
	newMonsterPos, err := monsterPos.GetUpdateFromDirection(direction)
	if err != nil {
		return false, 0, nil // invalid postion, but don't return error, just check next direction
	}
	manDist = monsterPos.ManhattenDistance(playerPos)

	found, id, err := newMonsterPos.GetEntityIDByPosition(world)
	if err != nil {
		return false, 0, err
	}
	if found {
		colType, err := cardinal.GetComponent[comp.Collidable](world, id)
		if err != nil {
			return false, 0, err
		}
		switch colType.Type {
		case comp.ItemCollide:
			return true, manDist, nil
		default:
			return false, -1, nil
		}
	}
	// no entity found, so it's a valid move
	return true, manDist, nil

}
