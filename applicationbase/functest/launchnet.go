///
// Copyright 2020 Insolar Technologies GmbH
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
///

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
var Root CommonUser

type CommonUser struct {
	Ref     string
	PrivKey string
	PubKey  string
}

func (m *CommonUser) GetReference() string {
	return m.Ref
}

func (m *CommonUser) GetPrivateKey() string {
	return m.PrivKey
}

func (m *CommonUser) GetPublicKey() string {
	return m.PubKey
}

func loadMemberKeys(keysPath string, member *CommonUser) error {
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
