/*
 *    Copyright 2018 INS Ecosystem
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package messagerouter

import (
	"errors"
	"testing"

	"github.com/insolar/insolar/network/host"
	"github.com/insolar/insolar/network/host/connection"
	"github.com/insolar/insolar/network/host/id"
	"github.com/insolar/insolar/network/host/node"
	"github.com/insolar/insolar/network/host/relay"
	"github.com/insolar/insolar/network/host/rpc"
	"github.com/insolar/insolar/network/host/store"
	"github.com/insolar/insolar/network/host/transport"
	"github.com/stretchr/testify/assert"
)

type req struct {
	ref    string
	method string
	args   []byte
}

type resp struct {
	data []byte
	res  []byte
	err  error
}

type runner struct {
	requests  []req
	responses []resp
}

const closedMessage = "closed" // "broken pipe" for kcpTransport

func dhtParams(ids []id.ID, address string) (store.Store, *node.Origin, transport.Transport, rpc.RPC, error) {
	st := store.NewMemoryStore()
	addr, _ := node.NewAddress(address)
	origin, _ := node.NewOrigin(ids, addr)
	conn, _ := connection.NewConnectionFactory().Create(address)
	tp, err := transport.NewUTPTransport(conn, relay.NewProxy())
	r := rpc.NewRPC()
	return st, origin, tp, r, err
}

func getDefaultCtx(dht *host.DHT) host.Context {
	ctx, _ := host.NewContextBuilder(dht).SetDefaultNode().Build()
	return ctx
}

func NewNode() (*host.DHT, error) {
	var ids []id.ID
	id1, _ := id.NewID(nil)
	ids = append(ids, id1)
	st, s, tp, r, err := dhtParams(ids, "127.0.0.1:16000")
	if err != nil {
		return nil, err
	}

	return host.NewDHT(st, s, tp, r, &host.Options{}, relay.NewProxy())
}

type mockRpc struct {
}

func (r *mockRpc) RemoteProcedureCall(ctx host.Context, target string, method string, args [][]byte) (result []byte, err error) {
	return nil, errors.New("not implemented in mock")
}

func (r *mockRpc) RemoteProcedureRegister(name string, method host.RemoteProcedure) {
	return
}

func (r *runner) Execute(ref string, method string, args []byte) ([]byte, []byte, error) {
	if len(r.responses) == 0 {
		panic("no request expected")
	}

	r.requests = append(r.requests, req{ref, method, args})

	resp := r.responses[0]
	r.responses = r.responses[1:]

	return resp.data, resp.res, resp.err
}

func TestNew(t *testing.T) {
	mr, err := New(new(runner), new(mockRpc))
	if err != nil {
		t.Fatal(err)
	}
	if mr == nil {
		t.Fatal("no object created")
	}
}

func TestRoute(t *testing.T) {
	r := new(runner)
	r.requests = make([]req, 0)
	r.responses = make([]resp, 0)

	dht, _ := NewNode()
	ctx := getDefaultCtx(dht)

	mr, _ := New(r, dht)
	reference := dht.GetOriginID(ctx)

	t.Run("success", func(t *testing.T) {
		r.responses = append(r.responses, resp{[]byte("data"), []byte("result"), nil})
		resp, err := mr.Route(ctx, Message{Reference: reference, Method: "SomeMethod", Arguments: []byte("args")})
		if err != nil {
			t.Fatal(err)
		}
		if string(resp.Data) != "data" {
			t.Fatal("unexpected data")
		}
		if string(resp.Result) != "result" {
			t.Fatal("unexpected data")
		}
		if len(r.requests) != 1 {
			t.Fatal("unexpected number of requests registered")
		}
		req := r.requests[0]
		r.requests = r.requests[1:]

		if req.ref != reference {
			t.Fatal("unexpected data")
		}
		if req.method != "SomeMethod" {
			t.Fatal("unexpected data")
		}
		if string(req.args) != "args" {
			t.Fatal("unexpected data")
		}
	})
	t.Run("error", func(t *testing.T) {
		r.responses = append(r.responses, resp{[]byte{}, []byte{}, errors.New("wtf")})
		_, err := mr.Route(ctx, Message{Reference: reference, Method: "SomeMethod", Arguments: []byte("args")})
		if err == nil {
			t.Fatal("error expected")
		}

		if len(r.requests) != 1 {
			t.Fatal("unexpected number of requests registered")
		}
		req := r.requests[0]
		r.requests = r.requests[1:]

		if req.ref != reference {
			t.Fatal("unexpected data")
		}
		if req.method != "SomeMethod" {
			t.Fatal("unexpected data")
		}
		if string(req.args) != "args" {
			t.Fatal("unexpected data")
		}
	})

	t.Run("referenceNotFound", func(t *testing.T) {
		_, err := mr.Route(ctx, Message{Reference: "refNotFound", Method: "SomeMethod", Arguments: []byte("args")})
		assert.Error(t, err)
	})
}
