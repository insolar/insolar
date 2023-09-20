package extractor

import (
	"github.com/pkg/errors"

	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

// CallResponse extracts response of Call
func CallResponse(data []byte) (interface{}, *foundation.Error, error) {
	var result interface{}
	var contractErr *foundation.Error
	err := foundation.UnmarshalMethodResultSimplified(data, &result, &contractErr)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ CallResponse ] Can't unmarshal response ")
	}

	return result, contractErr, nil
}

// PublicKeyResponse extracts response of GetPublicKey
func PublicKeyResponse(data []byte) (string, error) {
	return stringResponse(data)
}
