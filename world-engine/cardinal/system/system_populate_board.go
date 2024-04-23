package system

import (
	comp "cinco-paus/component"
	"cinco-paus/seismic/client"
	"log"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

const PlayerMaxHealth = 5

// createPlayer creates a player entity for the given gameID at the given position
func createPlayer(world cardinal.WorldContext, gameID types.EntityID, x, y int) error {
	_, err := cardinal.Create(world,
		comp.Player{},
		comp.Collidable{Type: comp.PlayerCollide},
		comp.Health{
			MaxHealth:  PlayerMaxHealth,
			CurrHealth: PlayerMaxHealth,
		},
		comp.Position{
			X: x,
			Y: y,
		},
		comp.GameObj{GameID: gameID},
	)
	if err != nil {
		return err
	}

	return nil
}

// createWall creates a wall entity for the given gameID at the given position
func createWall(world cardinal.WorldContext, gameId types.EntityID, x, y int) error {
	_, err := cardinal.Create(world,
		comp.Wall{Type: comp.WALL},
		comp.Collidable{Type: comp.WallCollide},
		comp.Position{
			X: x,
			Y: y,
		},
		comp.GameObj{GameID: gameId},
	)
	if err != nil {
		return err
	}
	return nil
}

// createMonster creates a monster entity for the given gameID at the given position with the given type
// monster max health is set to monsterType + 1
func createMonster(world cardinal.WorldContext, gameId types.EntityID, x int, y int, monType comp.MonsterType) error {
	health := int(monType) + 1
	_, err := cardinal.Create(world,
		comp.Monster{Type: monType},
		comp.Collidable{Type: comp.MonsterCollide},
		comp.Health{
			MaxHealth:  health,
			CurrHealth: health,
		},
		comp.Position{
			X: x,
			Y: y,
		},
		comp.GameObj{GameID: gameId},
	)
	if err != nil {
		return err
	}
	return nil
}

// spawnWands creates NumWands wands for the given gameID
// to track when abilities are spent. Wands start as available
func spawnWands(world cardinal.WorldContext, gameID types.EntityID) error {
	for i := 0; i < client.NumWands; i++ {
		// w := comp.NewRandomWandCore()
		_, err := cardinal.Create(world,
			comp.WandCore{
				Number: i,
			},
			comp.Available{IsAvailable: true},
			comp.GameObj{GameID: gameID},
		)

		if err != nil {
			return err
		}
	}
	return nil
}

// spawnWallFrame creates the default wall structure for the given gameID
// this includes outer walls and dead squares that nothing should stand on
func spawnWallFrame(world cardinal.WorldContext, gameID types.EntityID) error {
	var err error
	// create outer walls and dead squares
	for i := 0; i < 11; i++ {
		for j := 0; j < 11; j++ {
			if i == 0 || i == 10 || j == 0 || j == 10 {
				err = createWall(world, gameID, i, j)
				if err != nil {
					return err
				}
			} else {
				if i%2 == 0 && j%2 == 0 {
					err = createWall(world, gameID, i, j)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// populateLevel1 populates the board for the first level
// this function differs from other populate functions
// because it also creates the player and wands
func populateLevel1(world cardinal.WorldContext, gameID types.EntityID) error {
	log.Println("populateLevel1()")
	// spawn player
	createPlayer(world, gameID, 1, 1)
	// create wands to track when abilities are spent
	spawnWands(world, gameID)

	populateDebugLevel(world, gameID)
	// // create walls
	// spawnWallFrame(world, gameID)
	// createWall(world, gameID, 3, 2)
	// createWall(world, gameID, 5, 2)
	// createWall(world, gameID, 2, 5)
	// createWall(world, gameID, 2, 7)
	// createWall(world, gameID, 4, 7)

	// // create monsters
	// // createMonster(world, gameID, 9, 1, comp.HEAVY)
	// // createMonster(world, gameID, 1, 9, comp.MEDIUM)
	// // createMonster(world, gameID, 3, 9, comp.LIGHT)
	// // createMonster(world, gameID, 3, 3, comp.LIGHT)
	// // createMonster(world, gameID, 9, 7, comp.LIGHT)
	// createMonster(world, gameID, 9, 1, comp.LIGHT)

	PrintStateToTerminal(world, gameID)

	return nil
}

func populateLevel2(world cardinal.WorldContext, gameID types.EntityID) error {
	spawnWallFrame(world, gameID)
	createWall(world, gameID, 6, 1)
	createWall(world, gameID, 2, 5)
	createWall(world, gameID, 3, 6)
	createWall(world, gameID, 2, 9)
	createWall(world, gameID, 9, 8)

	// create monsters
	createMonster(world, gameID, 3, 5, comp.MEDIUM)
	createMonster(world, gameID, 9, 7, comp.MEDIUM)
	createMonster(world, gameID, 3, 7, comp.LIGHT)
	createMonster(world, gameID, 7, 7, comp.LIGHT)
	createMonster(world, gameID, 1, 9, comp.LIGHT)

	PrintStateToTerminal(world, gameID)
	return nil
}

func populateDebugLevel(world cardinal.WorldContext, gameID types.EntityID) error {
	log.Print("populateDebugLevel()")
	// create walls

	// create walls
	spawnWallFrame(world, gameID)
	// createWall(world, gameID, 3, 2)
	// createWall(world, gameID, 5, 2)
	// createWall(world, gameID, 2, 5)
	// createWall(world, gameID, 2, 7)
	// createWall(world, gameID, 4, 7)

	// create monsters
	createMonster(world, gameID, 9, 1, comp.LIGHT)
	// createMonster(world, gameID, 1, 9, comp.MEDIUM)
	// createMonster(world, gameID, 3, 9, comp.LIGHT)
	// createMonster(world, gameID, 3, 3, comp.LIGHT)
	// createMonster(world, gameID, 9, 7, comp.LIGHT)
	// createMonster(world, gameID, 9, 1, comp.LIGHT)

	return nil
}
