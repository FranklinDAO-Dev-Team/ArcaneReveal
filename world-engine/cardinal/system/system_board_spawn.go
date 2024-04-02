package system

import (
	comp "cinco-paus/component"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
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
		comp.Monster{
			Type: "light",
		},
		comp.Health{
			MaxHealth:  1,
			CurrHealth: 1,
		},
		comp.Position{
			X: 0,
			Y: 1,
		},
	)

	_, err = cardinal.Create(world,
		comp.Monster{
			Type: "heavy",
		},
		comp.Health{
			MaxHealth:  3,
			CurrHealth: 3,
		},
		comp.Position{
			X: 0,
			Y: 2,
		},
	)

	return nil
}

// deletes all monsters, eventually delete everything except player
func clear_board(world cardinal.WorldContext) error {
	var err error
	searchErr := cardinal.NewSearch(
		world,
		filter.Contains(comp.Monster{}),
	).Each(func(id types.EntityID) bool {
		err = cardinal.Remove(world, id)
		if err != nil {
			return false
		}
		return true
	})
	return searchErr
}

// func PlayerHealth(world cardinal.WorldContext, req *PlayerHealthRequest) (*PlayerHealthResponse, error) {
// 	var playerHealth *comp.Health
// 	var err error
// 	searchErr := cardinal.NewSearch(
// 		world,
// 		filter.Exact(comp.Character{}, comp.Health{})).
// 		Each(func(id types.EntityID) bool {
// 			var character *comp.Character
// 			character, err = cardinal.GetComponent[comp.Character](world, id)
// 			if err != nil {
// 				return false
// 			}

// 			// Terminates the search if the player is found
// 			if player.Nickname == req.Nickname {
// 				playerHealth, err = cardinal.GetComponent[comp.Health](world, id)
// 				if err != nil {
// 					return false
// 				}
// 				return false
// 			}

// 			// Continue searching if the player is not the target player
// 			return true
// 		})
// 	if searchErr != nil {
// 		return nil, searchErr
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	if playerHealth == nil {
// 		return nil, fmt.Errorf("player %s does not exist", req.Nickname)
// 	}

// 	return &PlayerHealthResponse{HP: playerHealth.HP}, nil
// }
