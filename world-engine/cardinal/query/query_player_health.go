package query

type PlayerHealthRequest struct {
	Nickname string
}

type PlayerHealthResponse struct {
	HP int
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
