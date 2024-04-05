package component

import (
	"pkg.world.dev/world-engine/cardinal"
)

const Ability3ID = 3

type Ability3 struct{}

var _ Ability = &Ability3{}

func (Ability3) GetAbilityID() int {
	return Ability3ID
}

func (a Ability3) Resolve(
	world cardinal.WorldContext,
	spellPosition *Position,
	direction Direction,
	executeUpdates bool,
	eventLogList *[]GameEventLog,
) (reveal bool, err error) {
	found, id, err := spellPosition.GetEntityIDByPosition(world)
	if err != nil {
		return false, err
	}
	if found {
		colType, err := cardinal.GetComponent[Collidable](world, id)
		if err != nil {
			return false, err
		}
		if colType.Type == WallCollide {
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					damageDealt, err := damageAtPostion(world, spellPosition, executeUpdates, false)
					if err != nil {
						return false, err
					}
					if damageDealt {
						*eventLogList = append(*eventLogList, GameEventLog{X: spellPosition.X, Y: spellPosition.Y, Event: GameEventSpellDamage})
					}
					reveal = true
				}
			}
		}
	}
	return reveal, nil
}
