package _example

import (
	"context"
	"strings"
	"time"

	"github.com/insolar/insolar/insolar/belt"
)

type SendCall struct {
	Request string
	Token   string

	// Keep internal state unexported.
	subCallRequest chan string
	done           chan struct{}

	result struct {
		SubCall string
		Result  string
	}
}

func (a *SendCall) Proceed(context.Context) bool {
	if a.done == nil {
		a.done = make(chan struct{})

		// This is a simulation of external call. It will not be in the client code.
		a.subCallRequest = make(chan string)
		go func() {
			for range time.NewTicker(100 * time.Millisecond).C {
				a.subCallRequest <- "this is sa sub-call request received from external source."
			}
		}()

		// Long operation simulation
		go func() {
			<-time.After(time.Second * 2)
			close(a.done)
		}()
		return false
	}

	select {
	case <-a.done:
		return true
	case a.result.SubCall = <-a.subCallRequest:
		return false
	}
}

type GetToken struct {
	Request string

	Result struct {
		Token string
	}
}

func (a *GetToken) Proceed(context.Context) bool {
	a.Result.Token = "I allow!"
	return true
}

type ReturnCall struct {
	Message string
}

func (a *ReturnCall) Proceed(context.Context) bool {
	return true
}

// =====================================================================================================================

type CallMethod struct {
	Request string

	// Keep internal state unexported.
	call     *SendCall
	subCalls []string
}

func (s *CallMethod) Present(ctx context.Context, FLOW belt.Flow) {
	s.call = &SendCall{Request: s.Request}
	// Collect all the sub-calls.
	for FLOW.Yield(s.ContinueWithToken, s.call) {
		s.subCalls = append(s.subCalls, s.call.result.SubCall)
	}
}

func (s *CallMethod) ContinueWithToken(ctx context.Context, FLOW belt.Flow) {
	// Receive token for sub-calls.
	token := &GetToken{Request: s.Request}
	FLOW.Yield(nil, token)
	s.call.Token = token.Result.Token

	// Continue collecting sub-calls.
	for FLOW.Yield(s.ContinueWithToken, s.call) {
		s.subCalls = append(s.subCalls, s.call.result.SubCall)
	}

	FLOW.Yield(nil, &ReturnCall{Message: "My calls: " + strings.Join(s.subCalls, ",")})
}
