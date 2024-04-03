package component

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

type Health struct {
	MaxHealth  int
	CurrHealth int
}

func (Health) Name() string {
	return "Health"
}

func DecrementHealth(world cardinal.WorldContext, entityID types.EntityID) (err error) {
	health, err := cardinal.GetComponent[Health](world, entityID)
	if err != nil {
		return err
	}
	health.CurrHealth--

	if health.CurrHealth <= 0 {
		// remove health entity
		err = cardinal.Remove(world, entityID)
		if err != nil {
			return err
		}
		fmt.Printf("entity %s died\n", fmt.Sprint(entityID))
	} else {
		// update health component
		err = cardinal.SetComponent[Health](world, entityID, health)
		if err != nil {
			return err
		}
	}

	return nil
}
