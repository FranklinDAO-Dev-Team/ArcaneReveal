package query

import (
	// comp "cinco-paus/component"
	comp "cinco-paus/component"
	"fmt"
	"log"

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

func GameState(world cardinal.WorldContext, _ *GameStateRequest) (*GameStateResponse, error) {
	playerData := &PlayerData{}
	wands := &[]WandData{}
	walls := &[]WallData{}
	monsters := &[]MonsterData{}

	var outsideErr error

	searchErr := cardinal.NewSearch(
		world,
		filter.Contains()).
		Each(func(id types.EntityID) bool {
			// log.Printf("id: %v\n", id)

			outsideErr = getPlayerData(world, id, playerData)
			if outsideErr != nil {
				log.Printf("error getting player data: %v\n", outsideErr)
				return false
			}

			outsideErr = getWandData(world, id, wands)
			if outsideErr != nil {
				log.Printf("error getting wand data: %v\n", outsideErr)
				return false
			}

			outsideErr = getWallData(world, id, walls)
			if outsideErr != nil {
				log.Printf("error getting wall data: %v\n", outsideErr)
				return false
			}

			outsideErr = getMonsterData(world, id, monsters)
			if outsideErr != nil {
				log.Printf("error getting monster data: %v\n", outsideErr)
				return false
			}

			// always check next entity
			return true
		})

	if searchErr != nil {
		log.Printf("searchErr: %v\n", searchErr)
		return nil, searchErr
	}
	if outsideErr != nil {
		log.Printf("outsideErr: %v\n", outsideErr)
		return nil, outsideErr
	}

	return &GameStateResponse{
		Player:   *playerData,
		Wands:    *wands,
		Walls:    *walls,
		Monsters: *monsters,
	}, nil
}

func getPlayerData(world cardinal.WorldContext, id types.EntityID, playerData *PlayerData) error {
	player, _ := cardinal.GetComponent[comp.Player](world, id) // don't error check, want to ignore unfound errors
	if player != nil {
		pos, err := cardinal.GetComponent[comp.Position](world, id)
		if err != nil {
			return fmt.Errorf("failed to get position component for player: %w", err)
		}
		health, err := cardinal.GetComponent[comp.Health](world, id)
		if err != nil {
			return fmt.Errorf("failed to get position component for health: %w", err)
		}
		// found the player
		playerData.X = pos.X
		playerData.Y = pos.Y
		playerData.MaxHealth = health.MaxHealth
		playerData.CurrHealth = health.CurrHealth
	}
	return nil
}

func getWandData(world cardinal.WorldContext, id types.EntityID, wands *[]WandData) error {
	wand, _ := cardinal.GetComponent[comp.WandCore](world, id)
	if wand != nil {
		availableObj, err := cardinal.GetComponent[comp.Available](world, id)
		if err != nil {
			return fmt.Errorf("failed to get available component for wand: %w", err)
		}
		*wands = append(*wands, WandData{
			Number:      wand.Number,
			IsAvailable: availableObj.IsAvailable,
		})
	}
	return nil
}

func getWallData(world cardinal.WorldContext, id types.EntityID, walls *[]WallData) error {
	pos, err := cardinal.GetComponent[comp.Position](world, id)
	if err != nil {
		return fmt.Errorf("failed to get position component for wall: %w", err)
	}
	*walls = append(*walls, WallData{
		X:    pos.X,
		Y:    pos.Y,
		Type: int(comp.WallCollide),
	})
	return nil
}

func getMonsterData(world cardinal.WorldContext, id types.EntityID, monsters *[]MonsterData) error {
	pos, err := cardinal.GetComponent[comp.Position](world, id)
	if err != nil {
		return fmt.Errorf("failed to get position component for monster: %w", err)
	}
	*monsters = append(*monsters, MonsterData{
		X:    pos.X,
		Y:    pos.Y,
		Type: int(comp.MonsterCollide),
	})
	return nil
}
