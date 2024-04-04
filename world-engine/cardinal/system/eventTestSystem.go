package system

import (
	comp "cinco-paus/component"
	"cinco-paus/msg"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/message"
)

func EventTestSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage[msg.EventTestMsg, msg.EventTestMsgResult](
		world,
		func(turn message.TxData[msg.EventTestMsg]) (msg.EventTestMsgResult, error) {

			fmt.Println("starting event test system")
			eventLogList := []comp.GameEventLog{
				comp.GameEventLog{
					X:     3,
					Y:     1,
					Event: comp.GameEventSpellBeam,
				},
				comp.GameEventLog{
					X:     5,
					Y:     1,
					Event: comp.GameEventSpellBeam,
				},
				comp.GameEventLog{
					X:     7,
					Y:     1,
					Event: comp.GameEventSpellBeam,
				},
				comp.GameEventLog{
					X:     9,
					Y:     1,
					Event: comp.GameEventSpellBeam,
				},
				comp.GameEventLog{
					X:     9,
					Y:     1,
					Event: comp.GameEventSpellDamage,
				},
				comp.GameEventLog{
					X:     9,
					Y:     3,
					Event: comp.GameEventSpellDamage,
				},
				comp.GameEventLog{
					X:     10,
					Y:     3,
					Event: comp.GameEventSpellDisappate,
				},
			}

			eventSliceList := eventLogListToSliceList(eventLogList)
			eventMap := make(map[string]any)
			eventMap["turnEvent"] = eventSliceList
			fmt.Println("eventMap", eventMap)
			world.EmitEvent(eventMap)

			return msg.EventTestMsgResult{}, nil
		})
}

func eventLogToSlice(eventLog comp.GameEventLog) []int {
	return []int{eventLog.X, eventLog.Y, int(eventLog.Event)}
}

func eventLogListToSliceList(eventLogList []comp.GameEventLog) [][]int {
	sliceList := make([][]int, len(eventLogList))
	for i, eventLog := range eventLogList {
		sliceList[i] = eventLogToSlice(eventLog)
	}
	return sliceList
}
