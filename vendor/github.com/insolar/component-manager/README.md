# Overview 

Component Manager provides dependency injection and lifecycle management for component-based monolith architecture apps.

See [Demo application](https://github.com/AndreyBronin/golang-di-sandbox)

[![Build Status](https://travis-ci.org/insolar/component-manager.svg?branch=master)](https://travis-ci.org/insolar/component-manager)
[![GolangCI](https://golangci.com/badges/github.com/insolar/component-manager.svg)](https://golangci.com/r/github.com/insolar/component-manager/)
[![Go Report Card](https://goreportcard.com/badge/github.com/insolar/component-manager)](https://goreportcard.com/report/github.com/insolar/component-manager)
[![GoDoc](https://godoc.org/github.com/insolar/component-manager?status.svg)](https://godoc.org/github.com/insolar/component-manager)
[![codecov](https://codecov.io/gh/insolar/component-manager/branch/master/graph/badge.svg)](https://codecov.io/gh/insolar/component-manager)


### Features 
- two step initialization
- reflect based dependency injection for interfaces
- resolving circular dependency 
- components lifecycle support
- ordered start, gracefully stop with reverse order
- easy component and integration tests with mock
- subcomponents support
- reduce boilerplate code

## Contetns
- [Basic usage](#basic-usage)
    * [Installing](#installing)
	* [Component definition](#component-definition)
	* [Component lifecycle](#component-lifecycle)


## Basic usage

## Installing
To start using Component Manager, install Go 1.9 or above and run `go get`:

```sh
$ go get github.com/insolar/component-manager
```


## Component definition

A Component is a struct which can have dependencies and/or can implement lifecycle interfaces.

Dependencies defined as fields in the struct and must be an interface type.
have to be exportable because reflect can set only exportable struct fields.
Also Dependencies must have tag `inject:""`.

```go
    type Supermarket struct {
        Warehouse core.Warehouse `inject:""`
    }

	cm := component.NewManager(nil)
	cm.Register(producer.NewFarm(), producer.NewDoorFactory())
	cm.Register(&supermarket.Supermarket{}, &warehouse.Warehouse{})
	cm.Register(NewCustomer("Bob"), NewCustomer("Alice"))
	cm.Inject()
```

## Component lifecycle

Usually components lives from app process executes till process finished. 

- new(instance created, first initialization) 
- inject(required dependency injected)
- init(second initialization)
- start(component can call their dependency interfaces, run goroutines)
- prepare stop(optional)
- stop (gracefully stop goroutines, close descriptors)

### Component constructor

Constructor with config as param.

### Init and start
When should use Init and when Start?
What does it means.

### Stop and gracefully stop

tbd

## intefaces 

```go
type Initer interface {
	Init(ctx context.Context) error
}

type Starter interface {
	Start(ctx context.Context) error
}

type GracefulStopper interface {
	GracefulStop(ctx context.Context) error
}

type Stopper interface {
	Stop(ctx context.Context) error
}
```


## Similar projects

- [facebookgo/inject](https://github.com/facebookgo/inject) - reflect based dependency injector
- [Uber FX](https://github.com/uber-go/fx) - A dependency injection based application framework
- [Google Wire](https://github.com/google/wire) - Compile-time Dependency Injection based on code generation
- [jwells131313/dargo](https://github.com/jwells131313/dargo) - Dependency Injector for GO inspired by Java [HK2](https://javaee.github.io/hk2/)
- [sarulabs/di](https://github.com/sarulabs/di) - Dependency injection framework for go programs
                                                   
