package system

import (
	"cinco-paus/component"
	comp "cinco-paus/component"
	"cinco-paus/msg"
	"cinco-paus/seismic/client"
	"errors"
	"fmt"
	"log"
	"strconv"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/message"
	"pkg.world.dev/world-engine/cardinal/types"
)

func PlayerTurnSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage[msg.PlayerTurnMsg, msg.PlayerTurnResult](
		world,
		func(turn message.TxData[msg.PlayerTurnMsg]) (msg.PlayerTurnResult, error) {
			var err error
			log.Println("starting player turn system")

			gameIdInt, err := strconv.Atoi(turn.Msg.GameIDStr)
			if err != nil {
				return msg.PlayerTurnResult{}, fmt.Errorf("error converting string to int: %w", err)
			}
			gameID := types.EntityID(gameIdInt)

			// check that the msg.sender is the game owner
			err = confirmGameOwnership(world, turn.Tx.PersonaTag, gameID)
			if err != nil {
				return msg.PlayerTurnResult{}, err
			}

			// check msg strings are well formatted
			err = turn.Msg.ValFmt()
			if err != nil {
				return msg.PlayerTurnResult{}, fmt.Errorf("error with msg format: %w", err)
			}

			direction, err := comp.StringToDirection(turn.Msg.Direction)
			if err != nil {
				return msg.PlayerTurnResult{}, err
			}
			eventLogList := &[]comp.GameEventLog{}

			err = playerTurnAction(world, gameID, turn, direction, eventLogList)
			if err != nil {
				return msg.PlayerTurnResult{}, err
			}

			// emit event log to client
			err = world.EmitEvent(map[string]any{
				"event":     "player_turn",
				"action":    turn.Msg.Action,
				"direction": direction,
			})
			if err != nil {
				return msg.PlayerTurnResult{}, err
			}

			// log
			PrintStateToTerminal(world, gameID)
			// debug prints
			// for _, logEntry := range *eventLogList {
			// 	log.Printf("X: %d, Y: %d, Event: %d\n",
			// 		logEntry.X, logEntry.Y, logEntry.Event)
			// }

			// return success
			result := msg.PlayerTurnResult{Success: true}
			return result, nil
		})
}

func playerTurnAction(
	world cardinal.WorldContext,
	gameID types.EntityID,
	turn message.TxData[msg.PlayerTurnMsg],
	direction comp.Direction,
	eventLogList *[]comp.GameEventLog,
) error {
	var err error
	switch turn.Msg.Action {
	case "attack":
		err = playerTurnAttack(world, gameID, direction, eventLogList)
		if err != nil {
			return err
		}
		err = MonsterTurnSystem(world, gameID, eventLogList)
		if err != nil {
			return err
		}

		// emit after attck and move
		eventMap := make(map[string]any)
		eventMap["turnEvent"] = *eventLogList
		err = world.EmitEvent(eventMap)
		if err != nil {
			return err
		}

	case "wand":
		wandnum, err := strconv.Atoi(turn.Msg.WandNum)
		if err != nil {
			return fmt.Errorf("error converting string to int: %w", err)
		}
		castID, potentialAbilities, err := playerTurnWand(world, gameID, direction, wandnum)
		if err != nil {
			return err
		}
		log.Printf("playerTurnAction potentialAbilities: %v \n", potentialAbilities)

		gameID, err := strconv.Atoi(turn.Msg.GameIDStr)
		if err != nil {
			return fmt.Errorf("error converting string to int: %w", err)
		}
		revealRequest := client.RevealRequest{
			PersonaTag:         turn.Tx.PersonaTag,
			GameID:             types.EntityID(gameID),
			CastID:             castID,
			WandNum:            wandnum,
			PotentialAbilities: *potentialAbilities,
		}

		// Send the reveal request to the Seismic server
		revealRequestCh <- revealRequest

	case "move":
		err = playerTurnMove(world, gameID, direction, eventLogList)
		if err != nil {
			return fmt.Errorf("PlayerTurnSystem err: %w", err)
		}
		err = MonsterTurnSystem(world, gameID, eventLogList)
		if err != nil {
			return fmt.Errorf("MonsterTurnSystem err: %w", err)
		}
		eventMap := make(map[string]any)
		eventMap["turnEvent"] = *eventLogList
		err = world.EmitEvent(eventMap)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("PlayerTurnSystem err: Invalid action")
	}
	return nil
}

func playerTurnAttack(
	world cardinal.WorldContext,
	gameID types.EntityID,
	direction comp.Direction,
	eventLogList *[]comp.GameEventLog,
) error {
	playerID, err := comp.QueryPlayerID(world, gameID)
	if err != nil {
		return err
	}
	playerPos, err := cardinal.GetComponent[comp.Position](world, playerID)
	if err != nil {
		return err
	}
	blankPos, err := playerPos.GetUpdateFromDirection(direction)
	if err != nil {
		return err
	}
	attackPos, err := blankPos.GetUpdateFromDirection(direction)
	if err != nil {
		return err
	}
	log.Printf("attackPos: %v\n", attackPos)

	found, id, err := attackPos.GetEntityIDByPosition(world, gameID)
	if err != nil {
		return err
	}
	if found {
		colType, err := cardinal.GetComponent[comp.Collidable](world, id)
		if err != nil {
			return err
		}
		log.Printf("colType: %v\n", colType)
		switch colType.Type {
		case comp.MonsterCollide:
			gameEvent := comp.GameEventLog{X: playerPos.X, Y: playerPos.Y, Event: comp.GameEventPlayerAttack}
			*eventLogList = append(*eventLogList, gameEvent)
			return comp.DecrementHealth(world, id)
		default:
			return fmt.Errorf("attempting to attack %s", colType.ToString())
		}
	} else {
		return errors.New("attempting to attack empty stace")
	}
}

