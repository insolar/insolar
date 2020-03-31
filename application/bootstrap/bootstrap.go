// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package bootstrap

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/applicationbase/bootstrap"

	"github.com/insolar/insolar/insolar/secrets"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// CreateGenesisContractsConfig generates bootstrap data.
//
// 1. read application-related keys files.
// 2. generates genesis contracts config for heavy node.
func CreateGenesisContractsConfig(ctx context.Context, configFile string) (map[string]interface{}, error) {
	config, err := ParseContractsConfig(configFile)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[ bootstrap ] config:\n%v\n", bootstrap.DumpAsJSON(config))

	inslog := inslogger.FromContext(ctx)

	inslog.Info("[ bootstrap ] read keys files")
	rootPublicKey, err := secrets.GetPublicKeyFromFile(filepath.Join(config.MembersKeysDir, "root_member_keys.json"))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get root keys")
	}

	return map[string]interface{}{
		"RootBalance":   config.RootBalance,
		"RootPublicKey": rootPublicKey,
	}, nil
}
