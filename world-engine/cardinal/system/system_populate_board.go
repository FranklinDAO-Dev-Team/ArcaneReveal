package system

import (
	comp "cinco-paus/component"

	"pkg.world.dev/world-engine/cardinal"
)

func PopulateBoardSystem(world cardinal.WorldContext) error {
	// update player pos
	p_id, err := queryPlayerID(world)
	if err != nil {
		return err
	}
	err = cardinal.SetComponent[comp.Position](world, p_id, &comp.Position{X: 0, Y: 0})
	if err != nil {
		return err
	}

	// create outer walls and dead squares
	for i := 0; i < 11; i++ {
		for j := 0; j < 11; j++ {
			if i == 0 || i == 10 || j == 0 || j == 10 {
				err = createWall(world, i, j)
				if err != nil {
					return err
				}
			} else {
				if i%2 == 0 && j%2 == 0 {
					err = createWall(world, i, j)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	// other walls
	createWall(world, 3, 2)
	createWall(world, 5, 2)
	createWall(world, 2, 5)
	createWall(world, 2, 7)
	createWall(world, 4, 7)

	// create monsters
	createMonster(world, 1, 9, comp.LIGHT)
	createMonster(world, 9, 1, comp.LIGHT)
	createMonster(world, 7, 5, comp.LIGHT)
	createMonster(world, 5, 3, comp.MEDIUM)
	createMonster(world, 9, 7, comp.MEDIUM)

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

func createWall(world cardinal.WorldContext, x, y int) error {
	_, err := cardinal.Create(world,
		comp.Wall{Type: comp.WALL},
		comp.Collidable{Type: comp.WallCollide},
		comp.Position{
			X: x,
			Y: y,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func createMonster(world cardinal.WorldContext, x int, y int, monType comp.MonsterType) error {
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
	)
	if err != nil {
		return err
	}
	return nil
}
