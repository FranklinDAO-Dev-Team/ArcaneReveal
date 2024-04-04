package component

import "cinco-paus/seismic/client"

type Spell struct {
	WandNumber int
	Abilities  *[client.TotalAbilities]bool // Array of 5 integers
	Expired    bool
	Direction  Direction
}

func (Spell) Name() string {
	return "Spell"
}
