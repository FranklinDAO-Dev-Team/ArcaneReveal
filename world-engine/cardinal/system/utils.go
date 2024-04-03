package system

import (
	comp "cinco-paus/component"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

func queryPlayerID(world cardinal.WorldContext) (types.EntityID, error) {
	return 0, nil // current hardcoded, player is always first entity created
}

func getWandByNumber(world cardinal.WorldContext, targetNum int) (types.EntityID, *comp.Wand, error) {
	var wandID types.EntityID
	var wand *comp.Wand
	var err error
	searchErr := cardinal.NewSearch(world, filter.Contains(comp.Wand{})).Each(
		func(id types.EntityID) bool {
			wand, err = cardinal.GetComponent[comp.Wand](world, id)
			if err != nil {
				return false
			}

			// Terminates the search if the player is found
			if wand.Number == targetNum {
				wandID = id
				return false
			}

			// Continue searching if the player is not the target player
			return true
		},
	)
	if searchErr != nil {
		return 0, nil, err
	}
	if err != nil {
		return 0, nil, err
	}

	return wandID, wand, nil
}
