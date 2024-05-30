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

// setPlayerPosition gets the game's player position and sets it to (x, y)
func setPlayerPosition(world cardinal.WorldContext, gameID types.EntityID, x, y int) error {
	playerID, err := comp.QueryPlayerID(world, gameID)
	if err != nil {
		return err
	}
	err = cardinal.SetComponent[comp.Position](world, playerID, &comp.Position{x, y})
	return err
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

	// if debug:
	// populateDebugLevel(world, gameID)

	// create walls
	spawnWallFrame(world, gameID)
	createWall(world, gameID, 3, 2)
	createWall(world, gameID, 5, 2)
	createWall(world, gameID, 2, 5)
	createWall(world, gameID, 2, 7)
	createWall(world, gameID, 4, 7)

	// create monsters
	createMonster(world, gameID, 1, 9, comp.MEDIUM)
	createMonster(world, gameID, 3, 3, comp.LIGHT)
	createMonster(world, gameID, 3, 9, comp.LIGHT)
	createMonster(world, gameID, 9, 1, comp.LIGHT)

	PrintStateToTerminal(world, gameID)

	return nil
}

func populateLevel2(world cardinal.WorldContext, gameID types.EntityID) error {
	log.Printf("game %d populating level 2\n", gameID)
	setPlayerPosition(world, gameID, 1, 9)
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
	// createMonster(world, gameID, 1, 9, comp.LIGHT)

	PrintStateToTerminal(world, gameID)
	return nil
}

func populateLevel3(world cardinal.WorldContext, gameID types.EntityID) error {
	setPlayerPosition(world, gameID, 9, 9)
	spawnWallFrame(world, gameID)
	createWall(world, gameID, 6, 1)
	createWall(world, gameID, 8, 7)
	createWall(world, gameID, 8, 5)
	createWall(world, gameID, 8, 1)
	createWall(world, gameID, 4, 1)
	createWall(world, gameID, 4, 5)
	createWall(world, gameID, 4, 7)
	createWall(world, gameID, 4, 9)
	createWall(world, gameID, 2, 5)

	// create monsters
	createMonster(world, gameID, 5, 9, comp.MEDIUM)
	createMonster(world, gameID, 3, 9, comp.MEDIUM)
	createMonster(world, gameID, 1, 9, comp.MEDIUM)
	createMonster(world, gameID, 9, 5, comp.LIGHT)
	createMonster(world, gameID, 1, 3, comp.LIGHT)

	PrintStateToTerminal(world, gameID)
	return nil
}

func populateLevel4(world cardinal.WorldContext, gameID types.EntityID) error {
	setPlayerPosition(world, gameID, 1, 9)
	spawnWallFrame(world, gameID)
	createWall(world, gameID, 2, 1)
	createWall(world, gameID, 4, 1)
	createWall(world, gameID, 8, 1)
	createWall(world, gameID, 7, 2)
	createWall(world, gameID, 7, 4)
	createWall(world, gameID, 2, 5)
	createWall(world, gameID, 4, 5)
	createWall(world, gameID, 6, 5)
	createWall(world, gameID, 7, 8)
	createWall(world, gameID, 9, 8)

	// create monsters
	createMonster(world, gameID, 1, 3, comp.MEDIUM)
	createMonster(world, gameID, 5, 7, comp.MEDIUM)
	createMonster(world, gameID, 3, 1, comp.LIGHT)
	createMonster(world, gameID, 5, 3, comp.LIGHT)
	createMonster(world, gameID, 9, 5, comp.LIGHT)

	PrintStateToTerminal(world, gameID)
	return nil
}

func populateLevel5(world cardinal.WorldContext, gameID types.EntityID) error {
	setPlayerPosition(world, gameID, 5, 1)

	// create walls
	spawnWallFrame(world, gameID)
	createWall(world, gameID, 2, 3)
	createWall(world, gameID, 4, 3)
	createWall(world, gameID, 1, 6)
	createWall(world, gameID, 2, 9)
	createWall(world, gameID, 4, 7)

	// create monsters
	createMonster(world, gameID, 1, 7, comp.MEDIUM)
	createMonster(world, gameID, 5, 9, comp.MEDIUM)
	createMonster(world, gameID, 3, 9, comp.LIGHT)
	createMonster(world, gameID, 9, 1, comp.LIGHT)
	createMonster(world, gameID, 1, 5, comp.LIGHT)

	PrintStateToTerminal(world, gameID)
	return nil
}

