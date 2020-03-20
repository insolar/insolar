# insconfig
Config management library.
This is the wrapper on https://github.com/spf13/viper library

Key features:
- .yaml format
- No default config path, path is explicitly set by --config/-c flag. Optionally you can override this by implementing ConfigPathGetter (look at tests)
- Environment overrides file values
- Can use only ENV, without file at all
- Optionally prints config to log on start
- No default values, all values are set explicitly, if not - returns error
- No unnecessary values (both in file and ENV), if not - returns error
- Supports custom flags, go flags and pflags
- Doesn't support overriding config by flags
- [wip] Generates empty yaml file with descriptions
- [wip] By default adds 2 flags --config и --gen-config
- Doesn't support overriding config on runtime
- Supports custom viper decode hooks

# Running example 
```
go run ./example/example.go --config="./example/example_config.yaml"
```

# Usage

With custom go flags (from example.go)
```go
    var testflag1 = flag.String("testflag1", "", "testflag1")
	mconf := Config{}
	params := insconfig.Params{
		EnvPrefix:    "example",
		ConfigPathGetter: &insconfig.FlagPathGetter{
			GoFlags: flag.CommandLine,
		},
	}
    insConfigurator := insconfig.New(params)
    _ = insConfigurator.Load(&mconf)
    fmt.Println(testflag1)
```

With custom spf13/pflags
```go
    var testflag1 = pflag.String("testflag1", "", "testflag1")
    mconf := Config{}
    params := insconfig.Params{
        EnvPrefix:    "example",
        ConfigPathGetter: &insconfig.PFlagPathGetter{
            PFlags: pflag.CommandLine,
        },
    }
    insConfigurator := insconfig.New(params)
    _ = insConfigurator.Load(&mconf)
    fmt.Println(testflag1)
```

With spf13/cobra. Cobra doesn't provide tools to manage flags parsing, so you need to add config flag yourself

```go
func main () {
    var configPath string
    rootCmd := &cobra.Command{
        Use: "insolard",
    }
    rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to config file")
    _ = rootCmd.MarkPersistentFlagRequired("config")
    err := rootCmd.Execute()

    // ...

    // To set your path from flag to insconfig you need to implement simple ConfigPathGetter interface and return path 
    type stringPathGetter struct {
        Path string
    }
    
    func (g *stringPathGetter) GetConfigPath() string {
        return g.Path
    }
}

func read(){
    cfg := ConfigStruct{}
    params := insconfig.Params{
        EnvPrefix:        "InsolarEnvPrefix",
        ConfigPathGetter: &stringPathGetter{Path: configPath},
        FileRequired:     false,
    }
    insConfigurator := insconfig.NewInsConfigurator(h.Params)
    err := insConfigurator.Load(&cfg)
    println(insconfig.ToString(cfg))
}
```