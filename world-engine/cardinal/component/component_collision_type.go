package component

type CollideType int

const (
	PlayerCollide CollideType = iota
	MonsterCollide
	WallCollide
	ItemCollide
)

type Collidable struct {
	Type CollideType `json:"type"`
}

func (Collidable) Name() string {
	return "Collidable"
}

func (c Collidable) ToString() string {
	switch c.Type {
	case PlayerCollide:
		return "PlayerCollide"
	case MonsterCollide:
		return "MonsterCollide"
	case WallCollide:
		return "WallCollide"
	case ItemCollide:
		return "ItemCollide"
	default:
		return "idk wat"
	}
}
