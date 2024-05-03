package component

import (
	"errors"
	"fmt"
	"math"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

const MaxX = 10
const MaxY = 10

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
	switch direction {
	case LEFT:
		if p.X == -1 {
			return nil, fmt.Errorf("moving out of bounds: %s, %d", p.String(), int(direction))
		}
		p.X--
	case RIGHT:
		if p.X == MaxX {
			return nil, fmt.Errorf("moving out of bounds: %s, %d", p.String(), int(direction))
		}
		p.X++
	case UP:
		if p.Y == -1 {
			return nil, fmt.Errorf("moving out of bounds: %s, %d", p.String(), int(direction))
		}
		p.Y--
	case DOWN:
		if p.Y == MaxY {
			return nil, fmt.Errorf("moving out of bounds: %s, %d", p.String(), int(direction))
		}
		p.Y++
	default:
		return nil, fmt.Errorf("invalid direction")
	}

	return &p, nil
}

type EntityAtLocation struct {
	Monsters []types.EntityID
	Players  []types.EntityID
}

func (p *Position) GetEntityIDByPosition(
	world cardinal.WorldContext,
	gameID types.EntityID,
) (found bool, eID types.EntityID, searchErr error) {
	if p == nil {
		return false, 0, errors.New("attempting GetEntityIDByPosition with nil input")
	}
	searchErr = cardinal.NewSearch(
		world,
		filter.Contains(Position{}, GameObj{})).Each(
		func(id types.EntityID) bool {
			pos, err := cardinal.GetComponent[Position](world, id)
			if err != nil {
				return false
			}
			gameObjTag, err := cardinal.GetComponent[GameObj](world, id)
			if err != nil {
				return false
			}

			// Terminates the search if the entity is found
			if pos.X == p.X && pos.Y == p.Y && gameObjTag.GameID == gameID {
				eID = id
				found = true
				return false
			}

			// Continue searching if position (x, y) does not match
			return true
		},
	)
	if searchErr != nil {
		return false, 0, searchErr
	}

	return found, eID, nil
}

func (p *Position) ManhattenDistance(other *Position) int {
	return int(math.Abs(float64(p.X-other.X)) + math.Abs(float64(p.Y-other.Y)))
}

func (p *Position) Towards(other *Position) (Direction, error) {
	dx := other.X - p.X
	dy := other.Y - p.Y

	switch {
	case dx == 0 && dy < 0:
		return UP, nil
	case dx == 0 && dy > 0:
		return DOWN, nil
	case dx < 0 && dy == 0:
		return LEFT, nil
	case dx > 0 && dy == 0:
		return RIGHT, nil
	default:
		return -1, errors.New("other position is not in a single direction")
	}
}

func (d Direction) rotateClockwise() Direction {
	switch d {
	case UP:
		return RIGHT
	case RIGHT:
		return DOWN
	case DOWN:
		return LEFT
	case LEFT:
		return UP
	default:
		return -1
	}
}

func (p Position) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}
