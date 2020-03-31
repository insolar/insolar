// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package functest

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/application/sdk"
	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
)

const insolarRootMemberKeys = "root_member_keys.json"

var AppPath = []string{"insolar", "insolar"}

var info *sdk.InfoResponse
var Root AppUser

type AppUser struct {
	Ref              string
	PrivKey          string
	PubKey           string
	MigrationAddress string
}

func (m *AppUser) GetReference() string {
	return m.Ref
}

func (m *AppUser) GetPrivateKey() string {
	return m.PrivKey
}

func (m *AppUser) GetPublicKey() string {
	return m.PubKey
}

func loadMemberKeys(keysPath string, member *AppUser) error {
	text, err := ioutil.ReadFile(keysPath)
	if err != nil {
		return errors.Wrapf(err, "[ loadMemberKeys ] could't load member keys")
	}
	var data map[string]string
	err = json.Unmarshal(text, &data)
	if err != nil {
		return errors.Wrapf(err, "[ loadMemberKeys ] could't unmarshal member keys")
	}
	if data["private_key"] == "" || data["public_key"] == "" {
		return errors.New("[ loadMemberKeys ] could't find any keys")
	}
	member.PrivKey = data["private_key"]
	member.PubKey = data["public_key"]

	return nil
}

func LoadAllMembersKeys() error {
	path, err := launchnet.LaunchnetPath(AppPath, "configs", insolarRootMemberKeys)
	if err != nil {
		return err
	}
	err = loadMemberKeys(path, &Root)
	if err != nil {
		return err
	}
	return nil
}

func SetInfo() error {
	var err error
	info, err = sdk.Info(launchnet.TestRPCUrl)
	if err != nil {
		return errors.Wrap(err, "[ setInfo ] error sending request")
	}
	return nil
}

func AfterSetup() {
	Root.Ref = info.RootMember
}
