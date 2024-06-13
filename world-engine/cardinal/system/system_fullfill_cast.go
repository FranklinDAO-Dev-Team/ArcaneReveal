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
	"pkg.world.dev/world-engine/cardinal/types"
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
				return msg.FulfillCastMsgResult{}, fmt.Errorf("failed to marshal result to JSON: %w", err)
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
			log.Printf("FulfillCastSystem. Starting spellPos: (%d, %d)", spellPos.X, spellPos.Y)
			castGameObj, err := cardinal.GetComponent[comp.GameObj](world, turn.Msg.Result.CastID)
			if err != nil {
				return msg.FulfillCastMsgResult{}, err
			}
			gameID := castGameObj.GameID
			if turn.Msg.Result.GameID != gameID {
				err = fmt.Errorf("GameID %d in msg.Result does not match GameID %d in castGameObj", turn.Msg.Result.GameID, gameID)
				return msg.FulfillCastMsgResult{}, err
			}

			// remove cast entity
			err = cardinal.Remove(world, turn.Msg.Result.CastID)
			if err != nil {
				return msg.FulfillCastMsgResult{}, err
			}

			//  Check salts and that &turn.Msg.Abilities makes sense
			// also update the game reveals array
			err = verifySalts(world, gameID, turn, spell)
			if err != nil {
				return msg.FulfillCastMsgResult{}, err
			}
			log.Println("Commitments verified")

			// resolve abilities and update chain state
			// pass eventLogList to record executed resolutions
			eventLogList := &[]comp.GameEventLog{}
			updateChainState := true
			printWhatActivated(turn.Msg.Result.Abilities)
			err = resolveAbilities(world, gameID, spell, spellPos, &turn.Msg.Result.Abilities, updateChainState, eventLogList)
			if err != nil {
				return msg.FulfillCastMsgResult{}, err
			}

			// Monster Turn occurs after abilities are resolved
			err = MonsterTurnSystem(world, gameID, eventLogList)
			if err != nil {
				return msg.FulfillCastMsgResult{}, err
			}

			// Emit activated abilities and spell log to client
			err = world.EmitEvent(map[string]any{
				"event":  "turn-event",
				"log":    *eventLogList,
				"gameID": gameID,
			})
			if err != nil {
				return msg.FulfillCastMsgResult{}, err
			}
			PrintStateToTerminal(world, gameID)

			// return successfully
			// note: this msg returns to Seismic as the caller, not the player client
			return msg.FulfillCastMsgResult{
				LogEntry: *eventLogList,
			}, nil
		})
}

// verifies the salts agains the commitments in the game and updates the game reveals array
func verifySalts(
	world cardinal.WorldContext,
	gameID types.EntityID,
	turn message.TxData[msg.FulfillCastMsg],
	spell *comp.Spell,
) error {
	for i, canCastAbilityI := range turn.Msg.Result.Abilities {
		if canCastAbilityI {
			// calculate the commitment based on the revealed salt
			i64 := int64(i)
			salt := big.NewInt(0)
			base := 10
			salt.SetString(turn.Msg.Result.Salts[i], base)
			commitment, err := poseidon.Hash([]*big.Int{big.NewInt(i64), salt})
			if err != nil {
				return fmt.Errorf("failed to hash salt: %w", err)
			}

			// verify the calculated commitment matches the commitment in the game
			game, err := cardinal.GetComponent[comp.Game](world, gameID)
			if err != nil {
				return err
			}
			var wandCommits = (*game.Commitments)[spell.WandNumber]
			commitIndex := findIndex(wandCommits, commitment.String())
			if commitIndex == -1 {
				return fmt.Errorf("commitment %d does not match", i)
			}

			log.Printf("Ability %d (%s) activted \n", i, comp.AbilityMap[i+1].GetAbilityName())

			// update the reveals array
			(*game.Reveals)[spell.WandNumber][commitIndex] = i
			cardinal.SetComponent[comp.Game](world, gameID, game)
		}
	}
	return nil
}

func findIndex(arr []string, val string) int {
	for i, v := range arr {
		if v == val {
			return i
		}
	}
	return -1
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

func printWhatActivated(abilities [client.TotalAbilities]bool) {
	abilityNameMap := []string{"PureDamage, SideDamage, WallDamage, Explosion, UpHeal, RightHeal, DownHeal, LeftHeal"}
	for i := 0; i < len(abilityNameMap); i++ {
		if abilities[i] {
			log.Println(abilityNameMap[i])
		}
	}
}
