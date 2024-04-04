package msg

import (
	"cinco-paus/seismic/client"

	"pkg.world.dev/world-engine/cardinal/types"
)

type FulfillCastMsg struct {
	Result    client.RevealReqResponse      `json:"res"`
	GameID    types.EntityID                `json:"gameID"`
	Success   bool                          `json:"success"`
	Abilities [client.TotalAbilities]bool   `json:"abilities"`
	Salts     [client.TotalAbilities]string `json:"salts"`
}

type FulfillCastMsgResult struct{}
