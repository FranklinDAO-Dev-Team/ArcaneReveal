package system

import (
	comp "cinco-paus/component"
	"cinco-paus/msg"
	"cinco-paus/seismic/client"
	"errors"
	"fmt"
	"strconv"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/message"
)

type gameEvent int

const (
	gameEventBeam gameEvent = iota
	gameEventDamage
	gameEventSpellDisappate
	gameEventMonsterAttack
	gameEventMonsterUp
	gameEventMonsterRight
	gameEventMonsterDown
	gameEventMonsterLeft
)

type gameEventLog struct {
	x     int
	y     int
	event gameEvent
}

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

			eventLogList := &[]gameEventLog{}

			switch turn.Msg.Action {
			case "attack":
				err = player_turn_attack(world, direction)
				if err != nil {
					return msg.PlayerTurnResult{Success: false}, err
				}
			case "wand":
				wandnum, err := strconv.Atoi(turn.Msg.WandNum)
				if err != nil {
					return msg.PlayerTurnResult{}, fmt.Errorf("error converting string to int: %w", err)
				}
				err = player_turn_wand(world, direction, wandnum, eventLogList)
				if err != nil {
					return msg.PlayerTurnResult{Success: false}, err
				}
			case "move":
				err = playerTurnMove(world, direction)
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

			result := msg.PlayerTurnResult{Success: true}

			MonsterTurnSystem(world, eventLogList)
			println()
			println("eventLogList: ", eventLogList)
			println("len(eventLogList): ", len(*eventLogList))
			for _, logEntry := range *eventLogList {
				fmt.Printf("X: %d, Y: %d, Event: %d\n",
					logEntry.x, logEntry.y, logEntry.event)
			}

			return result, nil
		})
}

func player_turn_attack(world cardinal.WorldContext, direction comp.Direction) error {
	playerPos, err := cardinal.GetComponent[comp.Position](world, 0)
	if err != nil {
		return err
	}
	attackPos, err := playerPos.GetUpdateFromDirection(direction)
	if err != nil {
		return err
	}

	found, id, err := attackPos.GetEntityIDByPosition(world)
	if err != nil {
		return err
	}
	if found {
		colType, err := cardinal.GetComponent[comp.Collidable](world, id)
		if err != nil {
			return err
		}
		switch colType.Type {
		case comp.MonsterCollide:
			return comp.DecrementHealth(world, id)
		default:
			return fmt.Errorf("attempting to attack %s", colType.ToString())
		}
	} else {
		return errors.New("attempting to attack empty stace")
	}
}

func player_turn_wand(world cardinal.WorldContext, direction comp.Direction, wandnum int, eventLogList *[]gameEventLog) error {
	playerPos, err := cardinal.GetComponent[comp.Position](world, 0)
	if err != nil {
		return err
	}
	spellPos, err := playerPos.GetUpdateFromDirection(direction)
	if err != nil {
		return err
	}

	wandID, wand, available, err := getWandByNumber(world, wandnum)
	if err != nil {
		return err
	}
	if !available.IsAvailable {
		return fmt.Errorf("wand %d already expired", wandnum)
	}
	// set the wand to not ready (do early as it may potentially be refreshed by abilities)
	cardinal.SetComponent[comp.Available](world, wandID, &comp.Available{IsAvailable: false})

	// hardcoding the ability for now instead of using wands
	spell := comp.Spell{
		Expired:   false,
		Abilities: wand.Abilities,
		Direction: direction,
	}

	potentialAbilities := &[client.TotalAbilities]bool{}
	updateChainState := false
	err = resolveAbilities(world, &spell, spellPos, potentialAbilities, updateChainState)
	if err != nil {
		return err
	}
	// TODO: call seismic client to resolve abilities
	seismic_response := &[client.TotalAbilities]bool{true, false}

	// acivate abilities returned by Seismic
	updateChainState = true
	spell.Expired = false
	err = resolveAbilities(world, &spell, spellPos, seismic_response, updateChainState)
	if err != nil {
		return err
	}
	// TODO: emit activated abilities and spell log to client

	return nil
}

func playerTurnMove(world cardinal.WorldContext, direction comp.Direction) error {
	playerID, err := queryPlayerID(world)
	if err != nil {
		return err
	}
	currPos, err := cardinal.GetComponent[comp.Position](world, playerID)
	if err != nil {
		return err
	}
	updatePos, err := currPos.GetUpdateFromDirection(direction)
	if err != nil {
		return err
	}

	found, id, err := updatePos.GetEntityIDByPosition(world)
	if err != nil {
		return err
	}
	if found {
		colType, err := cardinal.GetComponent[comp.Collidable](world, id)
		if err != nil {
			return err
		}
		if colType.Type != comp.ItemCollide {
			return fmt.Errorf("attempting to move onto an object of type %s", colType.ToString())
		}
	}

	cardinal.SetComponent[comp.Position](world, playerID, updatePos)

	return nil
}

func resolveAbilities(
	world cardinal.WorldContext,
	spell *comp.Spell,
	spellPos *comp.Position,
	potentialAbilities *[client.TotalAbilities]bool,
	updateChainState bool,
) error {
	for !spell.Expired {
		// record abilities that could activate a current square
		err := resolveAbilitiesAtPosition(world, spellPos, spell.Direction, potentialAbilities, updateChainState)
		if err != nil {
			return err
		}

		// get next spell position
		spellPos, err = spellPos.GetUpdateFromDirection(spell.Direction)
		if err != nil {
			spell.Expired = true
		}
		if spellPos == nil {
			spell.Expired = true
			break
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

func resolveAbilitiesAtPosition(
	world cardinal.WorldContext,
	spellPos *comp.Position,
	direction comp.Direction,
	potentialAbilities *[client.TotalAbilities]bool,
	updateChainState bool,
) error {
	for i := 0; i < len(*potentialAbilities); i++ {
		a := comp.AbilityMap[i+1]
		if a == nil {
			return errors.New("unknown ability called")
		}
		activated, err := a.Resolve(world, spellPos, direction, updateChainState)
		// only overwrite if ability activated
		(*potentialAbilities)[i] = activated || (*potentialAbilities)[i]
		if err != nil {
			return err
		}
	}

	return nil
}
