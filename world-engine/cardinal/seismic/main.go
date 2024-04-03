package main

import (
	"cinco-paus/seismic/client"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"
)

func main() {
	seismicProver, _ := client.NewProver()

	playerSource := big.NewInt(int64(10))
	game, _ := client.NewGameState(playerSource)

	startTime := time.Now()
	proof, err := seismicProver.Prove(game)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	log.Printf("Prove function took %s", time.Since(startTime))
	proofJSON, _ := json.Marshal(proof)
	fmt.Println(string(proofJSON))

	res := client.Verify(*proof)
	fmt.Println("Verification Result:", res)
}
