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

package genesis

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/insolar/insolar/application/contract/member"
	"github.com/insolar/insolar/application/contract/nodedomain"
	rootdomaincontract "github.com/insolar/insolar/application/contract/rootdomain"
	walletcontract "github.com/insolar/insolar/application/contract/wallet"
	"github.com/insolar/insolar/bootstrap"
	"github.com/insolar/insolar/insolar"
	"github.com/pkg/errors"
)

func generatePlugins(outDir string, insgoccBin string) error {
	args := []string{
		"compile-genesis-plugins",
		"-o", outDir,
	}

	fmt.Println(insgoccBin, strings.Join(args, " "))
	gocc := exec.Command(insgoccBin, args...)
	gocc.Stderr = os.Stderr
	gocc.Stdout = os.Stdout
	return gocc.Run()
}

type memout struct {
	name   string
	memory interface{}
}

func generateMemoryFiles(outDir string, rootPubKey string, rootBalance uint) error {
	var outs []memout

	outs = append(outs, memout{
		name: insolar.GenesisNameRootDomain,
		memory: &rootdomaincontract.RootDomain{
			RootMember:    bootstrap.ContractRootMember,
			NodeDomainRef: bootstrap.ContractNodeDomain,
		},
	})

	nd, _ := nodedomain.NewNodeDomain()
	outs = append(outs, memout{
		name:   insolar.GenesisNameNodeDomain,
		memory: nd,
	})

	m, err := member.New("RootMember", rootPubKey)
	if err != nil {
		return errors.Wrap(err, "root member constructor failed")
	}
	outs = append(outs, memout{
		name:   insolar.GenesisNameRootMember,
		memory: m,
	})

	w, err := walletcontract.New(rootBalance)
	if err != nil {
		return errors.Wrap(err, "failed to create wallet instance")
	}
	outs = append(outs, memout{
		name:   insolar.GenesisNameRootWallet,
		memory: w,
	})

	for _, o := range outs {
		memFile := filepath.Join(outDir, o.name+".bin")
		err := generateMemoryFile(memFile, o.memory)
		if err != nil {
			return errors.Wrapf(err, "failed to store domain memory for %v in file %v", o.name, memFile)
		}
	}

	return nil
}

func generateMemoryFile(memfile string, data interface{}) error {
	b, err := insolar.Serialize(data)
	if err != nil {
		return errors.Wrap(err, "[ activateNodeDomain ] node domain serialization")
	}

	return errors.Wrapf(
		ioutil.WriteFile(memfile, b, 0600),
		"can't write to file %v", memfile,
	)
}
