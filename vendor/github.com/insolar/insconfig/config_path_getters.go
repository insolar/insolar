package insconfig

import (
	goflag "flag"

	flag "github.com/spf13/pflag"
)

// DefaultPathGetter adds "--config" flag and read path from it
type DefaultPathGetter struct {
	GoFlags *goflag.FlagSet
}

func (g *DefaultPathGetter) GetConfigPath() string {
	configPath := flag.String("config", "", "path to config")
	flag.Parse()
	return *configPath
}

// FlagPathGetter made for go flags compatibility
// Adds "--config" flag and read path from it, custom go flags should be created before and set to GoFlags
type FlagPathGetter struct {
	GoFlags *goflag.FlagSet
}

func (g *FlagPathGetter) GetConfigPath() string {
	if g.GoFlags != nil {
		flag.CommandLine.AddGoFlagSet(g.GoFlags)
	}
	configPath := flag.String("config", "", "path to config")
	flag.Parse()
	return *configPath
}

// PFlagPathGetter made for spf13/pflags compatibility.
// Adds "--config" flag and read path from it, custom pflags should be created before and set to PFlags
type PFlagPathGetter struct {
	PFlags *flag.FlagSet
}

func (g *PFlagPathGetter) GetConfigPath() string {
	if g.PFlags != nil {
		flag.CommandLine.AddFlagSet(g.PFlags)
	}
	configPath := flag.String("config", "", "path to config")
	flag.Parse()
	return *configPath
}
