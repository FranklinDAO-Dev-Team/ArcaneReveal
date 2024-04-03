package component

type Spell struct {
	Abilities [NUM_ABILITIES]int // Array of 5 integers
	Expired   bool
	Direction Direction
}

func (Spell) Name() string {
	return "Spell"
}