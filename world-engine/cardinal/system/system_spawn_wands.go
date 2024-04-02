package system

import (
	comp "cinco-paus/component"
	"math/rand"

	"pkg.world.dev/world-engine/cardinal"
)

func SpawnWandsSystem(world cardinal.WorldContext) error {
	for i := 0; i < comp.NUM_WANDS; i++ {
		w := generateRandomWand()
		_, err := cardinal.Create(world,
			// abilites, err := generateWandAbilites(comp.NUM_WANDS, 0, 0)
			comp.Wand{
				Number:    i,
				Abilities: w.Abilities,
				Revealed:  w.Revealed,
				IsReady:   w.IsReady,
			},
		)

		if err != nil {
			return err
		}
	}
	return nil
}

// generates a random want for the start of a game
// want number defaults to 0
// abilites in range(0, TOTAL_ABILITIES)
// revealed is array of -1
// isReady is array of true
func generateRandomWand() comp.Wand {
	var w = comp.Wand{}
	// Set Revealed to all -1
	for i := range w.Revealed {
		w.Revealed[i] = -1
	}

	// Generate unique random numbers for Abilities
	uniqueNumbers := make(map[int]bool)
	for i := 0; i < comp.NUM_ABILITIES; {
		num := rand.Intn(comp.TOTAL_ABILITIES) // Random number between 0 and 49
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

	return w
}
