package component

type Game struct {
	PersonaTag  string
	Commitments *[][]string
	Reveals     *[][]string
	Level       int
}

func (Game) Name() string {
	return "Game"
}
