package component

type Spell struct {
	Abilities [NUM_ABILITIES]Ability // Array of 5 integers
	Expired   bool
	Direction string
}

func (Spell) Name() string {
	return "Spell"
}
