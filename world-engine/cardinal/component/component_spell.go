package component

import "cinco-paus/seismic/client"

type Spell struct {
	Abilities [client.NumAbilities]int // Array of 5 integers
	Expired   bool
	Direction Direction
}

func (Spell) Name() string {
	return "Spell"
}
