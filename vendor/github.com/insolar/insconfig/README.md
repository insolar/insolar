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

#Usage

see example and tests