func playerTurnWand(
	world cardinal.WorldContext,
	gameID types.EntityID,
	direction comp.Direction,
	wandnum int,
) (castID types.EntityID, potentialAbilities *[client.TotalAbilities]bool, err error) {
	log.Println("playerTurnWand")
	playerID, err := comp.QueryPlayerID(world, gameID)
	if err != nil {
		return 0, nil, err
	}
	playerPos, err := cardinal.GetComponent[comp.Position](world, playerID)
	if err != nil {
		return 0, nil, err
	}
	spellPos, err := playerPos.GetUpdateFromDirection(direction)
	if err != nil {
		return 0, nil, err
	}
	wandID, _, available, err := getWandByNumber(world, gameID, wandnum)
	if err != nil {
		return 0, nil, err
	}

	// handle wand availability
	if !available.IsAvailable {
		return 0, nil, fmt.Errorf("wand %d already expired", wandnum)
	}
	// set the wand to not ready (do early as it may potentially be refreshed by abilities)
	cardinal.SetComponent[comp.Available](world, wandID, &comp.Available{IsAvailable: false})

	// set all abilities to false becasue assume false until simulated
	allAbilities := &[client.TotalAbilities]bool{}
	spell := &comp.Spell{
		WandNumber: wandnum,
		Expired:    false,
		Abilities:  allAbilities,
		Direction:  direction,
	}

	// simulate a cast to determine potential ability activations
	// log.Printf("playerTurnWand potentialAbilities BEFORE resolveAbilities: %v \n", spell.Abilities)
	updateChainState := false
	dummy := &[]comp.GameEventLog{} // dummy event log, not used for anything but to satisfy the function signature
	err = resolveAbilities(world, gameID, spell, playerPos, spell.Abilities, updateChainState, dummy)
	if err != nil {
		return 0, nil, err
	}
	log.Printf("playerTurnWand potentialAbilities AFTER resolveAbilities: %v \n\n", spell.Abilities)
	// log.Printf("dummy event log: %v \n", dummy)

	// create a new entity for the cast to later be resolved
	castID, err = cardinal.Create(
		world,
		comp.GameObj{GameID: gameID},
		comp.AwaitingReveal{IsAvailable: true},
		spell,
		spellPos,
	)
	if err != nil {
		return 0, nil, err
	}

	return castID, spell.Abilities, nil
}

func playerTurnMove(
	world cardinal.WorldContext,
	gameID types.EntityID,
	direction comp.Direction,
	eventLogList *[]comp.GameEventLog,
) error {
	playerID, err := comp.QueryPlayerID(world, gameID)
	if err != nil {
		return err
	}
	playerPos, err := cardinal.GetComponent[comp.Position](world, playerID)
	if err != nil {
		return err
	}
	gameEvent := comp.GameEventLog{X: playerPos.X, Y: playerPos.Y, Event: directionToGameEventPlayerMove(direction)}
	*eventLogList = append(*eventLogList, gameEvent)

	newPos, err := playerPos.GetUpdateFromDirection(direction)
	if err != nil {
		return err
	}
	valid, err := comp.IsCollisonThere(world, gameID, *newPos)
	if err != nil {
		return err
	} else if valid {
		// invalid position, but don't return error, just check next direction
		return fmt.Errorf("would collide at %s", newPos.String())
	}
	thing := directionToGameEventPlayerMove(direction)
	*eventLogList = append(*eventLogList, comp.GameEventLog{X: newPos.X, Y: newPos.Y, Event: thing})

	newNewPos, err := newPos.GetUpdateFromDirection(direction)
	if err != nil {
		return err
	}
	valid, err = comp.IsCollisonThere(world, gameID, *newNewPos)
	if err != nil {
		return err
	} else if valid {
		// invalid position, but don't return error, just check next direction
		return fmt.Errorf("would collide at %s", newNewPos.String())
	}

	cardinal.SetComponent[comp.Position](world, playerID, newNewPos)

	return nil
}

func directionToGameEventPlayerMove(direction comp.Direction) comp.GameEvent {
	switch direction {
	case comp.LEFT:
		return comp.GameEventPlayerLeft
	case comp.RIGHT:
		return comp.GameEventPlayerRight
	case comp.UP:
		return comp.GameEventPlayerUp
	case comp.DOWN:
		return comp.GameEventPlayerDown
	default:
		panic("invalid direction")
	}
}

// checks that the given personaTag owns the game
func confirmGameOwnership(world cardinal.WorldContext, personaTag string, gameID types.EntityID) error {
	game, err := cardinal.GetComponent[component.Game](world, types.EntityID(gameID))
	if err != nil {
		return fmt.Errorf("failed to find game %d:", gameID)
	}
	if game.PersonaTag != personaTag {
		return fmt.Errorf("personaTag %s does not own game %d", personaTag, gameID)
	}

	return nil
}
