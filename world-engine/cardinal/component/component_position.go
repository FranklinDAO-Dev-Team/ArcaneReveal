package component

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

const MAX_X = 4
const MAX_Y = 4

type Position struct {
	X int `json:"x"`
	Y int `josn:"y"`
}

func (Position) Name() string {
	return "Position"
}

func (p Position) UpdateFromDirection(direction string) error {
	switch direction {
	case "left":
		if p.X == 0 {
			return fmt.Errorf("moving out of bounds")
		}
		p.X--
	case "right":
		if p.X == MAX_X {
			return fmt.Errorf("moving out of bounds")
		}
		p.X++
	case "up":
		if p.Y == 0 {
			return fmt.Errorf("moving out of bounds")
		}
		p.Y--
	case "down":
		if p.Y == MAX_Y {
			return fmt.Errorf("moving out of bounds")
		}
		p.Y++
	default:
		return fmt.Errorf("invalid direction")
	}

	return nil
}

type EntityAtLocation struct {
	Monsters []types.EntityID
	Players  []types.EntityID
}

func (p *Position) getEntityIDByPosition(world cardinal.WorldContext) (types.EntityID, error) {
	var eID types.EntityID
	var err error

	// err := cardinal.NewSearch(world, filter.Contains(Monster{}, Position{})).Each(
	// 	// Check if the position is equals p.X and p.Y
	// 	// // If the position is equals p.X and p.Y add it to the EntityAtLocation.Monster
	// )

	searchErr := cardinal.NewSearch(world, filter.Contains(Position{})).Each(
		func(id types.EntityID) bool {
			pos, err := cardinal.GetComponent[Position](world, id)
			if err != nil {
				return false
			}

			// Terminates the search if the player is found
			if pos.X == p.X && pos.Y == p.Y {
				eID = id
				return false
			}

			// Continue searching if the player is not the target player
			return true
		},
	)
	if searchErr != nil {
		return 0, err
	}
	if err != nil {
		return 0, err
	}

	return eID, nil
}
