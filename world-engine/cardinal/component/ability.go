package component

import (
	"pkg.world.dev/world-engine/cardinal"
)

type Ability interface {
	GetAbilityID() int
	Resolve(cardinal.WorldContext, *Position) (bool, error)
}
