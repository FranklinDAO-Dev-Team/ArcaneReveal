package msg

import "cinco-paus/seismic/client"

type FulfillCastMsg struct {
	Result client.RevealReqResponse `json:"res"`
}

type FulfillCastMsgResult struct{}
