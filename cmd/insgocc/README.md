# Insgocc
Toolset for developing builtin platform's and application's contracts.

## Glossary

### builtin contract
Single golang file which contains contract's code executed by virtual nodes.

### contract's proxy
Golang package, generated for contract. Proxy package is needed for a remote call.
When you need to make call from one contract to another, you use proxy.
 
### contract's wrapper
Single golang file, generated for contract. Wrapper is needed to unpack values from packed memory (cbor) into native structures, call a real function and pack the resulting structure back into cbor.

### contracts initialization file
This file describes the mappings that are required to search and interact with contracts and their methods in the built-in executor at virtual nodes.
In addition, code descriptors and prototype descriptors are initialized.

Use functions from initialization file for builtinContracts param at initialization of VirtualServer component.

## How to use

## Commands

### regen-builtin
Generates code for builtin contracts, makes them usable for application.
* generates proxies
* generates wrappers
* generates initialization file
#### Flags:
        -c, --contractsPath string   dir path to builtin contracts
        -h, --help                   help for regen-builtin
        -i, --importPath string      import path for builtin contracts packages

#### Example:
        ./bin/insgocc regen-builtin -c=application/builtin/contract -i=github.com/insolar/insolar/application/builtin/contract

### proxy
Generates contract's proxy.
#### Flags:
        -r, --code-reference         reference to code of contract
        -m, --machine-type           machine type (one of builtin/go) (default go)
        -o, --output file            output file (use - for STDOUT)

Contract code reference is insolar.Reference of prototype, saved on ledger by genesis component.

#### Example:
        ./bin/insgocc proxy -r=<contract code reference> -m=builtin -o=contact_proxy.go <path to single contract file>

### wrapper
Generates contract's wrapper.
#### Flags:
        -p, --panic-logical          panics are logical errors (turned off by default)
        -m, --machine-type           machine type (one of builtin/go) (default go)
        -o, --output file            output file (use - for STDOUT) (default -)

#### Example:
        ./bin/insgocc wrapper -p -m=builtin -o=contact_wrapper.go <path to single contract file>
