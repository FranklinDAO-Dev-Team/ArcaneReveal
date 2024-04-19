package system

import (
	"errors"
	"log"

	"pkg.world.dev/world-engine/cardinal"

	comp "cinco-paus/component"
)

const PlayerMaxHealth = 5

// PlayerSpawnerSystem spawns players based on `CreatePlayer` transactions.
// This provides an example of a system that creates a new entity.
func SpawnPlayerSystem(world cardinal.WorldContext) error {
	// Create player
	_, err := cardinal.Create(world,
		comp.Player{},
		comp.Collidable{Type: comp.PlayerCollide},
		comp.Health{
			MaxHealth:  PlayerMaxHealth,
			CurrHealth: PlayerMaxHealth,
		},
		comp.Position{
			X: 1,
			Y: 1,
		},
	)
	log.Println("player spawned")
	p, _ := cardinal.GetComponent[comp.Position](world, 0)
	log.Println(p.X, p.Y)

	if err != nil {
		return errors.New("failed to create player")
	}

	return nil
}
