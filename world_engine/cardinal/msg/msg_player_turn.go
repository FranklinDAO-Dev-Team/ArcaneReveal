package msg

import "errors"

type PlayerTurnMsg struct {
	Nickname  string `json:"nickname"`
	Action    string `json:"action"`
	Direction string `json:"direction"`
}

type PlayerTurnResult struct {
	Success bool `json:"success"`
}

func (m PlayerTurnMsg) ValFmt() error {
	validActions := map[string]bool{
		"attack": true,
		"wand":   true,
		"move":   true,
	}

	validDirections := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	if !validActions[m.Action] {
		return errors.New("invalid action")
	}

	if !validDirections[m.Direction] {
		return errors.New("invalid direction")
	}

	return nil
}
