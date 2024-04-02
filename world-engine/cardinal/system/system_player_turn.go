package system

import (
	comp "cinco-paus/component"
	"cinco-paus/msg"
	"fmt"

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

			playerID, err := queryPlayerID(world)
			player, err := cardinal.GetComponent[comp.Position](world, playerID)

			if err != nil {
				return msg.PlayerTurnResult{}, fmt.Errorf("PlayerTurnSystem err: %w", err)
			}
			if player == nil {
				return msg.PlayerTurnResult{}, fmt.Errorf("PlayerTurnSystem err: Player Not Found")
			}

			switch turn.Msg.Action {
			case "attack":
				// player_turn_attack(*player, turn.Msg.Direction)
			case "wand":
				// player_turn_wand(*player, turn.Msg.Direction)
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

func player_turn_attack(player comp.Player, direction string) {
	fmt.Printf("attacking in the %s direction\n", direction)
}

func player_turn_wand(player comp.Player, direction string) {
	fmt.Printf("Waving wand in the %s direction\n", direction)
}

func player_turn_move(world cardinal.WorldContext, direction string) error {
	playerID, err := queryPlayerID(world)
	if err != nil {
		return err
	}
	currPos, err := cardinal.GetComponent[comp.Position](world, playerID)

	switch direction {
	case "left":
		if currPos.X == 0 {
			return fmt.Errorf("moving out of bounds")
		}
		currPos.X--
	case "right":
		if currPos.X == 4 {
			return fmt.Errorf("moving out of bounds")
		}
		currPos.X++
	case "up":
		if currPos.Y == 0 {
			return fmt.Errorf("moving out of bounds")
		}
		currPos.Y--
	case "down":
		if currPos.Y == 4 {
			return fmt.Errorf("moving out of bounds")
		}
		currPos.Y++
	default:
		return fmt.Errorf("invalid direction")
	}
	cardinal.SetComponent[comp.Position](world, playerID, currPos)

	return nil
}
