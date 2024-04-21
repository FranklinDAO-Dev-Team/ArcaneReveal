package system

import (
	"cinco-paus/query"
	"fmt"
	"log"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

func PrintStateToTerminal(world cardinal.WorldContext, gameID types.EntityID) {
	log.Println("entered PrintStateToTerminal")
	log.Println("gameID: ", gameID)
	fmt.Println()
	gameState, err := query.GameState(world, &query.GameStateRequest{GameID: gameID})
	if err != nil {
		log.Println("query.GameState err: ", err)
		return
	}

	// TODO:
	boardSize := 11
	board := make([][]string, boardSize)
	for i := range board {
		board[i] = make([]string, boardSize)
		for j := range board[i] {
			board[i][j] = " " // Empty space
		}
	}

	// Place player
	board[gameState.Player.Y][gameState.Player.X] = "P"

	// Place walls
	for _, wall := range gameState.Walls {
		board[wall.Y][wall.X] = "X"
	}

	// Place monsters
	for _, monster := range gameState.Monsters {
		board[monster.Y][monster.X] = "M"
	}

	// Print the board
	for _, row := range board {
		for _, cell := range row {
			switch cell {
			case "P":
				fmt.Print("\033[34mP\033[0m ") // Player in blue
			case "X":
				fmt.Print("\033[31mX\033[0m ") // Wall in red
			case "M":
				fmt.Print("\033[35mM\033[0m ") // Monster in magenta
			default:
				fmt.Print(cell + " ")
			}
		}
		fmt.Println()
	}
}
