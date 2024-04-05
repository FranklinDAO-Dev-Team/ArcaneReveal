package component

import "pkg.world.dev/world-engine/cardinal"

const Ability4ID = 4

type Ability4 struct{}

var _ Ability = &Ability4{}

func (Ability4) GetAbilityID() int {
	return Ability4ID
}

func (a Ability4) Resolve(
	world cardinal.WorldContext,
	spellPosition *Position,
	direction Direction,
	executeUpdates bool,
	eventLogList *[]GameEventLog,
) (reveal bool, err error) {
	// look up entity at spell position
	found, id, err := spellPosition.GetEntityIDByPosition(world)
	if err != nil {
		return false, err
	}
	if found {
		colType, err := cardinal.GetComponent[Collidable](world, id)
		if err != nil {
			return false, err
		}
		// if entity is a wall, then trigger explosion
		if colType.Type == WallCollide {
			topLeft := Position{X: spellPosition.X - 2, Y: spellPosition.Y - 2}
			for i := 0; i < 5; i++ {
				for j := 0; j < 5; j++ {
					damagePos := Position{X: topLeft.X + i, Y: topLeft.Y + j}
					damageDealt, err := damageAtPostion(world, &damagePos, executeUpdates, true)
					if err != nil {
						return false, err
					}
					if damageDealt {
						*eventLogList = append(*eventLogList, GameEventLog{X: damagePos.X, Y: damagePos.Y, Event: GameEventSpellDamage})
					}
					reveal = true
				}
			}
		}
	}
	return reveal, nil
}
