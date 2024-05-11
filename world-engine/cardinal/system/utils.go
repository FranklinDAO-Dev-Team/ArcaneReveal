package system

import (
	comp "cinco-paus/component"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

func getWandByNumber(
	world cardinal.WorldContext,
	gameID types.EntityID,
	targetNum int,
) (wandID types.EntityID, wandCore *comp.WandCore, available *comp.Available, err error) {
	var outerErr error
	searchErr := cardinal.NewSearch(world, filter.Contains(comp.WandCore{}, comp.GameObj{})).Each(
		func(id types.EntityID) bool {
			// make sure the wand is for the right game
			gameObjTag, err := cardinal.GetComponent[comp.GameObj](world, id)
			if err != nil {
				outerErr = err
				return false
			}
			if gameObjTag.GameID != gameID {
				// skip to next entity
				return true
			}

			wandCore, err = cardinal.GetComponent[comp.WandCore](world, id)
			if err != nil {
				outerErr = err
				return false
			}
			// Terminates the search if the wand is found
			if wandCore.Number == targetNum {
				// set return values
				wandID = id
				available, err = cardinal.GetComponent[comp.Available](world, id)
				if err != nil {
					outerErr = err
					return false
				}

				// successfully found the wand, stop searching
				return false
			}

			// Continue searching if the wand num is wrong
			return true
		},
	)
	if searchErr != nil {
		return 0, nil, nil, err
	}
	if outerErr != nil {
		return 0, nil, nil, outerErr
	}

	return wandID, wandCore, available, nil
}
