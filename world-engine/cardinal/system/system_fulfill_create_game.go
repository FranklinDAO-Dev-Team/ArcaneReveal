package system

import (
	"cinco-paus/component"
	"cinco-paus/msg"
	"cinco-paus/seismic/client"
	"fmt"
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/message"
)

func FulfillCreateGameSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage[msg.FulfillCreateGameMsg, msg.FulfillCreateGameMsgResult](
		world,
		func(req message.TxData[msg.FulfillCreateGameMsg]) (msg.FulfillCreateGameMsgResult, error) {

			pendingGameID := req.Msg.Result.PendingGameID
			pendingGame, hasPendingGame := cardinal.GetComponent[component.PendingGame](
				world,
				pendingGameID,
			)
			if hasPendingGame != nil {
				return msg.FulfillCreateGameMsgResult{}, fmt.Errorf("no pending game with id %v", pendingGameID)
			}

			// Delete pending game, regardless of proof success
			_ = cardinal.Remove(world, pendingGameID)

			if !req.Msg.Result.Success {
				return msg.FulfillCreateGameMsgResult{}, fmt.Errorf("proving failed")
			}

			pubSignals := req.Msg.Result.Proof.PubSignals
			circuitPlayerSource := pubSignals[0]

			if circuitPlayerSource != pendingGame.PlayerSource {
				return msg.FulfillCreateGameMsgResult{}, fmt.Errorf("playerSource public signal assigned incorrectly")
			}

			proofOk := client.Verify(req.Msg.Result.Proof)
			if !proofOk {
				return msg.FulfillCreateGameMsgResult{}, fmt.Errorf("zkp verification failed")
			}

			commitments := make([][]string, client.NumWands)
			for i := range commitments {
				commitments[i] = make([]string, client.NumAbilities)
				for j := range commitments[i] {
					commitments[i][j] = pubSignals[1+i*client.NumAbilities+j]
				}
			}

			_, err := cardinal.Create(world, component.Game{
				PersonaTag:  req.Msg.Result.PersonaTag,
				Commitments: &commitments,
			})
			if err != nil {
				return msg.FulfillCreateGameMsgResult{}, fmt.Errorf("failed to create Game component: %v", err)
			}

			return msg.FulfillCreateGameMsgResult{}, nil
		})
}
