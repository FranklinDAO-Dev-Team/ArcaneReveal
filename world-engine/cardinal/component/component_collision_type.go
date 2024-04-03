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
