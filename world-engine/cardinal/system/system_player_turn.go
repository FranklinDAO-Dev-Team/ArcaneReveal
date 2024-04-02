package system

import (
	comp "cinco-paus/component"
	"cinco-paus/msg"
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

			switch turn.Msg.Action {
			case "attack":
				player_turn_attack(world, turn.Msg.Direction)
			case "wand":
				wandnum, err := strconv.Atoi(turn.Msg.WandNum)
				if err != nil {
					return msg.PlayerTurnResult{}, fmt.Errorf("Error converting string to int: %w", err)
				}
				player_turn_wand(world, turn.Msg.Direction, wandnum)
			case "move":
				err = player_turn_move(world, turn.Msg.Direction)
				if err != nil {
					return msg.PlayerTurnResult{}, fmt.Errorf("PlayerTurnSystem err: %w", err)
				}
			default:
				return msg.PlayerTurnResult{}, fmt.Errorf("PlayerTurnSystem err: Invalid action")
			}

			err = world.EmitEvent(map[string]any{
				"event":     "player_turn",
				"action":    turn.Msg.Action,
				"direction": turn.Msg.Direction,
			})
			if err != nil {
				return msg.PlayerTurnResult{}, err
			}
			return msg.PlayerTurnResult{Success: true}, nil
		})
}

func player_turn_attack(world cardinal.WorldContext, direction string) error {
	fmt.Printf("attacking in the %s direction\n", direction)
	return nil
}

func player_turn_wand(world cardinal.WorldContext, direction string, wandnum int) error {
	playerPos, err := cardinal.GetComponent[comp.Position](world, 0)
	if err != nil {
		return err
	}
	spellPos, err := playerPos.GetUpdateFromDirection(direction)
	if err != nil {
		return err
	}
	spell, err := cardinal.Create(world,
		comp.Spell{},
		spellPos,
	)

	// not done, do stop it from erroring
	fmt.Println("spell: %d", spell)

	a1 := &comp.Ability_1{}
	fmt.Println(a1.GetAbilityID())
	a1.Resolve(world, spellPos)

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

func player_turn_move(world cardinal.WorldContext, direction string) error {
	playerID, err := queryPlayerID(world)
	if err != nil {
		return err
	}
	currPos, err := cardinal.GetComponent[comp.Position](world, playerID)
	updatePos, err := currPos.GetUpdateFromDirection(direction)
	cardinal.SetComponent[comp.Position](world, playerID, updatePos)

	return nil
}
