package system

import (
	comp "cinco-paus/component"
	"cinco-paus/msg"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/message"
)

func FulfillCastSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage[msg.FulfillCastMsg, msg.FulfillCastMsgResult](
		world,
		func(turn message.TxData[msg.FulfillCastMsg]) (msg.FulfillCastMsgResult, error) {
			// get relevant info about the cast
			spell, err := cardinal.GetComponent[comp.Spell](world, turn.Msg.CastID)
			if err != nil {
				return msg.FulfillCastMsgResult{}, err
			}
			spell.Expired = false

			spellPos, err := cardinal.GetComponent[comp.Position](world, turn.Msg.CastID)
			if err != nil {
				return msg.FulfillCastMsgResult{}, err
			}

			// TODO: Check salts and that &turn.Msg.Abilities makes sense
			println("TODO: Check salts and that &turn.Msg.Abilities makes sense")

			// resolve abilities and update chain state
			eventLogList := &[]comp.GameEventLog{}
			updateChainState := true
			err = resolveAbilities(world, spell, spellPos, &turn.Msg.Abilities, updateChainState, eventLogList) // pass eventLogList to record executed resolutions
			if err != nil {
				return msg.FulfillCastMsgResult{}, err
			}

			// Monster Turn occurs after abilities are resolved
			MonsterTurnSystem(world, eventLogList)

			// TODO: emit activated abilities and spell log to client
			println("TODO: emit activated abilities and spell log to client")
			for _, logEntry := range *eventLogList {
				fmt.Printf("X: %d, Y: %d, Event: %d\n",
					logEntry.X, logEntry.Y, logEntry.Event)
			}

			// remove cast entity
			err = cardinal.Remove(world, turn.Msg.CastID)
			if err != nil {
				return msg.FulfillCastMsgResult{}, err
			}

			// return successfully
			return msg.FulfillCastMsgResult{Success: true}, nil
		})
}
