package component

type Player struct {
	Nickname      string `json:"nickname"`
	MaxHealth     int    `json:"maxhealth"`
	CurrentHealth int    `json:"currenthealth"`
	X             int    `json:"x"`
	Y             int    `josn:"y"`
}

func (Player) Name() string {
	return "Player"
}
