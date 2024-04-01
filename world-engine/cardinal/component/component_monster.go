package component

type Monster struct {
	Type       string   `json:"type"`
	StatusList []string `json:"statusList"`
}

func (Monster) Name() string {
	return "Monster"
}
