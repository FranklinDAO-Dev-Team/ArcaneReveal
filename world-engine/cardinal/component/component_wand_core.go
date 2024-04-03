package component

import "math/rand"

const NUM_WANDS = 2
const NUM_ABILITIES = 1
const TOTAL_ABILITIES = 2

type WandCore struct {
	Number    int
	Abilities [NUM_ABILITIES]int // Array of 5 integers
	Revealed  [NUM_ABILITIES]int // Slice of integers
}

func (WandCore) Name() string {
	return "Wand"
}

func NewRandomWandCore() WandCore {
	w := WandCore{}

	// Set Revealed to all -1
	for i := range w.Revealed {
		w.Revealed[i] = -1
	}

	// Generate unique random numbers for Abilities
	uniqueNumbers := make(map[int]bool)
	for i := 0; i < NUM_ABILITIES; {
		num := rand.Intn(TOTAL_ABILITIES) + 1 // Random number between 1 and 50
		if !uniqueNumbers[num] {
			uniqueNumbers[num] = true
			w.Abilities[i] = num
			i++
		}
	}

	return w
}