func populateLevel6(world cardinal.WorldContext, gameID types.EntityID) error {
	setPlayerPosition(world, gameID, 9, 5)

	// create walls
	spawnWallFrame(world, gameID)
	createWall(world, gameID, 8, 5)
	createWall(world, gameID, 8, 1)
	createWall(world, gameID, 9, 8)
	createWall(world, gameID, 5, 8)
	createWall(world, gameID, 4, 8)
	createWall(world, gameID, 5, 6)
	createWall(world, gameID, 3, 6)
	createWall(world, gameID, 1, 4)
	createWall(world, gameID, 4, 5)
	createWall(world, gameID, 4, 3)

	// create monsters
	createMonster(world, gameID, 3, 1, comp.MEDIUM)
	createMonster(world, gameID, 1, 7, comp.MEDIUM)
	createMonster(world, gameID, 3, 7, comp.MEDIUM)
	createMonster(world, gameID, 7, 5, comp.LIGHT)
	createMonster(world, gameID, 7, 3, comp.LIGHT)

	PrintStateToTerminal(world, gameID)
	return nil
}

func populateLevel7(world cardinal.WorldContext, gameID types.EntityID) error {
	setPlayerPosition(world, gameID, 5, 9)

	// create walls
	spawnWallFrame(world, gameID)
	createWall(world, gameID, 2, 1)
	createWall(world, gameID, 5, 2)
	createWall(world, gameID, 7, 2)
	createWall(world, gameID, 1, 4)
	createWall(world, gameID, 3, 4)
	createWall(world, gameID, 4, 5)
	createWall(world, gameID, 4, 7)
	createWall(world, gameID, 6, 5)
	createWall(world, gameID, 2, 9)
	createWall(world, gameID, 9, 6)
	createWall(world, gameID, 7, 8)

	// create monsters
	createMonster(world, gameID, 7, 1, comp.MEDIUM)
	createMonster(world, gameID, 3, 3, comp.MEDIUM)
	createMonster(world, gameID, 1, 9, comp.MEDIUM)
	createMonster(world, gameID, 1, 3, comp.LIGHT)

	return nil
}

func populateLevel8(world cardinal.WorldContext, gameID types.EntityID) error {
	setPlayerPosition(world, gameID, 1, 5)

	// create walls
	spawnWallFrame(world, gameID)
	createWall(world, gameID, 1, 2)
	createWall(world, gameID, 3, 2)
	createWall(world, gameID, 2, 3)
	createWall(world, gameID, 3, 4)
	createWall(world, gameID, 5, 6)
	createWall(world, gameID, 4, 7)
	createWall(world, gameID, 6, 5)
	createWall(world, gameID, 7, 4)

	// create monsters
	createMonster(world, gameID, 9, 9, comp.MEDIUM)
	createMonster(world, gameID, 5, 9, comp.MEDIUM)
	createMonster(world, gameID, 3, 7, comp.LIGHT)
	createMonster(world, gameID, 7, 3, comp.LIGHT)

	return nil
}

func populateLevel9(world cardinal.WorldContext, gameID types.EntityID) error {
	setPlayerPosition(world, gameID, 5, 5)

	spawnWallFrame(world, gameID)
	createWall(world, gameID, 1, 2)
	createWall(world, gameID, 4, 1)
	createWall(world, gameID, 5, 4)
	createWall(world, gameID, 7, 4)
	createWall(world, gameID, 4, 7)
	createWall(world, gameID, 4, 9)
	createWall(world, gameID, 6, 7)
	createWall(world, gameID, 7, 8)

	// create monsters
	createMonster(world, gameID, 3, 9, comp.MEDIUM)
	createMonster(world, gameID, 5, 9, comp.MEDIUM)
	createMonster(world, gameID, 7, 1, comp.MEDIUM)
	createMonster(world, gameID, 3, 9, comp.MEDIUM)
	createMonster(world, gameID, 7, 9, comp.LIGHT)

	PrintStateToTerminal(world, gameID)
	return nil
}

func populateLevel10(world cardinal.WorldContext, gameID types.EntityID) error {
	setPlayerPosition(world, gameID, 5, 5)

	spawnWallFrame(world, gameID)
	createWall(world, gameID, 4, 3)
	createWall(world, gameID, 5, 2)
	createWall(world, gameID, 7, 2)
	createWall(world, gameID, 8, 3)
	createWall(world, gameID, 7, 4)
	createWall(world, gameID, 6, 5)
	createWall(world, gameID, 5, 6)
	createWall(world, gameID, 4, 7)
	createWall(world, gameID, 3, 8)
	createWall(world, gameID, 5, 8)
	createWall(world, gameID, 7, 8)
	createWall(world, gameID, 9, 6)

	// create monsters
	createMonster(world, gameID, 1, 3, comp.MEDIUM)
	createMonster(world, gameID, 5, 1, comp.MEDIUM)
	createMonster(world, gameID, 9, 3, comp.MEDIUM)
	createMonster(world, gameID, 3, 9, comp.LIGHT)
	createMonster(world, gameID, 1, 5, comp.LIGHT)

	PrintStateToTerminal(world, gameID)
	return nil
}

// populateDebugLevel is used for debugging purposes only
func populateDebugLevel(world cardinal.WorldContext, gameID types.EntityID) error {
	log.Print("populateDebugLevel()")
	setPlayerPosition(world, gameID, 1, 1)
	spawnWallFrame(world, gameID)
	createMonster(world, gameID, 3, 1, comp.LIGHT)

	return nil
}
