package system

import (
	"errors"

	"pkg.world.dev/world-engine/cardinal"

	comp "cinco-paus/component"
)

const PLAYER_MAX_HEALTH = 5

// PlayerSpawnerSystem spawns players based on `CreatePlayer` transactions.
// This provides an example of a system that creates a new entity.
func PlayerSpawnerSystem(world cardinal.WorldContext) error {
	maxHp := 100
	_, err := cardinal.Create(world,
		comp.Player{
			Nickname:      "Spencer",
			MaxHealth:     PLAYER_MAX_HEALTH,
			CurrentHealth: PLAYER_MAX_HEALTH,
		},
		comp.Health{HP: maxHp},
	)
	if err != nil {
		return errors.New("failed to create player")
	}

	return nil
}
