package component

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

const Ability9ID = 9

type Ability9 struct{}

var _ Ability = &Ability9{}

func (Ability9) GetAbilityID() int {
	return Ability9ID
}

// heals monster if it is below max health
func (a Ability9) Resolve(
	world cardinal.WorldContext,
	gameID types.EntityID,
	spellPosition *Position,
	_ Direction,
	executeUpdates bool,
	eventLogList *[]GameEventLog,
) (reveal bool, err error) {
	// Lookup if entity exists
	found, id, err := spellPosition.GetEntityIDByPosition(world, gameID)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}

	// check if its a monster
	colType, err := cardinal.GetComponent[Collidable](world, id)
	if err != nil {
		return false, err
	}
	if colType.Type != MonsterCollide {
		return false, nil
	}

	// check if monster is below max health
	monsterHealth, err := cardinal.GetComponent[Health](world, id)
	if err != nil {
		return false, err
	}
	if monsterHealth.CurrHealth == monsterHealth.MaxHealth {
		return false, err // ability cannot activate if monster is at max health
	}

	// Monster is below max health. If executeUpdates is true, heal it
	if executeUpdates {
		maxedHealth := Health{
			MaxHealth:  monsterHealth.MaxHealth,
			CurrHealth: monsterHealth.MaxHealth,
		}
		cardinal.SetComponent[Health](world, id, &maxedHealth)
		gameEvent := GameEventLog{X: spellPosition.X, Y: spellPosition.Y, Event: GameEventMonsterHeal}
		*eventLogList = append(*eventLogList, gameEvent)
	}

	// it's a monster below max health, so ability should reveal
	return true, nil
}
