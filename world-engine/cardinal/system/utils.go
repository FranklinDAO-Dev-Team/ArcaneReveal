package system

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// import (
// 	"fmt"

// 	"pkg.world.dev/world-engine/cardinal"
// 	"pkg.world.dev/world-engine/cardinal/search/filter"
// 	"pkg.world.dev/world-engine/cardinal/types"

// 	comp "cinco-paus/component"
// )

// queryPlayer queries for the target player's info
func queryPlayerID(world cardinal.WorldContext) (types.EntityID, error) {
	return 0, nil // current hardcoded
}

// // queryTargetPlayer queries for the target player's entity ID and health component.
// func queryTargetPlayer(world cardinal.WorldContext, targetNickname string) (types.EntityID, *comp.Health, error) {
// 	var playerID types.EntityID
// 	var playerHealth *comp.Health
// 	var err error
// 	searchErr := cardinal.NewSearch(world, filter.Exact(comp.Player{}, comp.Health{})).Each(
// 		func(id types.EntityID) bool {
// 			var player *comp.Player
// 			player, err = cardinal.GetComponent[comp.Player](world, id)
// 			if err != nil {
// 				return false
// 			}

// 			// Terminates the search if the player is found
// 			if player.Nickname == targetNickname {
// 				playerID = id
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
// 		return 0, nil, err
// 	}
// 	if err != nil {
// 		return 0, nil, err
// 	}
// 	if playerHealth == nil {
// 		return 0, nil, fmt.Errorf("player %q does not exist", targetNickname)
// 	}

// 	return playerID, playerHealth, err
// }

// // queryPlayer queries for the target player's info
// func queryPlayer(world cardinal.WorldContext, nickname string) (*comp.Player, error) {
// 	fmt.Println("entered queryPlayer")
// 	var p *comp.Player
// 	var err error
// 	searchErr := cardinal.NewSearch(world, filter.Exact(comp.Player{})).Each(
// 		func(id types.EntityID) bool {
// 			var player *comp.Player
// 			player, err = cardinal.GetComponent[comp.Player](world, id)
// 			if err != nil {
// 				return false
// 			}

// 			fmt.Println(player)
// 			fmt.Printf("in queryPlayer search. nickname = %s, player.Nickname = %s\n", nickname, player.Nickname)

// 			// Terminates the search if the player is found
// 			if player.Nickname == nickname {
// 				p = player
// 				return false
// 			}

// 			// Continue searching if the player is not the target player
// 			return true
// 		})
// 	if searchErr != nil {
// 		return nil, err
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	return p, err
// }
