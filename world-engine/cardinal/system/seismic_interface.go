package system

import (
	"cinco-paus/msg"
	"cinco-paus/seismic/client"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/sign"
)

var (
	proofRequestCh  = make(chan client.ProofRequest)
	proofReturnCh   = make(chan client.ProofReqResponse)
	revealRequestCh = make(chan client.RevealRequest)
	revealReturnCh  = make(chan client.RevealReqResponse)
)

func Initialize(world *cardinal.World) *client.SeismicClient {
	fulFillCreateMsg, ok := world.GetMessageByFullName("game.fulfill-create-game")
	if !ok {
		fmt.Printf("error: no 'fulfill-create-game' message")
		return nil
	}

	fulFillCastMsg, ok := world.GetMessageByFullName("game.fulfill-cast")
	if !ok {
		fmt.Printf("error: no 'fulfill-cast' message")
		return nil
	}

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case proofRes := <-proofReturnCh:
				if proofRes.Success {
					ok := client.Verify(proofRes.Proof)
					fmt.Println("verification res:", ok)
				}

				payload := msg.FulfillCreateGameMsg{Result: proofRes}
				sig, err := sign.NewTransaction(privateKey, "Seismic.Systems", world.Namespace(), 0, payload)
				if err != nil {
					fmt.Printf("failed to sign new tx: %v", err)
				}

				world.AddTransaction(fulFillCreateMsg.ID(), payload, sig)

			case revealRes := <-revealReturnCh:
				payload := msg.FulfillCastMsg{Result: revealRes}
				sig, err := sign.NewTransaction(privateKey, "Seismic.Systems", world.Namespace(), 0, payload)
				if err != nil {
					fmt.Printf("failed to sign new tx: %v", err)
				}

				world.AddTransaction(fulFillCastMsg.ID(), payload, sig)
			}
		}
	}()

	return client.New(proofRequestCh, proofReturnCh, revealRequestCh, revealReturnCh)
}
