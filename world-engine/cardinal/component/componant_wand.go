package component

import "math/rand"

const NUM_WANDS = 5
const NUM_ABILITIES = 5
const TOTAL_ABILITIES = 50

type Wand struct {
	Number    int
	Abilities [NUM_ABILITIES]int // Array of 5 integers
	Revealed  [NUM_ABILITIES]int // Slice of integers
	IsReady   [NUM_ABILITIES]bool
}

func (Wand) Name() string {
	return "Wand"
}

func (w *Wand) Init() error {
	// Set Revealed to all -1
	for i := range w.Revealed {
		w.Revealed[i] = -1
	}

	// Generate unique random numbers for Abilities
	uniqueNumbers := make(map[int]bool)
	for i := 0; i < NUM_ABILITIES; {
		num := rand.Intn(50) + 1 // Random number between 1 and 50
		if !uniqueNumbers[num] {
			uniqueNumbers[num] = true
			w.Abilities[i] = num
			i++
		}
	}

	// Set IsReady to true
	for i := range w.Revealed {
		w.IsReady[i] = true
	}

	return nil
}
