package query

import (
	// comp "cinco-paus/component"
	comp "cinco-paus/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

type PlayerData struct {
	X          int `json:"x"`
	Y          int `json:"y"`
	MaxHealth  int `json:"maxHealth"`
	CurrHealth int `json:"currHealth"`
}

type WandData struct {
	Number      int  `json:"number"`
	IsAvailable bool `json:"isAvailable"`
}

type WallData struct {
	X    int `json:"x"`
	Y    int `json:"y"`
	Type int `json:"type"`
}

type MonsterData struct {
	X    int `json:"x"`
	Y    int `json:"y"`
	Type int `json:"type"`
}

type GameStateRequest struct {
}

type GameStateResponse struct {
	Player   PlayerData    `json:"player"`
	Wands    []WandData    `json:"wands"`
	Walls    []WallData    `json:"walls"`
	Monsters []MonsterData `json:"monsters"`
}

func GameState(world cardinal.WorldContext, req *GameStateRequest) (*GameStateResponse, error) {
	fmt.Println("made it to query")
	PlayerData := PlayerData{}
	Wands := []WandData{}
	Walls := []WallData{}
	Monsters := []MonsterData{}

	var outsideErr error

	searchErr := cardinal.NewSearch(
		world,
		filter.Contains()).
		Each(func(id types.EntityID) bool {
			fmt.Printf("id: %v\n", id)
			player, _ := cardinal.GetComponent[comp.Player](world, id) // don't error check, want to ignore unfound errors
			if player != nil {
				pos, err := cardinal.GetComponent[comp.Position](world, id)
				if err != nil {
					fmt.Printf("outsideErr 2")
					outsideErr = fmt.Errorf("failed to get position component for player: %w", err)
				}
				health, err := cardinal.GetComponent[comp.Health](world, id)
				if err != nil {
					fmt.Printf("outsideErr 3")
					outsideErr = fmt.Errorf("failed to get position component for health: %w", err)
				}
				// found the player
				PlayerData.X = pos.X
				PlayerData.Y = pos.Y
				PlayerData.MaxHealth = health.MaxHealth
				PlayerData.CurrHealth = health.CurrHealth
			}

			wand, _ := cardinal.GetComponent[comp.WandCore](world, id)
			if wand != nil {
				availableObj, err := cardinal.GetComponent[comp.Available](world, id)
				if err != nil {
					fmt.Printf("outsideErr 5")
					outsideErr = fmt.Errorf("failed to get available component for wand: %w", err)
					return false
				} else {
					Wands = append(Wands, WandData{
						Number:      wand.Number,
						IsAvailable: availableObj.IsAvailable,
					})

				}
			}

			wall, _ := cardinal.GetComponent[comp.Wall](world, id)
			if wall != nil {
				pos, err := cardinal.GetComponent[comp.Position](world, id)
				if err != nil {
					fmt.Printf("outsideErr 7")
					outsideErr = fmt.Errorf("failed to get position component for wall: %w", err)
				}
				Walls = append(Walls, WallData{
					X:    pos.X,
					Y:    pos.Y,
					Type: int(wall.Type),
				})

			}

			monster, _ := cardinal.GetComponent[comp.Monster](world, id)
			if monster != nil {
				pos, err := cardinal.GetComponent[comp.Position](world, id)
				if err != nil {
					fmt.Printf("outsideErr 9")
					outsideErr = fmt.Errorf("failed to get position component for monster: %w", err)
				}
				Monsters = append(Monsters, MonsterData{
					X:    pos.X,
					Y:    pos.Y,
					Type: int(monster.Type),
				})

			}

			return true
		})

	if searchErr != nil {
		fmt.Printf("searchErr: %v\n", searchErr)
		return nil, searchErr
	}
	if outsideErr != nil {
		fmt.Printf("outsideErr: %v\n", outsideErr)
		return nil, outsideErr
	}

	fmt.Println("GameStateResponse")
	fmt.Println(PlayerData)
	fmt.Println(Wands)
	fmt.Println(Walls)
	fmt.Println(Monsters)
	return &GameStateResponse{
		Player:   PlayerData,
		Wands:    Wands,
		Walls:    Walls,
		Monsters: Monsters,
	}, nil

}
