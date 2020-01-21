// Copyright 2020 Insolar Network Ltd.
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

package main

import (
	"encoding/json"
	"fmt"
)

type jsonCfg struct {
	pretty bool
	indent string
	prefix string
}

type jsonOpt func(*jsonCfg)

func setPretty(pretty bool) jsonOpt {
	return func(cfg *jsonCfg) { cfg.pretty = pretty }
}

func jsonPrefix(prefix string) jsonOpt {
	return func(cfg *jsonCfg) { cfg.prefix = prefix }
}

func printJSON(v interface{}, opts ...jsonOpt) {
	cfg := &jsonCfg{
		indent: "  ",
		prefix: "",
	}
	for _, o := range opts {
		o(cfg)
	}

	var b []byte
	var err error
	if cfg.pretty {
		b, err = json.MarshalIndent(v, cfg.prefix, cfg.indent)
	} else {
		b, err = json.Marshal(v)
	}
	if err != nil {
		panic(err)
	}

	if cfg.pretty {
		fmt.Printf(cfg.prefix)
	}
	fmt.Printf("%s\n", b)
}
