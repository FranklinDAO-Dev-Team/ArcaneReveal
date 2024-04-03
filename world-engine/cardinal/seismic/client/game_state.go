package client

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/iden3/go-iden3-crypto/poseidon"
)

type GameState struct {
	PlayerSource  string     `json:"playerSource"`
	SeismicSource string     `json:"seismicSource"`
	Commitments   [][]string `json:"commitments"`
	//priv
	Attributes [][]string `json:"attributes"`
	Salts      [][]string `json:"salts"`
}

func NewGameState(playerSource *big.Int) (*GameState, error) {
	maxBigInt := big.NewInt(0)
	maxBigInt.SetString(SnarkFieldSize, 10)

	seismicSource, err := GetRandomSource()
	if err != nil {
		return nil, fmt.Errorf("failed to generate random seed: %v", err)
	}

	attributes := make([][]string, NStaffs)
	salts := make([][]string, NStaffs)
	commitments := make([][]string, NStaffs)

	matchSource, err := poseidon.Hash([]*big.Int{playerSource, seismicSource})
	if err != nil {
		return nil, fmt.Errorf("failed to hash sourced randomness: %v", err)
	}

	allSpells := Permutate(matchSource, NAllSpells)

	for i := range attributes {
		attributes[i] = make([]string, NStaffSpells)
		salts[i] = make([]string, NStaffSpells)
		commitments[i] = make([]string, NStaffSpells)
		for j := range attributes[i] {
			index := i*NStaffSpells + j
			attribute := big.NewInt(int64(allSpells[index]))
			salt, err := poseidon.Hash([]*big.Int{big.NewInt(int64(index)), seismicSource})
			if err != nil {
				return nil, fmt.Errorf("failed to generate salt: %v", err)
			}

			commitment, err := poseidon.Hash([]*big.Int{attribute, salt})
			if err != nil {
				return nil, fmt.Errorf("failed to compute poseidon hash: %v", err)
			}

			attributes[i][j] = attribute.String()
			salts[i][j] = salt.String()
			commitments[i][j] = commitment.String()
		}
	}

	return &GameState{
		PlayerSource:  playerSource.String(),
		SeismicSource: seismicSource.String(),
		Attributes:    attributes,
		Salts:         salts,
		Commitments:   commitments,
	}, nil
}

func GetRandomSource() (*big.Int, error) {
	maxBigInt := big.NewInt(0)
	maxBigInt.SetString(SnarkFieldSize, 10)
	r, err := rand.Int(rand.Reader, maxBigInt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %v", err)
	}
	return r, nil
}

func (game *GameState) ToCircuitInputs() map[string]interface{} {
	input := map[string]interface{}{
		"playerSource":  game.PlayerSource,
		"seismicSource": game.SeismicSource,
		"attributes":    game.Attributes,
		"salts":         game.Salts,
		"commitments":   game.Commitments,
	}

	return input
}
