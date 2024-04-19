package component

import (
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
	} else {
		// update health component
		err = cardinal.SetComponent[Health](world, entityID, health)
		if err != nil {
			return err
		}
	}

	return nil
}

func IncrementHealth(world cardinal.WorldContext, entityID types.EntityID) (err error) {
	health, err := cardinal.GetComponent[Health](world, entityID)
	if err != nil {
		return err
	}

	if health.CurrHealth < health.MaxHealth {
		// Health is not at max, increment it
		health.CurrHealth++
		err = cardinal.SetComponent[Health](world, entityID, health)
		if err != nil {
			return err
		}
	}
	return nil
}
