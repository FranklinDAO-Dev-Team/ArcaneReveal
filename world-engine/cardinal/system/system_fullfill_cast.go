package system

// import (
// 	comp "cinco-paus/component"
// 	"cinco-paus/msg"

// 	"pkg.world.dev/world-engine/cardinal"
// 	"pkg.world.dev/world-engine/cardinal/message"
// )

// func FulfillCastSystem(world cardinal.WorldContext) error {
// 	return cardinal.EachMessage[msg.PlayerTurnMsg, msg.PlayerTurnResult](
// 		world,
// 		func(turn message.TxData[msg.PlayerTurnMsg]) (msg.PlayerTurnResult, error) {

// 			// acivate abilities returned by Seismic
// 			spell := comp.Spell{
// 				Expired:   false,
// 				Abilities: wand.Abilities,
// 				Direction: direction,
// 			}
// 			updateChainState := true

// 			// resolve abilities and update chain state
// 			err := resolveAbilities(world, &spell, spellPos, seismic_response, updateChainState, eventLogList) // pass eventLogList to record executed resolutions
// 			if err != nil {
// 				return msg.PlayerTurnResult{Success: false}, err
// 			}

// 			// TODO: emit activated abilities and spell log to client
// 			// MonsterTurnSystem(world, eventLogList)
// 			// println()
// 			// println("eventLogList: ", eventLogList)
// 			// println("len(eventLogList): ", len(*eventLogList))
// 			// for _, logEntry := range *eventLogList {
// 			// 	fmt.Printf("X: %d, Y: %d, Event: %d\n",
// 			// 		logEntry.X, logEntry.Y, logEntry.Event)
// 			// }

// 			// check salts

// 			panic("unimplemented")
// 			// return successfully
// 			return msg.PlayerTurnResult{Success: true}, nil
// 		})
// }
