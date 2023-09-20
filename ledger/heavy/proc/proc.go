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
