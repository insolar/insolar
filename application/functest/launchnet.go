// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package functest

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	yaml "gopkg.in/yaml.v2"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/application/sdk"
	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
)

const insolarRootMemberKeys = "root_member_keys.json"
const insolarMigrationAdminMemberKeys = "migration_admin_member_keys.json"
const insolarFeeMemberKeys = "fee_member_keys.json"

var ApplicationIncentives [application.GenesisAmountApplicationIncentivesMembers]*AppUser
var NetworkIncentives [application.GenesisAmountNetworkIncentivesMembers]*AppUser
var Enterprise [application.GenesisAmountEnterpriseMembers]*AppUser
var Foundation [application.GenesisAmountFoundationMembers]*AppUser

var AppPath = []string{"insolar", "insolar"}

var info *sdk.InfoResponse
var Root AppUser
var MigrationAdmin AppUser
var FeeMember AppUser
var MigrationDaemons [application.GenesisAmountMigrationDaemonMembers]*AppUser

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

func GetNumShards() (int, error) {
	type bootstrapConf struct {
		PKShardCount int `yaml:"ma_shard_count"`
	}

	var conf bootstrapConf

	path, err := launchnet.LaunchnetPath(AppPath, "bootstrap.yaml")
	if err != nil {
		return 0, err
	}
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, errors.Wrap(err, "[ GetNumShards ] Can't read bootstrap config")
	}

	err = yaml.Unmarshal(buff, &conf)
	if err != nil {
		return 0, errors.Wrap(err, "[ GetNumShards ] Can't parse bootstrap config")
	}

	return conf.PKShardCount, nil
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
	path, err = launchnet.LaunchnetPath(AppPath, "configs", insolarFeeMemberKeys)
	if err != nil {
		return err
	}
	err = loadMemberKeys(path, &FeeMember)
	if err != nil {
		return err
	}
	path, err = launchnet.LaunchnetPath(AppPath, "configs", insolarMigrationAdminMemberKeys)
	if err != nil {
		return err
	}
	err = loadMemberKeys(path, &MigrationAdmin)
	if err != nil {
		return err
	}
	for i := range MigrationDaemons {
		path, err := launchnet.LaunchnetPath(AppPath, "configs", "migration_daemon_"+strconv.Itoa(i)+"_member_keys.json")
		if err != nil {
			return err
		}
		var md AppUser
		err = loadMemberKeys(path, &md)
		if err != nil {
			return err
		}
		MigrationDaemons[i] = &md
	}

	for i := 0; i < application.GenesisAmountApplicationIncentivesMembers; i++ {
		path, err := launchnet.LaunchnetPath(AppPath, "configs", "application_incentives_"+strconv.Itoa(i)+"_member_keys.json")
		if err != nil {
			return err
		}
		var md AppUser
		err = loadMemberKeys(path, &md)
		if err != nil {
			return err
		}
		ApplicationIncentives[i] = &md
	}

	for i := 0; i < application.GenesisAmountNetworkIncentivesMembers; i++ {
		path, err := launchnet.LaunchnetPath(AppPath, "configs", "network_incentives_"+strconv.Itoa(i)+"_member_keys.json")
		if err != nil {
			return err
		}
		var md AppUser
		err = loadMemberKeys(path, &md)
		if err != nil {
			return err
		}
		NetworkIncentives[i] = &md
	}

	for i := 0; i < application.GenesisAmountFoundationMembers; i++ {
		path, err := launchnet.LaunchnetPath(AppPath, "configs", "foundation_"+strconv.Itoa(i)+"_member_keys.json")
		if err != nil {
			return err
		}
		var md AppUser
		err = loadMemberKeys(path, &md)
		if err != nil {
			return err
		}
		Foundation[i] = &md
	}

	for i := 0; i < application.GenesisAmountEnterpriseMembers; i++ {
		path, err := launchnet.LaunchnetPath(AppPath, "configs", "enterprise_"+strconv.Itoa(i)+"_member_keys.json")
		if err != nil {
			return err
		}
		var md AppUser
		err = loadMemberKeys(path, &md)
		if err != nil {
			return err
		}
		Enterprise[i] = &md
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
	MigrationAdmin.Ref = info.MigrationAdminMember
}
