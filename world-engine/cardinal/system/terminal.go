package system

import (
	"cinco-paus/query"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
)

func PrintStateToTerminal(world cardinal.WorldContext) {
	fmt.Println("hello world")
	gameState, _ := query.GameState(world, &query.GameStateRequest{})

	// TODO:

	fmt.Println("GameState: ", gameState)

}
