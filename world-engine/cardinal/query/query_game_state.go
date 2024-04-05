package query

import "pkg.world.dev/world-engine/cardinal"

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
	return nil, nil
}
