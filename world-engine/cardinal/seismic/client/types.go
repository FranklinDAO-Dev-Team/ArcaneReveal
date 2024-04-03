package client

import "github.com/iden3/go-rapidsnark/types"

type ProofRequest struct {
	PersonaTag   string
	PlayerSource string
}

type ProofReqResponse struct {
	Success bool
	Proof   types.ZKProof
	Error   string
}

func NewProofSuccessResponse(proof types.ZKProof) ProofReqResponse {
	return ProofReqResponse{
		Success: true,
		Proof:   proof,
	}
}

func NewProofFailResponse(msg string) ProofReqResponse {
	return ProofReqResponse{
		Success: false,
		Error:   msg,
	}
}
