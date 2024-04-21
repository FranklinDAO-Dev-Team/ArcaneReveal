package system

import (
	comp "cinco-paus/component"
	"cinco-paus/seismic/client"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

const PlayerMaxHealth = 5

func PopulateBoard(world cardinal.WorldContext, gameID types.EntityID) error {
	var err error

	// spawn player
	createPlayer(world, gameID, 1, 1)

	// create monsters
	createMonster(world, gameID, 9, 1, comp.MEDIUM)
	createMonster(world, gameID, 1, 9, comp.LIGHT)
	createMonster(world, gameID, 3, 9, comp.LIGHT)
	createMonster(world, gameID, 3, 3, comp.LIGHT)
	createMonster(world, gameID, 9, 7, comp.LIGHT)

	// other walls
	createWall(world, gameID, 3, 2)
	createWall(world, gameID, 5, 2)
	createWall(world, gameID, 2, 5)
	createWall(world, gameID, 2, 7)
	createWall(world, gameID, 4, 7)

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

	// create wands to track when abilities are spent
	SpawnWands(world, gameID)

	PrintStateToTerminal(world, gameID)

	return nil
}

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

func createMonster(world cardinal.WorldContext, gameId types.EntityID, x int, y int, monType comp.MonsterType) error {
	health := int(monType) + 1
	_, err := cardinal.Create(world,
		comp.Monster{Type: comp.HEAVY},
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

func SpawnWands(world cardinal.WorldContext, gameID types.EntityID) error {
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

// // deletes all monsters, eventually delete everything except player
// func clear_board(world cardinal.WorldContext) error {
// 	var err error
// 	searchErr := cardinal.NewSearch(
// 		world,
// 		filter.Contains(comp.Monster{}),
// 	).Each(func(id types.EntityID) bool {
// 		err = cardinal.Remove(world, id)
// 		return err == nil
// 	})
// 	return searchErr
// }
