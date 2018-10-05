Insolar – Configuration
===============

[![GoDoc](https://godoc.org/github.com/insolar/insolar/configuration?status.svg)](https://godoc.org/github.com/insolar/insolar/configuration)


Package provides configuration params for all Insolar components and helper for config resources management.

### Configuration

Configuration struct is a root registry for all components config.
It provides constructor method `NewConfiguration()` which creates new instance of configuration object filled with default values.

Each root level Insolar component has a constructor with config as argument.
Each components should have its own config struct with the same name in this package.
Each config struct should have constructor which returns instance with default params.

### Holder

Package also provides [Holder](https://godoc.org/github.com/insolar/insolar/configuration#Holder) to easily manage config resources. 
It based on [Viper config solution for Golang](https://github.com/spf13/viper) and helps to Marshal\Unmarshal config structs, manage files, ENV and command line variables.

Holder provides functionality to merge configuration from different sources.

#### Merge priority

1. command line flags
2. ENV variables
3. yaml file
4. Default config

### Manage configuration from cli

Insolar cli tool helps user to manage configuration.

```
insolar config --help
```