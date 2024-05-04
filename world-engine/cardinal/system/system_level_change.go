package system

import (
	comp "cinco-paus/component"
	"fmt"
	"log"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

const MaxLevels = 10

// LevelChangeSystem checks if the player has completed the current level
// and if so, switches to the next level.
// If the player beats the last level, the player wins the game.
func LevelChangeSystem(world cardinal.WorldContext) error {
	var outerErr error
	searchErr := cardinal.NewSearch(
		world,
		filter.Contains(comp.Game{})).
		Each(func(gameID types.EntityID) bool {
			levelCompleted, err := checkLevelCompleted(world, gameID)
			if err != nil {
				outerErr = err
				return false
			}

			if levelCompleted {
				log.Printf("LevelChangeSystem() game %d completed level\n", gameID)
				world.EmitEvent(map[string]any{
					"event":  "level-completed",
					"gameID": gameID,
				})
				switchLevels(world, gameID)
			}

			playerDied, err := checkPlayerDied(world, gameID)
			if err != nil {
				outerErr = err
				return false
			}

			if playerDied {
				log.Printf("LevelChangeSystem() game %d player died\n", gameID)
				world.EmitEvent(map[string]any{
					"event":  "game-over",
					"gameID": gameID,
				})
				err = removeGameInstance(world, gameID)
				if err != nil {
					log.Printf("LevelChangeSystem() failed to remove game instance %d: %w\n", gameID, err)
					return false
				}
				return false
			}

			// check next game
			return true
		})
	if searchErr != nil {
		return searchErr
	}
	if outerErr != nil {
		return outerErr
	}
	return nil
}

func checkPlayerDied(world cardinal.WorldContext, gameID types.EntityID) (bool, error) {
	playerID, err := comp.QueryPlayerID(world, gameID)
	if err != nil {
		return false, err
	}
	playerHealth, err := cardinal.GetComponent[comp.Health](world, playerID)
	if err != nil {
		return false, err
	}
	if playerHealth.CurrHealth == 0 {
		return true, nil
	}
	return false, nil
}

// checkLevelCompleted checks if the player has completed the current level
// A level is completed if all monsters have been killed
func checkLevelCompleted(world cardinal.WorldContext, gameID types.EntityID) (bool, error) {
	// iterate through all monsters on the shard,
	// and check if any are attacted to the given gameID
	var outerErr error
	var levelCompleted = true // level assumed to be completed until a monster is found
	searchErr := cardinal.NewSearch(
		world,
		filter.Contains(comp.Monster{})).
		Each(func(monID types.EntityID) bool {
			gameIDtag, err := cardinal.GetComponent[comp.GameObj](world, monID)
			if err != nil {
				outerErr = err
				return false
			}
			// monster exists for current game, level not completed
			if gameIDtag.GameID == gameID {
				levelCompleted = false
				return false
			}
			return true
		})
	if searchErr != nil {
		return false, searchErr
	}
	if outerErr != nil {
		return false, outerErr
	}
	return levelCompleted, nil
}

// switchLevels sets up the game for the next level by
// 1) incrementing game level, and checking if the player has won the game
// 2) removing all gameObjs attatched to the gameID except the player, then
// 1) moving the player back to the starting tile
// 2) reseting wand availability
// 3) populating the board with monsters and walls (and reset player position)
func switchLevels(world cardinal.WorldContext, gameID types.EntityID) error {
	var err error

	updatedLevel, err := updateGameLevel(world, gameID)
	if err != nil {
		return err
	}

	// clear board
	err = clearBoard(world, gameID)
	if err != nil {
		return err
	}

	// other stuff
	err = resetWandAvailability(world, gameID)
	if err != nil {
		return err
	}
	err = populateBoard(world, gameID, updatedLevel)
	if err != nil {
		return err
	}
	return nil

}

// updateGameLevel updates the game level for the given gameID
// if the player is already at the max level, emit a win event and delete the game
func updateGameLevel(world cardinal.WorldContext, gameID types.EntityID) (int, error) {
	game, err := cardinal.GetComponent[comp.Game](world, gameID)
	if err != nil {
		return -1, err
	}
	if game.Level == MaxLevels {
		log.Printf("updateGameLevel() game %d was won!", gameID)
		world.EmitEvent(map[string]any{
			"event":  "game-won",
			"gameID": gameID,
		})
		err = removeGameInstance(world, gameID)
		if err != nil {
			log.Printf("updateGameLevel() failed to remove game instance %d: %w\n", gameID, err)
			return -1, err
		}
		return -1, nil
	}

	// player is not at max level, increment level
	newLevel := game.Level + 1
	updatedGame := comp.Game{
		PersonaTag:  game.PersonaTag,
		Commitments: game.Commitments,
		Level:       newLevel,
	}
	err = cardinal.SetComponent[comp.Game](world, gameID, &updatedGame)
	if err != nil {
		return -1, err
	}
	world.EmitEvent(map[string]any{
		"event":    "level-won",
		"gameID":   gameID,
		"newLevel": newLevel,
	})

	return newLevel, nil
}

// clearBoard removes all gameObjs attatched to the gameID except the player
func clearBoard(world cardinal.WorldContext, gameID types.EntityID) error {
	var outerErr error
	searchErr := cardinal.NewSearch(
		world,
		filter.Contains(comp.Collidable{}),
	).Each(func(id types.EntityID) bool {
		// check if entity is attached to the current game
		// if not, skip to next entity
		gameIDtag, err := cardinal.GetComponent[comp.GameObj](world, id)
		if err != nil {
			outerErr = err
			return false
		}
		if gameIDtag.GameID != gameID {
			return true
		}

		// remove entity if it is a monster or wall
		// if so, remove it
		colType, err := cardinal.GetComponent[comp.Collidable](world, id)
		if err != nil {
			outerErr = err
			return false
		}
		if colType.Type == comp.MonsterCollide || colType.Type == comp.WallCollide {
			err = cardinal.Remove(world, id)
			if err != nil {
				outerErr = err
				return false
			}
		}

		// check next entity
		return true
	})
	if searchErr != nil {
		return searchErr
	}
	if outerErr != nil {
		return outerErr
	}
	return nil
}

func resetWandAvailability(world cardinal.WorldContext, gameID types.EntityID) error {
	var outerErr error
	searchErr := cardinal.NewSearch(
		world,
		filter.Contains(comp.WandCore{}),
	).Each(func(id types.EntityID) bool {
		// check if entity is attached to the current game
		// if not, skip to next entity
		gameIDtag, err := cardinal.GetComponent[comp.GameObj](world, id)
		if err != nil {
			outerErr = err
			return false
		}
		if gameIDtag.GameID != gameID {
			return true
		}

		// set availabile component to true
		err = cardinal.SetComponent[comp.Available](world, id, &comp.Available{IsAvailable: true})
		if err != nil {
			outerErr = err
			return false
		}

		// check next entity
		return true
	})
	if searchErr != nil {
		return searchErr
	}
	if outerErr != nil {
		return outerErr
	}
	return nil
}

// populateBoard populates the board with monsters and walls based on the input level
func populateBoard(world cardinal.WorldContext, gameID types.EntityID, level int) error {
	switch level {
	case 1:
		return populateLevel1(world, gameID)
	case 2:
		return populateLevel2(world, gameID)
	case 3:
		return populateLevel3(world, gameID)
	case 4:
		return populateLevel4(world, gameID)
	case 5:
		return populateLevel5(world, gameID)
	case 6:
		return populateLevel6(world, gameID)
	case 7:
		return populateLevel7(world, gameID)
	case 8:
		return populateLevel8(world, gameID)
	case 9:
		return populateLevel9(world, gameID)
	case 10:
		return populateLevel10(world, gameID)
	default:
		return fmt.Errorf("invalid level")
	}
}

func removeGameInstance(world cardinal.WorldContext, gameID types.EntityID) error {
	log.Printf("removeGameInstance() removing game instance %d", gameID)
	var outerErr error
	searchErr := cardinal.NewSearch(
		world,
		filter.Contains(comp.GameObj{}),
	).Each(func(id types.EntityID) bool {
		// check if entity is attached to the current game
		// if not, skip to next entity
		gameIDtag, err := cardinal.GetComponent[comp.GameObj](world, id)
		if err != nil {
			outerErr = err
			return false
		}
		if gameIDtag.GameID != gameID {
			return true
		}

		// remove entity
		err = cardinal.Remove(world, id)
		if err != nil {
			outerErr = err
			return false
		}

		// check next entity
		return true
	})
	if searchErr != nil {
		return searchErr
	}
	if outerErr != nil {
		return outerErr
	}

	// remove game itself
	err := cardinal.Remove(world, gameID)
	if err != nil {
		return err
	}

	// exit successfully
	return nil

}
