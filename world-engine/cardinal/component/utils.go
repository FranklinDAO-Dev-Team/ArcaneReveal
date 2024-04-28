package component

import (
	"fmt"
	"log"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

func IsCollisonThere(world cardinal.WorldContext, pos Position) (bool, error) {
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
		case ItemCollide:
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

func QueryPlayerID(world cardinal.WorldContext, gameID types.EntityID) (types.EntityID, error) {
	var outerErr error
	var playerID types.EntityID
	searchErr := cardinal.NewSearch(
		world,
		filter.Contains(Player{})).
		Each(func(id types.EntityID) bool {
			playerGameObj, err := cardinal.GetComponent[GameObj](world, id)
			if err != nil {
				log.Println("QueryPlayerID err 1: ", err)
				outerErr = err
				return false
			}
			if playerGameObj.GameID == gameID {
				playerID = id
				return false
			}

			return true
		})

	if playerID == 0 {
		log.Println("QueryPlayerID no player found")
		return 0, fmt.Errorf("QueryPlayerID no player found")
	}

	if searchErr != nil {
		log.Println("QueryPlayerID err 2: ", searchErr)
		return 0, searchErr
	}
	if outerErr != nil {
		log.Println("QueryPlayerID err 3: ", outerErr)
		return 0, outerErr
	}
	return playerID, nil
}
