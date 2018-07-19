Insolar
===============
Blockchain platform

[![Build Status](https://travis-ci.org/insolar/network.svg?branch=master)](https://travis-ci.org/insolar/network)
[![Go Report Card](https://goreportcard.com/badge/github.com/insolar/network)](https://goreportcard.com/report/github.com/insolar/network)
[![GoDoc](https://godoc.org/github.com/insolar/network?status.svg)](https://godoc.org/github.com/insolar/network)

_This project is still in early development state.
It is not recommended to use it in production environment._

Overview
--------
**Insolar** is a blockchain platform developed by INS

We took [Kademlia DHT](https://en.wikipedia.org/wiki/Kademlia) original specifications and made significant improvements to make it ready
for real world application by enterprises.

#### Key features of our blockchain network layer:
 - **Support of heterogeneous network topology** with different types of nodes being able to communicate with each other.
   In classic peer-to-peer networks, any node can communicate directly with any other node on the network.
   In a real enterprise environment, this condition is often unacceptable for a variety of reasons including security.
 - **Network routing with a node or node group becoming relays** for others nodes.
   The network can continue to function despite various network restrictions such as firewalls, NATs, etc.
 - **Ability to limit number of gateways to corporate node group via relays** to keep the node group secure while being
   able to interact with the rest of the network through relays. This feature mitigates the risk of DDoS attacks.


Installation
------------

    go get github.com/insolar/network


Contributing
------------

Please feel free to submit issues, fork the repository and send pull requests!

When submitting an issue, we ask that you please include a complete test function that demonstrates the issue.
