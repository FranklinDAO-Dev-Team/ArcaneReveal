package msg

import (
	"cinco-paus/component"
	"cinco-paus/seismic/client"
)

type FulfillCastMsg struct {
	Result client.RevealReqResponse `json:"res"`
}

type FulfillCastMsgResult struct {
	LogEntry []component.GameEventLog
}
