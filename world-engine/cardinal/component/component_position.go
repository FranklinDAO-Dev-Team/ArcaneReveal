package component

type Position struct {
	X int `json:"x"`
	Y int `josn:"y"`
}

func (Position) Name() string {
	return "Position"
}
