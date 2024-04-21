package component

import "pkg.world.dev/world-engine/cardinal/types"

type GameObj struct {
	GameID types.EntityID
}

func (GameObj) Name() string {
	return "GameObj"
}
