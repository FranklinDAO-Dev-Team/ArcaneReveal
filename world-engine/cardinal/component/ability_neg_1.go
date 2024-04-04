package component

import "pkg.world.dev/world-engine/cardinal"

const AbilityNegOneID = -1

type AbilityNegOne struct{}

var _ Ability = &AbilityNegOne{}

func (AbilityNegOne) GetAbilityID() int {
	return AbilityNegOneID
}

// a NOP ability for hidden abilities
func (a AbilityNegOne) Resolve(
	world cardinal.WorldContext,
	spellPosition *Position,
	_ Direction,
	executeUpdates bool,
	eventLogList *[]GameEventLog,
) (reveal bool, err error) {
	return false, nil
}
