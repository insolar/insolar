//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package extractor

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"

	"github.com/pkg/errors"
)

func Generic(data []byte, returns ...interface{}) error {
	res := foundation.Result{
		Returns: returns,
	}
	err := insolar.Deserialize(data, &res)
	if err != nil {
		return errors.Wrap(err, "[ StringResponse ] Can't unmarshal response")
	}

	// this magic helper that injects logic error into one of returns, just sugar
	if res.Error != nil {
		found := false
		for i := 0; i < len(returns); i++ {
			if e, ok := returns[i].(**foundation.Error); ok {
				if e == nil {
					return errors.New("nil pointer in unmarshal")
				}
				*e = res.Error
				found = true
			}
		}
		if !found {
			return errors.New("no place for error in returns")
		}
	}

	return nil
}

func stringResponse(data []byte) (string, error) {
	var result string
	var contractErr *foundation.Error
	err := Generic(data, &result, &contractErr)
	if err != nil {
		return "", errors.Wrap(err, "[ StringResponse ] Can't unmarshal response ")
	}
	if contractErr != nil {
		return "", errors.Wrap(contractErr, "[ StringResponse ] Has error in response")
	}

	return result, nil
}
