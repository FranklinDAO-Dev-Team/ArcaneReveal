package component

type Game struct {
	PersonaTag  string
	Commitments *[][]string
	Reveals     *[][]int
	Level       int
}

func (Game) Name() string {
	return "Game"
}
