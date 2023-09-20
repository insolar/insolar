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
