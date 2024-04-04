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
	revealRequestCh chan RevealRequest
	revealReturnCh  chan RevealReqResponse
}

func New(
	proofRequestCh chan ProofRequest,
	proofReturnCh chan ProofReqResponse,
	revealRequestCh chan RevealRequest,
	revealReturnCh chan RevealReqResponse,
) *SeismicClient {
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

				gameState, hasGame := sc.store.GetGameState(req.PersonaTag)
				if !hasGame {
					sc.revealReturnCh <- RevealReqResponse{
						PersonaTag: req.PersonaTag,
						GameID:     req.GameID,
						CastID:     req.CastID,
						Success:    false,
						Error:      "no game found",
					}
					continue
				}

				castedAbilities := [TotalAbilities]bool{}
				salts := [TotalAbilities]string{}
				for i, canCast := range req.PotentialAbilities {
					wandHasAbility, salt := gameState.WandHasAbility(req.WandNum, i)
					castedAbilities[i] = canCast && wandHasAbility
					salts[i] = salt
				}

				sc.revealReturnCh <- RevealReqResponse{
					PersonaTag: req.PersonaTag,
					GameID:     req.GameID,
					CastID:     req.CastID,
					Success:    true,
					Abilities:  castedAbilities,
					Salts:      salts,
				}
			}
		}
	}()
}
