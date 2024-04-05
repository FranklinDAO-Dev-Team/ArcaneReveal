package msg

import (
	"cinco-paus/seismic/client"
	"errors"
	"fmt"
	"strconv"
)

type PlayerTurnMsg struct {
	GameIDStr string `json:"gameIDStr"`
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
		return fmt.Errorf("error converting string to int: %w", err)
	}
	if m.Action == "wand" && (wandnum < 0 || wandnum >= client.NumWands) {
		return fmt.Errorf("invalid wand number: %d", wandnum)
	}

	return nil
}
