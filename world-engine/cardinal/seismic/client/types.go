package client

import (
	zkTypes "github.com/iden3/go-rapidsnark/types"
	worldEngineTypes "pkg.world.dev/world-engine/cardinal/types"
)

type ProofRequest struct {
	PersonaTag    string
	PendingGameID worldEngineTypes.EntityID
	PlayerSource  string
}

type ProofReqResponse struct {
	PersonaTag    string
	PendingGameID worldEngineTypes.EntityID
	Success       bool
	Proof         zkTypes.ZKProof
	Error         string
}

func NewProofSuccessResponse(req ProofRequest, proof zkTypes.ZKProof) ProofReqResponse {
	return ProofReqResponse{
		PersonaTag:    req.PersonaTag,
		PendingGameID: req.PendingGameID,
		Success:       true,
		Proof:         proof,
	}
}

func NewProofFailResponse(req ProofRequest, msg string) ProofReqResponse {
	return ProofReqResponse{
		PersonaTag:    req.PersonaTag,
		PendingGameID: req.PendingGameID,
		Success:       false,
		Error:         msg,
	}
}

type RevealRequest struct {
	PersonaTag string
	GameID     worldEngineTypes.EntityID
}

type RevealReqResponse struct {
	PersonaTag string
	GameID     worldEngineTypes.EntityID
}
