package component

import "fmt"

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
