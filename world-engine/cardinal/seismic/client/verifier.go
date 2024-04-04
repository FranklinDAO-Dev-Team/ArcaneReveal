package client

import (
	"cinco-paus/seismic/circuit"

	"github.com/iden3/go-rapidsnark/types"
	"github.com/iden3/go-rapidsnark/verifier"
)

func Verify(proof types.ZKProof) bool {
	verifierError := verifier.VerifyGroth16(proof, circuit.VerificationKey)
	if verifierError != nil {
		return false
	}
	return true
}
