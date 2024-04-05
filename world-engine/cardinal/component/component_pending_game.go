package component

type PendingGame struct {
	PlayerSource string `json:"playerSource"`
	PersonaTag   string `json:"personaTag"`
}

func (PendingGame) Name() string {
	return "PendingGame"
}
