package component

import (
	"log"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// polymorph
const Ability10ID = 10

type Ability10 struct{}

var _ Ability = &Ability10{}

func (Ability10) GetAbilityID() int {
	return Ability10ID
}
func (Ability10) GetAbilityName() string {
	return "polymorph"
}

// polymorphs the monster
func (Ability10) Resolve(
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

	if executeUpdates {
		err = polymorphMonster(world, gameID, id)
		if err != nil {
			log.Println("polymorphMonster() err: %w", err)
			return false, err
		}
		gameEvent := GameEventLog{X: spellPosition.X, Y: spellPosition.Y, Event: GameEventMonsterPolymorph}
		*eventLogList = append(*eventLogList, gameEvent)
	}

	// hit a monster, so ability should reveal
	return true, nil
}

func polymorphMonster(world cardinal.WorldContext, gameID types.EntityID, monID types.EntityID) error {
	log.Println("entered polymorphMonster()")
	// get monster type
	monster, err := cardinal.GetComponent[Monster](world, monID)
	if err != nil {
		return err
	}
	monsterPos, err := cardinal.GetComponent[Position](world, monID)
	if err != nil {
		return err
	}
	newMonsterType := (monster.Type + 1) % NumMonsterTypes
	log.Println("oldMonsterType: ", monster.Type)
	log.Println("newMonsterType: ", newMonsterType)
	// remove old monster
	err = cardinal.Remove(world, monID)
	if err != nil {
		return err
	}

	// create new monster
	newMonID, err := cardinal.Create(world,
		Monster{
			Type: newMonsterType,
		},
		Collidable{Type: MonsterCollide},
		Health{
			MaxHealth:  int(newMonsterType) + 1,
			CurrHealth: int(newMonsterType) + 1,
		},
		Position{
			X: monsterPos.X,
			Y: monsterPos.Y,
		},
		GameObj{GameID: gameID},
	)
	if err != nil {
		return err
	}
	log.Println("oldMonID: ", monID)
	log.Println("newMonID: ", newMonID)

	return nil

}
