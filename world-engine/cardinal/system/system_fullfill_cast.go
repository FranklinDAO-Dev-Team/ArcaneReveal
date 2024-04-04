package system

import (
	comp "cinco-paus/component"
	"cinco-paus/msg"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/iden3/go-iden3-crypto/poseidon"
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/message"
)

func FulfillCastSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage[msg.FulfillCastMsg, msg.FulfillCastMsgResult](
		world,
		func(turn message.TxData[msg.FulfillCastMsg]) (msg.FulfillCastMsgResult, error) {
			// debug prints
			fmt.Println("starting fulfill cast system")
			resultJSON, err := json.Marshal(turn.Msg.Result)
			if err != nil {
				fmt.Println("failed!!!!")
				return msg.FulfillCastMsgResult{}, fmt.Errorf("failed to marshal result to JSON: %v", err)
			}
			fmt.Printf("Result JSON: %s\n", resultJSON)

			// get relevant info about the cast
			spell, err := cardinal.GetComponent[comp.Spell](world, turn.Msg.Result.CastID)
			if err != nil {
				return msg.FulfillCastMsgResult{}, err
			}
			spell.Expired = false

			spellPos, err := cardinal.GetComponent[comp.Position](world, turn.Msg.Result.CastID)
			if err != nil {
				return msg.FulfillCastMsgResult{}, err
			}
			// remove cast entity
			err = cardinal.Remove(world, turn.Msg.Result.CastID)
			if err != nil {
				return msg.FulfillCastMsgResult{}, err
			}

			//  Check salts and that &turn.Msg.Abilities makes sense
			for i, canCastAbilityI := range turn.Msg.Result.Abilities {
				if canCastAbilityI {
					// fmt.Println("canCast", i)
					i64 := int64(i)
					salt := big.NewInt(0)
					salt.SetString(turn.Msg.Result.Salts[i], 10)
					commitment, error := poseidon.Hash([]*big.Int{big.NewInt(i64), salt})
					if error != nil {
						return msg.FulfillCastMsgResult{}, fmt.Errorf("failed to hash salt: %v", error)
					}

					game, err := cardinal.GetComponent[comp.Game](world, turn.Msg.Result.GameID)
					if err != nil {
						return msg.FulfillCastMsgResult{}, err
					}
					// fmt.Println("game.Commitments dimensions", len(*game.Commitments), len((*game.Commitments)[0]))
					// fmt.Printf("wand num: %d, i: %d\n", spell.WandNumber, i)
					if (*game.Commitments)[spell.WandNumber][0] != commitment.String() { // hardcoded to first index because wands have 1 ability rn
						return msg.FulfillCastMsgResult{}, fmt.Errorf("commitment %d does not match", i)
					}
				}
			}
			fmt.Println("Commitments verified")

			// resolve abilities and update chain state
			eventLogList := &[]comp.GameEventLog{}
			updateChainState := true
			err = resolveAbilities(world, spell, spellPos, &turn.Msg.Result.Abilities, updateChainState, eventLogList) // pass eventLogList to record executed resolutions
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

			// return successfully
			return msg.FulfillCastMsgResult{Success: true}, nil
		})
}
