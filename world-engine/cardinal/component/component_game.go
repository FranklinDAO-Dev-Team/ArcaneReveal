package component

type Game struct {
	PersonaTag  string
	Commitments *[][]string
}

func (Game) Name() string {
	return "Game"
}
