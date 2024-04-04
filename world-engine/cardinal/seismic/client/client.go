package client

import (
	"fmt"
	"math/big"
)

type SeismicClient struct {
	store           GameStateStore
	prover          *SeismicProver
	proofRequestCh  chan ProofRequest
	proofReturnCh   chan ProofReqResponse
	revealRequestCh chan int
	revealReturnCh  chan int
}

func New(proofRequestCh chan ProofRequest, proofReturnCh chan ProofReqResponse, revealRequestCh chan int, revealReturnCh chan int) *SeismicClient {
	store := NewGameStateStore()
	prover, err := NewProver()
	if err != nil {
		fmt.Printf("error creating prover: %v", err)
		panic(err)
	}
	return &SeismicClient{
		store:           store,
		prover:          prover,
		proofRequestCh:  proofRequestCh,
		proofReturnCh:   proofReturnCh,
		revealRequestCh: revealRequestCh,
		revealReturnCh:  revealReturnCh,
	}
}

func (sc *SeismicClient) Start() {
	go func() {
		for {
			select {
			case req := <-sc.proofRequestCh:
				playerSource, ok := new(big.Int).SetString(req.PlayerSource, 10)
				if !ok {
					sc.proofReturnCh <- NewProofFailResponse(req, "failed to parse playerSource")
					continue
				}

				gameState, err := NewGameState(playerSource)
				if err != nil {
					sc.proofReturnCh <- NewProofFailResponse(req, "failed to generate game state")
					continue
				}

				proof, err := sc.prover.Prove(gameState)
				if err != nil {
					sc.proofReturnCh <- NewProofFailResponse(req, "failed to prove")
					continue
				}

				sc.store.ReplaceGameState(req.PersonaTag, gameState)

				sc.proofReturnCh <- NewProofSuccessResponse(req, *proof)

			case req := <-sc.revealRequestCh:
				// TODO: legit response here, depending on game implementation
				sc.revealReturnCh <- 42
				fmt.Println("commit-reveal req:", req)
			}
		}
	}()
}
