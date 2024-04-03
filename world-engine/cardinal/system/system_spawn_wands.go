package system

import (
	comp "cinco-paus/component"

	"pkg.world.dev/world-engine/cardinal"
)

func SpawnWandsSystem(world cardinal.WorldContext) error {
	for i := 0; i < comp.NUM_WANDS; i++ {
		w := comp.NewRandomWandCore()
		_, err := cardinal.Create(world,
			comp.WandCore{
				Number:    i,
				Abilities: w.Abilities,
				Revealed:  w.Revealed,
			},
			comp.Available{IsAvailable: true},
		)

		if err != nil {
			return err
		}
	}
	return nil
}
