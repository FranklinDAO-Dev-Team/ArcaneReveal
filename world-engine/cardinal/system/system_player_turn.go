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

			// player, err := queryPlayer(world, turn.Msg.Nickname)
			fmt.Println("in PlayerTurnSystem")
			player, err := cardinal.GetComponent[comp.Player](world, 0)

			if err != nil {
				return msg.PlayerTurnResult{}, fmt.Errorf("PlayerTurnSystem err: %w", err)
			}
			if player == nil {
				return msg.PlayerTurnResult{}, fmt.Errorf("PlayerTurnSystem err: Player Not Found")
			}

			fmt.Printf("turn.Msg.Action: %s, action: %s", turn.Msg.Action, turn.Msg.Action)

			switch turn.Msg.Action {
			case "attack":
				player_turn_attack(*player, turn.Msg.Direction)
			case "wand":
				player_turn_wand(*player, turn.Msg.Direction)
			case "move":
				fmt.Println("in correct switch")
				err = player_turn_move(world, player, turn.Msg.Direction)
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

func player_turn_move(world cardinal.WorldContext, player *comp.Player, direction string) error {
	fmt.Println("entered player_turn_move")
	switch direction {
	case "left":
		if player.X == 0 {
			return fmt.Errorf("moving out of bounds")
		}
		player.X--
	case "right":
		if player.X == 4 {
			return fmt.Errorf("moving out of bounds")
		}
		player.X++
	case "up":
		if player.Y == 0 {
			return fmt.Errorf("moving out of bounds")
		}
		player.Y--
	case "down":
		if player.Y == 4 {
			return fmt.Errorf("moving out of bounds")
		}
		player.Y++
	default:
		return fmt.Errorf("invalid direction")
	}
	fmt.Printf("x: %d, y: %d", player.X, player.Y)
	cardinal.SetComponent[comp.Player](world, 0, player)

	return nil
}
