package system

import (
	comp "cinco-paus/component"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

func getWandByNumber(world cardinal.WorldContext, targetNum int) (wandID types.EntityID, wandCore *comp.WandCore, available *comp.Available, err error) {
	searchErr := cardinal.NewSearch(world, filter.Contains(comp.WandCore{})).Each(
		func(id types.EntityID) bool {
			wandCore, err = cardinal.GetComponent[comp.WandCore](world, id)
			if err != nil {
				return false
			}

			// Terminates the search if the player is found
			if wandCore.Number == targetNum {
				wandID = id
				available, err = cardinal.GetComponent[comp.Available](world, id)
				if err != nil {
					return false
				}
				return false
			}

			// Continue searching if the player is not the target player
			return true
		},
	)
	if searchErr != nil {
		return 0, nil, nil, err
	}

	// log.Printf("wandID: %s\n", fmt.Sprint(wandID))
	// log.Println("wandCore: ", wandCore)
	return wandID, wandCore, available, nil
}
