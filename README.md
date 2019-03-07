# Insolar

Enterprise-ready blockchain platform

[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/2150/badge)](https://bestpractices.coreinfrastructure.org/projects/2150)

[![Build Status](https://travis-ci.org/insolar/insolar.svg?branch=master)](https://travis-ci.org/insolar/insolar)
[![GolangCI](https://golangci.com/badges/github.com/insolar/insolar.svg)](https://golangci.com/r/github.com/insolar/insolar/)
[![Go Report Card](https://goreportcard.com/badge/github.com/insolar/insolar)](https://goreportcard.com/report/github.com/insolar/insolar)
[![GoDoc](https://godoc.org/github.com/insolar/insolar?status.svg)](https://godoc.org/github.com/insolar/insolar)
[![codecov](https://codecov.io/gh/insolar/insolar/branch/master/graph/badge.svg)](https://codecov.io/gh/insolar/insolar)

## Overview

**Insolar** is building a 4th generation blockchain platform for business aimed to enable seamless interactions between companies and unlock new growth opportunities. In addition to the blockchain platform, Insolar will provide blockchain services and ecosystem support for companies that are looking to develop and deploy blockchain solutions. Insolar will feature most complete and secure set of production-ready business blockchain tools and services to quickly build or launch blockchain enterprise applications, accelerating the progression path from initial proof-of-concept to full-scale production.

The world’s most innovative companies in finance, logistics, consumer goods, energy, healthcare, transportation, manufacturing and others will be turning to Insolar to create applications and networks that deliver tangible business success. They recognise that even in today’s digital economy, vast amounts of value continue to be trapped inside processes and organisations that don’t connect. Insolar is their remedy, helping them discover and design business value in blockchain networks — starting, accelerating and innovating strategies that replace longstanding business friction with trust and transparency. Delegating trust to a blockchain means that businesses can pursue broader networks, onboard new partners, and enter new ecosystems with ease. Blockchain-based networks that support multiparty collaboration around shared, trusted data and process automation across organisational boundaries bring benefits at many levels, starting with efficiency gains and culminating in reinventing how entire industry ecosystems operate.

Insolar is a global team of 60+ people in North America and Europe, including a 35-strong engineering team with practical blockchain engineering know-how, and 10 leading blockchain academics from major institutions (York University, ETH Zurich, Princeton).

## Components

### [Network](network)

Blockchain network layer.

* Support of heterogeneous network topology.
* Network routing with a host or host group becoming relays for others hosts.
* Ability to limit number of gateways to corporate host group via relays to keep the host group secure.

See [package readme](network/dhtnetwork) for more details.

### [Ledger](ledger)

Record storage engine backed by [BadgerDB](https://github.com/dgraph-io/badger).

### [Virtual machines](vm)

Various engines for smart contract execution:

* [wasm](vm/wasm) - WebAssembly implementation of smart contracts

### [Application layer](application)

Application module describes interaction of system components with each other.
Every component of the system is a `SmartContract`. Members of the system are given the opportunity to build their own dApps by publishing smart contracts in `Domain` instances.
Domains define the visibility scope for the child contracts and their interaction policies. Actually, `Domain` is subclass of `SmartContract`.

See [package readme](application) for more details.

### [Configuration](configuration)

Provides configuration params for all Insolar components and helper for config resources management.

### [Metrics](metrics)

Using Prometheus monitoring system and time series database for collecting and store metrics

## Installation

Download Insolar package

    go get github.com/insolar/insolar

Go to package directory

    cd $GOPATH/src/github.com/insolar/insolar

Install dependencies and build binaries

    make install-deps pre-build build

### Example

Run launcher:

    scripts/insolard/launchnet.sh -g

It will generate genesis data and launch a number of nodes. Default number is 5, you can uncomment more nodes in `scripts/insolard/genesis.yaml`.

After node processes are started you will see messages like “NODE 3 STARTED in background” in log and PulseWatcher will be started.
When you see `Ready` in Insolar State you can run test scripts and benchmarks:

    bin/apirequester -k=scripts/insolard/configs/root_member_keys.json -u=http://127.0.0.1:19101/api

This tool runs such scenario: it creates a number of users with wallets, then transfers some money between these users. First time script does it sequentially, second time — concurrently.
Options:
* `-k`: Path to root user keypair. All requests to create new user must be signed by root user.
* `-u`: Node API URL. By default first node listens on 127.0.0.1:19101. It can be changed in config.

Run benchmark

    bin/benchmark -c 2 -r 4 -k=scripts/insolard/configs/root_member_keys.json

Options:
* `-k`: Same as above, path to root user keypair.
* `-c`: Number of concurrent threads in which requests will be sent.
* `-r`: Number of transfer requests that will be sent in each thread.

After testing you can stop all nodes by pressing Ctrl+C.

#### See [apirequester](cmd/apirequester) and [benchmark](cmd/benchmark) readme for more details

## Contributing

See [Contributing Guidelines](.github/CONTRIBUTING.md).

## License

This project is licensed under the terms of the [Apache license 2.0](LICENSE), except for the [Network](network), [NetworkCoordinator](networkcoordinator) and [Consensus](consensus)  subdirectories, which are licensed under the terms of the [BSD 3-Clause Clear License](network/LICENSE.md).
