package system

import (
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
			log.Println("starting player turn system")
			err := turn.Msg.ValFmt()
			if err != nil {
				return msg.PlayerTurnResult{}, fmt.Errorf("error with msg format: %w", err)
			}

			direction, err := comp.StringToDirection(turn.Msg.Direction)
			if err != nil {
				return msg.PlayerTurnResult{}, err
			}

			eventLogList := &[]comp.GameEventLog{}

			err = playerTurnAction(world, turn, direction, eventLogList)
			if err != nil {
				return msg.PlayerTurnResult{}, err
			}

			err = world.EmitEvent(map[string]any{
				"event":     "player_turn",
				"action":    turn.Msg.Action,
				"direction": direction,
			})
			PrintStateToTerminal(world)

			if err != nil {
				return msg.PlayerTurnResult{}, err
			}

			// debug prints
			// for _, logEntry := range *eventLogList {
			// 	log.Printf("X: %d, Y: %d, Event: %d\n",
			// 		logEntry.X, logEntry.Y, logEntry.Event)
			// }

			result := msg.PlayerTurnResult{Success: true}
			return result, nil
		})
}

func playerTurnAction(
	world cardinal.WorldContext,
	turn message.TxData[msg.PlayerTurnMsg],
	direction comp.Direction,
	eventLogList *[]comp.GameEventLog,
) error {
	var err error
	switch turn.Msg.Action {
	case "attack":
		err = playerTurnAttack(world, direction, eventLogList)
		if err != nil {
			return err
		}
		err = MonsterTurnSystem(world, eventLogList)
		if err != nil {
			return err
		}

		// TODO: emit events to client
		log.Println("TODO: emit activated abilities and spell log to client")

	case "wand":
		wandnum, err := strconv.Atoi(turn.Msg.WandNum)
		if err != nil {
			return fmt.Errorf("error converting string to int: %w", err)
		}
		castID, potentialAbilities, err := playerTurnWand(world, direction, wandnum)
		if err != nil {
			return err
		}
		// log.Println("castID: ", castID)
		// log.Println("potentialAbilities: ", potentialAbilities)

		// log.Println("gameidstr:", turn.Msg.GameIDStr)

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
		// set all abilities to true since we don't know which ones will be activated
		for i := 0; i < len(revealRequest.PotentialAbilities); i++ {
			revealRequest.PotentialAbilities[i] = true
		}
		revealRequestCh <- revealRequest
		log.Println("PlayerTurnSystem *potentialAbilities", revealRequest.PotentialAbilities)

	case "move":
		err = playerTurnMove(world, direction, eventLogList)
		if err != nil {
			return fmt.Errorf("PlayerTurnSystem err: %w", err)
		}
		err = MonsterTurnSystem(world, eventLogList)
		if err != nil {
			return fmt.Errorf("MonsterTurnSystem err: %w", err)
		}
		log.Println("TODO: emit activated abilities and spell log to client")
	default:
		return fmt.Errorf("PlayerTurnSystem err: Invalid action")
	}
	return nil
}

func playerTurnAttack(
	world cardinal.WorldContext,
	direction comp.Direction,
	eventLogList *[]comp.GameEventLog,
) error {
	playerPos, err := cardinal.GetComponent[comp.Position](world, 0)
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

	found, id, err := attackPos.GetEntityIDByPosition(world)
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
	direction comp.Direction,
	wandnum int,
) (castID types.EntityID, potentialAbilities *[client.TotalAbilities]bool, err error) {
	log.Println("playerTurnWand")
	playerPos, err := cardinal.GetComponent[comp.Position](world, 0)
	if err != nil {
		return 0, nil, err
	}
	// log.Println("playerTurnWand 1")
	spellPos, err := playerPos.GetUpdateFromDirection(direction)
	if err != nil {
		return 0, nil, err
	}
	// log.Println("playerTurnWand 2")
	wandID, _, available, err := getWandByNumber(world, wandnum)
	if err != nil {
		return 0, nil, err
	}
	// log.Println("playerTurnWand 3")

	// handle wand availability
	if !available.IsAvailable {
		return 0, nil, fmt.Errorf("wand %d already expired", wandnum)
	}
	// set the wand to not ready (do early as it may potentially be refreshed by abilities)
	cardinal.SetComponent[comp.Available](world, wandID, &comp.Available{IsAvailable: false})

	// set all abilities to true since we don't know which ones will be activated
	allAbilities := &[client.TotalAbilities]bool{}
	for i := range allAbilities {
		allAbilities[i] = true
	}
	spell := &comp.Spell{
		WandNumber: wandnum,
		Expired:    false,
		Abilities:  allAbilities,
		Direction:  direction,
	}
	// log.Println("playerTurnWand 4")

	// simulate a cast to determine potential ability activations
	updateChainState := false
	dummy := &[]comp.GameEventLog{} // dummy event log, not used for anything but to satisfy the function signature
	err = resolveAbilities(world, spell, playerPos, spell.Abilities, updateChainState, dummy)
	if err != nil {
		return 0, nil, err
	}
	log.Println("playerTurnWand 5")

	// create a new entity for the cast to later be resolved
	castID, err = cardinal.Create(
		world,
		comp.AwaitingReveal{IsAvailable: true},
		spell,
		spellPos,
	)
	if err != nil {
		return 0, nil, err
	}

	return castID, spell.Abilities, nil
}

func playerTurnMove(world cardinal.WorldContext, direction comp.Direction, eventLogList *[]comp.GameEventLog) error {
	playerID, err := comp.QueryPlayerID(world)
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
	valid, err := comp.IsCollisonThere(world, *newPos)
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
	valid, err = comp.IsCollisonThere(world, *newNewPos)
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
