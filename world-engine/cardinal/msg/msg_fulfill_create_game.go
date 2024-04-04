package msg

import "cinco-paus/seismic/client"

type FulfillCreateGameMsg struct {
	Result client.ProofReqResponse `json:"res"`
}

type FulfillCreateGameMsgResult struct{}
