package client

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/iden3/go-rapidsnark/prover"
	"github.com/iden3/go-rapidsnark/types"
	"github.com/iden3/go-rapidsnark/witness"

	"cinco-paus/seismic/circuit"
)

type SeismicProver struct {
	witnessCalculator *witness.Circom2WitnessCalculator
}

func NewProver() (*SeismicProver, error) {
	calc, err := witness.NewCircom2WitnessCalculator(circuit.CircuitWasm, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create witness calculator: %v", err)
	}
	return &SeismicProver{
		witnessCalculator: calc,
	}, nil
}

func (seismicProver *SeismicProver) Prove(gameState *GameState) (*types.ZKProof, error) {
	inputs := gameState.ToCircuitInputs()

	jsonInputs, err := json.Marshal(inputs)
	if err != nil {
		return nil, fmt.Errorf("failed to JSONify inputs: %v", err)
	}
	parsedInputs, err := witness.ParseInputs(jsonInputs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse inputs: %v", err)
	}

	wtns, err := seismicProver.witnessCalculator.CalculateWTNSBin(parsedInputs, true)
	if err != nil {
		log.Println("failed to compute witness:", err)
		log.Println("Likely that size of expected circuit input and actual is misaligned")
		log.Println("Double check that the constants in constants.go and init.circom are the same")
		return nil, fmt.Errorf("failed to computed witness: %v", err)
	}

	proof, err := prover.Groth16Prover(circuit.CircuitZkey, wtns)
	if err != nil {
		return nil, fmt.Errorf("failed to generate proof: %v", err)
	}

	return proof, nil
}
