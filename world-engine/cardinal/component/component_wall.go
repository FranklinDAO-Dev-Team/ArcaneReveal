package component

type WallType int

const (
	WALL WallType = iota
	LOCK
	BOARDER
	ENTRY
	EXIT
)

type Wall struct {
	Type WallType
}

func (Wall) Name() string {
	return "Wall"
}
