// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proc

type Dependencies struct {
	PassState        func(*PassState)
	SendCode         func(*SendCode)
	SendRequests     func(*SendRequests)
	SendRequest      func(*SendRequest)
	Replication      func(*Replication)
	SendJet          func(*SendJet)
	SendIndex        func(*SendIndex)
	SearchIndex      func(*SearchIndex)
	SendInitialState func(*SendInitialState)
	SendPulse        func(*SendPulse)
}
