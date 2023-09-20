package sdk

import (
	"github.com/insolar/insolar/api/requester"
)

// InfoResponse represents response from rpc on network.getInfo method
type InfoResponse struct {
	RootDomain string `json:"rootDomain"`
	RootMember string `json:"rootMember"`
	NodeDomain string `json:"nodeDomain"`
	TraceID    string `json:"traceID"`
}

type rpcInfoResponse struct {
	requester.Response
	Result InfoResponse `json:"result"`
}
