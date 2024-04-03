package msg

import (
	"errors"
	"fmt"
	"strconv"

	comp "cinco-paus/component"
)

type PlayerTurnMsg struct {
	Action    string `json:"action"`
	Direction string `json:"direction"`
	WandNum   string `json:"wandnum"`
}

type PlayerTurnResult struct {
	Success bool
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

	wandnum, err := strconv.Atoi(m.WandNum)
	if err != nil {
		return fmt.Errorf("Error converting string to int: %w", err)
	}
	if m.Action == "wand" && (wandnum < 0 || wandnum >= comp.NUM_WANDS) {
		return fmt.Errorf("invalid wand number: %d", wandnum)
	}

	return nil
}
