// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package extractor

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

// InfoResponse represents response from Info() method of RootDomain contract
type Info struct {
	RootDomain              string   `json:"rootDomain"`
	RootMember              string   `json:"rootMember"`
	MigrationDaemonsMembers []string `json:"migrationDaemonMembers"`
	MigrationAdminMember    string   `json:"migrationAdminMember"`
	NodeDomain              string   `json:"nodeDomain"`
}

// InfoResponse returns response from Info() method of RootDomain contract
func InfoResponse(data []byte) (*Info, error) {
	var infoMap interface{}
	var contractErr *foundation.Error
	err := foundation.UnmarshalMethodResultSimplified(data, &infoMap, &contractErr)
	if err != nil {
		return nil, errors.Wrap(err, "[ InfoResponse ] Can't unmarshal")
	}
	if contractErr != nil {
		return nil, errors.Wrap(contractErr, "[ InfoResponse ] Has error in response")
	}

	var info Info
	data = infoMap.([]byte)
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, errors.Wrap(err, "[ InfoResponse ] Can't unmarshal response ")
	}

	return &info, nil
}
