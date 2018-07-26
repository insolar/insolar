Insolar – Platform Genesis
==================================================
Proof-Of-Concept, platform for realizations of dApps.

Overview
--------
Genesis module describes interaction of system components with each other.

Every component of the system is a `SmartContract`. Members of the system are given the opportunity to build their own dApps by publishing smart contracts in `Domain` instances.

Domains define the visibility scope for the child contracts and their interaction policies. Actually, `Domain` is subclass of `SmartContract`.


Base domains
--------------
### [ReferenceDomain](https://godoc.org/github.com/insolar/insolar/genesis/public/core/reference.go)

System domain that allow users to publish their `SmartContract`'s for public use. 
After publication, the reference becomes available for resolving by global resolver.

Usage:
```
factory := NewReferenceDomainFactory()
refDomain := factory.Create(nil)

record, err := refDomain.RegisterReference(someReference)
resolved, err := refDomain.ResolveReference(record)
```

Resolver mechanism
--------------
Resolver makes possible for smart contracts to interact with each other.
A method for resolving is provided for every type of reference.

### [Global](https://godoc.org/github.com/insolar/insolar/genesis/model/resolver/global.go)
Resolve references of `GlobalScope` type.
`GlobalScope` references are global reference, must be resolved via `ReferenceDomain`.
 Object with this type of reference consider as public. 


### [Context](https://godoc.org/github.com/insolar/insolar/genesis/model/resolver/context.go)
Resolve references of `ContextScope` type.
`ContextScope` references point to an object within a context provided by parent.
It may be a child or any other object available to the parent.


### [Child](https://godoc.org/github.com/insolar/insolar/genesis/model/resolver/child.go)
Resolve references of `ChildScope` type.
`ChildScope` references to a child of an object which resolving the reference.
