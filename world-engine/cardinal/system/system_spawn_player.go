package system

import (
	"errors"

	"pkg.world.dev/world-engine/cardinal"

	comp "cinco-paus/component"
)

const PLAYER_MAX_HEALTH = 5

// PlayerSpawnerSystem spawns players based on `CreatePlayer` transactions.
// This provides an example of a system that creates a new entity.
func SpawnPlayerSystem(world cardinal.WorldContext) error {
	// Create player
	_, err := cardinal.Create(world,
		comp.Player{},
		comp.Health{
			MaxHealth:  PLAYER_MAX_HEALTH,
			CurrHealth: PLAYER_MAX_HEALTH,
		},
		comp.Position{
			X: 0,
			Y: 0,
		},
	)

	if err != nil {
		return errors.New("failed to create player")
	}

	return nil
}
