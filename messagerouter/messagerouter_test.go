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
	mr, err := New(new(runner))
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

	mr, _ := New(r)

	t.Run("success", func(t *testing.T) {
		r.responses = append(r.responses, resp{[]byte("data"), []byte("result"), nil})
		resp, err := mr.Route(Message{Reference: "some.ref", Method: "SomeMethod", Arguments: []byte("args")})
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

		if req.ref != "some.ref" {
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
		_, err := mr.Route(Message{Reference: "some.ref", Method: "SomeMethod", Arguments: []byte("args")})
		if err == nil {
			t.Fatal("error expected")
		}

		if len(r.requests) != 1 {
			t.Fatal("unexpected number of requests registered")
		}
		req := r.requests[0]
		r.requests = r.requests[1:]

		if req.ref != "some.ref" {
			t.Fatal("unexpected data")
		}
		if req.method != "SomeMethod" {
			t.Fatal("unexpected data")
		}
		if string(req.args) != "args" {
			t.Fatal("unexpected data")
		}
	})
}
