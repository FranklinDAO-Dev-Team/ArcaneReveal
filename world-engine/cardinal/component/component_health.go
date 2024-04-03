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

func decrementHealth(world cardinal.WorldContext, entityID types.EntityID) error {
	health, err := cardinal.GetComponent[Health](world, entityID)
	if err != nil {
		return err
	}
	health.CurrHealth--
	err = cardinal.SetComponent[Health](world, entityID, health)
	if err != nil {
		return err
	}
	return nil
}
