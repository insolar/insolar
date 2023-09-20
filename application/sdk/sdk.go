package sdk

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/api/requester"
)

// Info makes rpc request to network.getInfo method and extracts it
func Info(url string) (*InfoResponse, error) {
	body, err := requester.GetResponseBodyPlatform(url, "network.getInfo", nil)
	if err != nil {
		return nil, errors.Wrap(err, "[ Info ]")
	}

	infoResp := rpcInfoResponse{}

	err = json.Unmarshal(body, &infoResp)
	if err != nil {
		return nil, errors.Wrap(err, "[ Info ] Can't unmarshal")
	}
	if infoResp.Error != nil {
		return nil, infoResp.Error
	}

	return &infoResp.Result, nil
}
