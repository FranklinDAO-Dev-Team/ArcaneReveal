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

	// spawn a monster
	_, err = cardinal.Create(world,
		comp.Monster{Type: comp.LIGHT},
		comp.Collidable{Type: comp.MonsterCollide},
		comp.Health{
			MaxHealth:  1,
			CurrHealth: 1,
		},
		comp.Position{
			X: 0,
			Y: 1,
		},
	)
	if err != nil {
		return err
	}

	_, err = cardinal.Create(world,
		comp.Monster{Type: comp.HEAVY},
		comp.Collidable{Type: comp.MonsterCollide},
		comp.Health{
			MaxHealth:  3,
			CurrHealth: 3,
		},
		comp.Position{
			X: 0,
			Y: 2,
		},
	)
	if err != nil {
		return err
	}

	_, err = cardinal.Create(world,
		comp.Monster{Type: comp.HEAVY},
		comp.Collidable{Type: comp.MonsterCollide},
		comp.Health{
			MaxHealth:  3,
			CurrHealth: 3,
		},
		comp.Position{
			X: 1,
			Y: 1,
		},
	)
	if err != nil {
		return err
	}

	_, err = cardinal.Create(world,
		comp.Wall{Type: comp.WALL},
		comp.Collidable{Type: comp.WallCollide},
		comp.Position{
			X: 3,
			Y: 0,
		},
	)
	if err != nil {
		return err
	}

	_, err = cardinal.Create(world,
		comp.Monster{Type: comp.HEAVY},
		comp.Collidable{Type: comp.MonsterCollide},
		comp.Health{
			MaxHealth:  3,
			CurrHealth: 3,
		},
		comp.Position{
			X: 4,
			Y: 0,
		},
	)
	if err != nil {
		return err
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
