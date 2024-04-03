package system

// import (
// 	comp "cinco-paus/component"

// 	"pkg.world.dev/world-engine/cardinal"
// 	"pkg.world.dev/world-engine/cardinal/search/filter"
// 	"pkg.world.dev/world-engine/cardinal/types"
// )

// func MonsterTurnSystem(world cardinal.WorldContext) error {
// 	playerPos, err := cardinal.GetComponent[comp.Position](world, 0)
// 	if err != nil {
// 		return err
// 	}
// 	searchErr := cardinal.NewSearch(
// 		world,
// 		filter.Contains(comp.Monster{}),
// 	).Each(func(id types.EntityID) bool {
// 		monsterPos, err := cardinal.GetComponent[comp.Position](world, id)
// 		if err != nil {
// 			return false
// 		}
// 		// first try to attack, need chack manhatten distance function

// 		// if can't, get move options (places that are legal, not moving into a wall, etc),
// 		// need to update collidable tag first though to be "hard collide" or something, items don't get this (do I want though? some spells do collide with items)

// 		return true // always return true to move on to the next monster
// 	})

// 	if searchErr != nil {
// 		return searchErr
// 	}
// 	return nil
// }
