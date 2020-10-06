[<img src="https://github.com/insolar/doc-pics/raw/master/st/github-readme-banner.png">](http://insolar.io/?utm_source=Github)

Insolar platform is the most secure, scalable, and comprehensive business-ready blockchain toolkit in the world. Insolar’s goal is to give businesses access to features and services that enable them to launch new decentralized applications quickly and easily. Whether a minimum viable product or full-scale production software, Insolar builds and integrates applications for your enterprise's existing systems.

[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/2150/badge)](https://bestpractices.coreinfrastructure.org/projects/2150)
[![GolangCI](https://golangci.com/badges/github.com/insolar/insolar.svg)](https://golangci.com/r/github.com/insolar/insolar/)
[![Go Report Card](https://goreportcard.com/badge/github.com/insolar/insolar)](https://goreportcard.com/report/github.com/insolar/insolar)
[![GoDoc](https://godoc.org/github.com/insolar/insolar?status.svg)](https://godoc.org/github.com/insolar/insolar)
[![codecov](https://codecov.io/gh/insolar/insolar/branch/master/graph/badge.svg)](https://codecov.io/gh/insolar/insolar)

# Insolar Platform 1.0

Insolar Platform 1.0 resides in this repository. Insolar MainNet application (and its benchmark) resides in the [insolar/mainnet repository](https://github.com/insolar/mainnet) and runs on top of Platform 1.0.

Assured Ledger, Insolar's Platform-as-a-Service solution, is being actively developed in the [insolar/assured-ledger repository](https://github.com/insolar/assured-ledger).

## Quick start

To learn what distinguishes Insolar from other blockchain projects, go through the [list of our features](https://insolar.io/platform?utm_source=Github). 

To get a grip on how Insolar works, take a look at the [big picture](https://docs.insolar.io/en/latest/basics.html) and explore the [architecture](https://docs.insolar.io/en/latest/architecture.html).

To run the Insolar Platform 1.0 locally, install it and deploy as described below.

### Install

1. Install the latest 1.15 version of the [Golang programming tools](https://golang.org/doc/install#install). Make sure the `$GOPATH` environment variable is set.

2. Download the Insolar package:

   ```
   go get github.com/insolar/insolar
   ```

3. Go to the package directory:

   ```
   cd $GOPATH/src/github.com/insolar/insolar
   ```

4. Install dependencies and build binaries:

   ```
   make
   ```

### Deploy locally

To deploy the Insolar network locally, run the launcher:

```
insolar-scripts/insolard/launchnet.sh -g
```

The launcher generates bootstrap data, starts a pulse watcher, and launches a number of nodes. In local setup, the "nodes" are simply services listening on different ports.
The default number of nodes is 5, you can uncomment more in `scripts/insolard/bootstrap_template.yaml`.

When the pulse watcher says `INSOLAR STATE: READY`, the network is up and running.

## Contribute!

Feel free to submit issues, fork the repository and send pull requests! 

To make the process smooth for both reviewers and contributors, familiarize yourself with the list of guidelines:

1. [Open source contributor guide](https://github.com/freeCodeCamp/how-to-contribute-to-open-source).
2. [Style guide: Effective Go](https://golang.org/doc/effective_go.html).
3. [List of shorthands for Go code review comments](https://github.com/golang/go/wiki/CodeReviewComments).

When submitting an issue, **include a complete test function** that demonstrates it.

Thank you for your intention to contribute to the Insolar project. As a company developing open-source code, we highly appreciate external contributions to our project.

## FAQ

For more information, check out our [FAQ](https://github.com/insolar/insolar/wiki/FAQ).

## Contacts

If you have any additional questions, join our [developers chat](https://t.me/InsolarTech).

Our social media:

[<img src="https://github.com/insolar/doc-pics/raw/master/st/ico-social-facebook.png" width="36" height="36">](https://facebook.com/insolario)
[<img src="https://github.com/insolar/doc-pics/raw/master/st/ico-social-twitter.png" width="36" height="36">](https://twitter.com/insolario)
[<img src="https://github.com/insolar/doc-pics/raw/master/st/ico-social-medium.png" width="36" height="36">](https://medium.com/insolar)
[<img src="https://github.com/insolar/doc-pics/raw/master/st/ico-social-youtube.png" width="36" height="36">](https://youtube.com/insolar)
[<img src="https://github.com/insolar/doc-pics/raw/master/st/ico-social-reddit.png" width="36" height="36">](https://www.reddit.com/r/insolar/)
[<img src="https://github.com/insolar/doc-pics/raw/master/st/ico-social-linkedin.png" width="36" height="36">](https://www.linkedin.com/company/insolario/)
[<img src="https://github.com/insolar/doc-pics/raw/master/st/ico-social-instagram.png" width="36" height="36">](https://instagram.com/insolario)
[<img src="https://github.com/insolar/doc-pics/raw/master/st/ico-social-telegram.png" width="36" height="36">](https://t.me/InsolarAnnouncements)

## License

This project is licensed under the terms of the [Insolar License 1.0](LICENSE.md).
