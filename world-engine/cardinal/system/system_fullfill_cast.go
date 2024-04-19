package system

import (
	comp "cinco-paus/component"
	"cinco-paus/msg"
	"cinco-paus/seismic/client"
	"encoding/json"
	"fmt"
	"log"

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
			log.Println("starting fulfill cast system")
			resultJSON, err := json.Marshal(turn.Msg.Result)
			if err != nil {
				// log.Println("failed!!!!")
				return msg.FulfillCastMsgResult{}, fmt.Errorf("failed to marshal result to JSON: %v", err)
			}
			log.Printf("Result JSON: %s\n", resultJSON)

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
					// log.Println("canCast", i)
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
					// log.Println("game.Commitments dimensions", len(*game.Commitments), len((*game.Commitments)[0]))
					// log.Printf("wand num: %d, i: %d\n", spell.WandNumber, i)

					if !contains((*game.Commitments)[spell.WandNumber], commitment.String()) { // hardcoded to first index because wands have 1 ability rn
						return msg.FulfillCastMsgResult{}, fmt.Errorf("commitment %d does not match", i)
					}
				}
			}
			log.Println("Commitments verified")

			// resolve abilities and update chain state
			eventLogList := &[]comp.GameEventLog{}
			updateChainState := true

			printWhatActivated(turn.Msg.Result.Abilities)

			err = resolveAbilities(world, spell, spellPos, &turn.Msg.Result.Abilities, updateChainState, eventLogList) // pass eventLogList to record executed resolutions
			if err != nil {
				return msg.FulfillCastMsgResult{}, err
			}

			// Monster Turn occurs after abilities are resolved
			MonsterTurnSystem(world, eventLogList)

			// Emit activated abilities and spell log to client
			eventMap := make(map[string]any)
			eventMap["turnEvent"] = *eventLogList
			world.EmitEvent(eventMap)
			PrintStateToTerminal(world)

			// log to console
			// for _, logEntry := range *eventLogList {
			// 	log.Printf("X: %d, Y: %d, Event: %d\n",
			// 		logEntry.X, logEntry.Y, logEntry.Event)
			// }

			// return successfully
			// note: this msg returns to Seismic as the caller, not the player client
			return msg.FulfillCastMsgResult{
				LogEntry: *eventLogList,
			}, nil
		})
}

// Function to check if a string is in an array
func contains(array []string, value string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func printWhatActivated(Abilities [client.TotalAbilities]bool) {
	abilityNameMap := []string{"pureDamage, SideDamage, WallDamage, Explosion, UpHeal, RightHeal, DownHeal, LeftHeal"}
	for i := 0; i < len(abilityNameMap); i++ {
		if Abilities[i] {
			log.Println(abilityNameMap[i])
		}
	}
}
