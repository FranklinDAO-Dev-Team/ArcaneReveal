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

type Direction int

const (
	UP Direction = iota
	RIGHT
	DOWN
	LEFT
)

func (Position) Name() string {
	return "Position"
}

func StringToDirection(dirStr string) (Direction, error) {
	switch dirStr {
	case "left":
		return LEFT, nil
	case "right":
		return RIGHT, nil
	case "up":
		return UP, nil
	case "down":
		return DOWN, nil
	default:
		return -1, fmt.Errorf("invalid direction string %s", dirStr)
	}
}

func (p Position) GetUpdateFromDirection(direction Direction) (*Position, error) {
	// fmt.Printf("In UpdateFromDirection. p = (%d, %d), dir = %s\n", p.X, p.Y, direction)
	switch direction {
	case LEFT:
		if p.X == 0 {
			return nil, fmt.Errorf("moving out of bounds")
		}
		p.X--
	case RIGHT:
		if p.X == MAX_X {
			return nil, fmt.Errorf("moving out of bounds")
		}
		p.X++
	case UP:
		if p.Y == 0 {
			return nil, fmt.Errorf("moving out of bounds")
		}
		p.Y--
	case DOWN:
		if p.Y == MAX_Y {
			return nil, fmt.Errorf("moving out of bounds")
		}
		p.Y++
	default:
		return nil, fmt.Errorf("invalid direction")
	}

	// fmt.Println("Exiting UpdateFromDirection. p = ", p)
	return &p, nil
}

type EntityAtLocation struct {
	Monsters []types.EntityID
	Players  []types.EntityID
}

func (p *Position) getEntityIDByPosition(world cardinal.WorldContext) (types.EntityID, error) {
	var eID types.EntityID

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
		return 0, searchErr
	}

	return eID, nil
}
