package system

import (
	"cinco-paus/component"
	"cinco-paus/msg"
	"cinco-paus/seismic/client"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/message"
)

func RequestGameSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage[msg.RequestGameMsg, msg.RequestGameMsgResult](
		world,
		func(req message.TxData[msg.RequestGameMsg]) (msg.RequestGameMsgResult, error) {
			personaTag := req.Tx.PersonaTag
			playerSource := req.Msg.PlayerSource

			id, err := cardinal.Create(world, component.PendingGame{
				PersonaTag:   personaTag,
				PlayerSource: playerSource,
			})
			if err != nil {
				return msg.RequestGameMsgResult{}, fmt.Errorf("failed to create pending game: %v", err)
			}

			proofReq := client.ProofRequest{
				PersonaTag:    personaTag,
				PendingGameID: id,
				PlayerSource:  playerSource,
			}
			proofRequestCh <- proofReq

			return msg.RequestGameMsgResult{}, nil
		})
}
