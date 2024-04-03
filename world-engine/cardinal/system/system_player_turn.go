package system

import (
	comp "cinco-paus/component"
	"cinco-paus/msg"
	"errors"
	"fmt"
	"strconv"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/message"
)

func PlayerTurnSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage[msg.PlayerTurnMsg, msg.PlayerTurnResult](
		world,
		func(turn message.TxData[msg.PlayerTurnMsg]) (msg.PlayerTurnResult, error) {
			err := turn.Msg.ValFmt()
			if err != nil {
				return msg.PlayerTurnResult{}, fmt.Errorf("error with msg format: %w", err)
			}

			direction, err := comp.StringToDirection(turn.Msg.Direction)
			if err != nil {
				return msg.PlayerTurnResult{}, err
			}

			switch turn.Msg.Action {
			case "attack":
				player_turn_attack(world, direction)
			case "wand":
				wandnum, err := strconv.Atoi(turn.Msg.WandNum)
				if err != nil {
					return msg.PlayerTurnResult{}, fmt.Errorf("Error converting string to int: %w", err)
				}
				player_turn_wand(world, direction, wandnum)
			case "move":
				err = player_turn_move(world, direction)
				if err != nil {
					return msg.PlayerTurnResult{}, fmt.Errorf("PlayerTurnSystem err: %w", err)
				}
			default:
				return msg.PlayerTurnResult{}, fmt.Errorf("PlayerTurnSystem err: Invalid action")
			}

			err = world.EmitEvent(map[string]any{
				"event":     "player_turn",
				"action":    turn.Msg.Action,
				"direction": direction,
			})
			if err != nil {
				return msg.PlayerTurnResult{}, err
			}
			return msg.PlayerTurnResult{Success: true}, nil
		})
}

func player_turn_attack(world cardinal.WorldContext, direction comp.Direction) error {
	fmt.Printf("attacking in the %s direction\n", direction)
	return nil
}

func player_turn_wand(world cardinal.WorldContext, direction comp.Direction, wandnum int) error {
	playerPos, err := cardinal.GetComponent[comp.Position](world, 0)
	if err != nil {
		return err
	}
	spellPos, err := playerPos.GetUpdateFromDirection(direction)
	if err != nil {
		return err
	}

	// hardcoding the ability for now instead of using wands
	spell := comp.Spell{
		Expired:   false,
		Abilities: [1]int{2},
		Direction: direction,
	}
	spell_entity, err := cardinal.Create(world,
		spell,
		spellPos,
	)

	// not done, do stop it from erroring
	fmt.Println("spell: %d", spell_entity)
	// a1 := &comp.Ability_1{}
	// fmt.Println(a1.GetAbilityID())
	// a1.Resolve(world, spellPos)

	for !spell.Expired {
		fmt.Println("Spell postion: ", spellPos)
		for i := 0; i < len(spell.Abilities); i++ {
			// fmt.Printf("Resolving ability %d\n", spell.Abilities[i])
			a := comp.AbilityMap[spell.Abilities[i]]
			if a == nil {
				return errors.New("unknown ability called")
			}
			a.Resolve(world, spellPos, spell.Direction)
		}

		// get next spell position
		spellPos, err = spellPos.GetUpdateFromDirection(spell.Direction)
		if err != nil {
			spell.Expired = true
		}

		// if wall entity at spellPos, stop
		found, id, err := spellPos.GetEntityIDByPosition(world)
		if err != nil {
			return err
		}
		if found {
			colType, err := cardinal.GetComponent[comp.Collidable](world, id)
			if err != nil {
				return err
			}
			if colType.Type == comp.WallCollide {
				spell.Expired = true
			}
		}

	}

	return nil
}

// func create_spellhead(world cardinal.WorldContext, direction string, wandnum int) (spellhead, error) {
// 	playerID, err := queryPlayerID(world)
// 	if err != nil {
// 		return spellhead{}, err
// 	}
// 	pos, err := cardinal.GetComponent[comp.Position](world, playerID)
// 	if err != nil {
// 		return spellhead{}, err
// 	}
// 	_, wand, err := getWandByNumber(world, wandnum)
// 	if err != nil {
// 		return spellhead{}, err
// 	}
// 	pos.UpdateFromDirection(direction)

// 	var head = spellhead{
// 		Pos:       pos,
// 		Abilities: wand.Abilities,
// 	}

// 	return head, err
// }

func player_turn_move(world cardinal.WorldContext, direction comp.Direction) error {
	playerID, err := queryPlayerID(world)
	if err != nil {
		return err
	}
	currPos, err := cardinal.GetComponent[comp.Position](world, playerID)
	updatePos, err := currPos.GetUpdateFromDirection(direction)
	cardinal.SetComponent[comp.Position](world, playerID, updatePos)

	return nil
}
