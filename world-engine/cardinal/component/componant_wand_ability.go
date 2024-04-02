package component

import "pkg.world.dev/world-engine/cardinal"

type SpellAbility interface {
	GetAbilityID() int
	Resolve(cardinal.WorldContext, Position) (bool, error)
}
