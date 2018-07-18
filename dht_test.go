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

package network

import (
	"bytes"
	"errors"
	"log"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/insolar/network/connection"
	"github.com/insolar/network/message"
	"github.com/insolar/network/node"
	"github.com/insolar/network/relay"
	"github.com/insolar/network/routing"
	"github.com/insolar/network/rpc"
	"github.com/insolar/network/store"
	"github.com/insolar/network/transport"

	"github.com/stretchr/testify/assert"
)

func getDefaultCtx(dht *DHT) Context {
	ctx, _ := NewContextBuilder(dht).SetDefaultNode().Build()
	return ctx
}

type mockFuture struct {
	result    chan *message.Message
	actor     *node.Node
	request   *message.Message
	requestID message.RequestID
}

func (f *mockFuture) ID() message.RequestID {
	return f.requestID
}

func (f *mockFuture) Actor() *node.Node {
	return f.actor
}

func (f *mockFuture) Request() *message.Message {
	return f.request
}

func (f *mockFuture) Result() <-chan *message.Message {
	return f.result
}

func (f *mockFuture) SetResult(msg *message.Message) {
	f.result <- msg
}

func (f *mockFuture) Cancel() {}

type mockTransport struct {
	recv     chan *message.Message
	send     chan *message.Message
	dc       chan bool
	msgChan  chan *message.Message
	failNext bool
	sequence *uint64
}

func newMockTransport() *mockTransport {
	net := &mockTransport{
		recv:     make(chan *message.Message),
		send:     make(chan *message.Message),
		dc:       make(chan bool),
		msgChan:  make(chan *message.Message),
		failNext: false,
		sequence: new(uint64),
	}
	return net
}

func (t *mockTransport) Start() error {
	return nil
}

func (t *mockTransport) Stop() {
	close(t.dc)
}

func (t *mockTransport) Close() {
	close(t.recv)
	close(t.send)
	close(t.msgChan)
}

func (t *mockTransport) Stopped() chan bool {
	return t.dc
}

func (t *mockTransport) Messages() chan *message.Message {
	return t.msgChan
}

func (t *mockTransport) failNextSendMessage() {
	t.failNext = true
}

func (t *mockTransport) SendRequest(q *message.Message) (transport.Future, error) {
	id := transport.AtomicLoadAndIncrementUint64(t.sequence)

	if t.failNext {
		t.failNext = false
		return nil, errors.New("MockNetworking Error")
	}
	t.recv <- q

	return &mockFuture{result: t.send, request: q, actor: q.Receiver, requestID: message.RequestID(id)}, nil
}

func (t *mockTransport) SendResponse(requestID message.RequestID, q *message.Message) error {
	if t.failNext {
		t.failNext = false
		return errors.New("MockNetworking Error")
	}
	return nil
}

func mockFindNodeResponse(request *message.Message, nextID []byte) *message.Message {
	r := &message.Message{}
	n := &node.Node{}
	n.ID = request.Sender.ID
	n.Address = request.Sender.Address
	r.Receiver = n
	netAddr, _ := node.NewAddress("0.0.0.0:3001")
	r.Sender = &node.Node{ID: request.Receiver.ID, Address: netAddr}
	r.Type = request.Type
	r.IsResponse = true
	responseData := &message.ResponseDataFindNode{}
	responseData.Closest = []*node.Node{{ID: nextID, Address: netAddr}}
	r.Data = responseData
	return r
}

func mockFindNodeResponseEmpty(request *message.Message) *message.Message {
	r := &message.Message{}
	n := &node.Node{}
	n.ID = request.Sender.ID
	n.Address = request.Sender.Address
	r.Receiver = n
	netAddr, _ := node.NewAddress("0.0.0.0:3001")
	r.Sender = &node.Node{ID: request.Receiver.ID, Address: netAddr}
	r.Type = request.Type
	r.IsResponse = true
	responseData := &message.ResponseDataFindNode{}
	responseData.Closest = []*node.Node{}
	r.Data = responseData
	return r
}

func dhtParams(ids []node.ID, address string) (store.Store, *node.Origin, transport.Transport, rpc.RPC, error) {
	st := store.NewMemoryStore()
	addr, _ := node.NewAddress(address)
	origin, err := node.NewOrigin(ids, addr)
	tp := newMockTransport()
	r := rpc.NewRPC()
	return st, origin, tp, r, err
}

func realDhtParams(ids []node.ID, address string) (store.Store, *node.Origin, transport.Transport, rpc.RPC, error) {
	st := store.NewMemoryStore()
	addr, _ := node.NewAddress(address)
	origin, _ := node.NewOrigin(ids, addr)
	conn, _ := connection.NewConnectionFactory().Create(address)
	tp, err := transport.NewUTPTransport(conn, relay.CreateProxy())
	r := rpc.NewRPC()
	return st, origin, tp, r, err
}

// Creates twenty DHTs and bootstraps each with the previous
// at the end all should know about each other
func TestBootstrapTwentyNodes(t *testing.T) {
	done := make(chan bool)
	port := 3000
	var dhts []*DHT
	for i := 0; i < 20; i++ {
		id, _ := node.NewIDs(1)
		st, s, tp, r, _ := realDhtParams(id, "127.0.0.1:"+strconv.Itoa(port))
		address, _ := node.NewAddress("127.0.0.1:" + strconv.Itoa(port-1))
		bootstrapNode := node.NewNode(address)
		dht, err := NewDHT(st, s, tp, r, &Options{
			BootstrapNodes: []*node.Node{
				bootstrapNode,
			},
		},
			relay.CreateProxy())
		port++
		dhts = append(dhts, dht)
		assert.NoError(t, err)
	}

	for _, dht := range dhts {
		ctx, _ := NewContextBuilder(dht).SetDefaultNode().Build()
		assert.Equal(t, 0, dht.NumNodes(ctx))
		go func(dht *DHT) {
			err := dht.Listen()
			assert.Equal(t, "closed", err.Error())
			done <- true
		}(dht)
		go func(dht *DHT) {
			err := dht.Bootstrap()
			assert.NoError(t, err)
		}(dht)
		time.Sleep(time.Millisecond * 200)
	}

	time.Sleep(time.Millisecond * 2000)

	for _, dht := range dhts {
		assert.Equal(t, 19, dht.NumNodes(getDefaultCtx(dht)))
		dht.Disconnect()
		<-done
	}
}

// Creates two DHTs, bootstrap one using the other, ensure that they both know
// about each other afterwards.
func TestBootstrapTwoNodes(t *testing.T) {
	done := make(chan bool)

	id1, _ := node.NewIDs(1)
	st, s, tp, r, err := realDhtParams(id1, "127.0.0.1:3000")
	dht1, _ := NewDHT(st, s, tp, r, &Options{}, relay.CreateProxy())
	assert.NoError(t, err)

	bootstrapAddr2, _ := node.NewAddress("127.0.0.1:3000")
	st2, s2, tp2, r2, err := realDhtParams(nil, "127.0.0.1:3001")
	dht2, _ := NewDHT(st2, s2, tp2, r2, &Options{
		BootstrapNodes: []*node.Node{
			{
				ID:      id1[0],
				Address: bootstrapAddr2,
			},
		},
	},
		relay.CreateProxy())

	assert.NoError(t, err)
	assert.Equal(t, 0, dht1.NumNodes(getDefaultCtx(dht1)))
	assert.Equal(t, 0, dht2.NumNodes(getDefaultCtx(dht2)))

	go func() {
		go func() {
			err2 := dht2.Bootstrap()
			assert.NoError(t, err2)

			time.Sleep(50 * time.Millisecond)

			dht2.Disconnect()
			dht1.Disconnect()
			done <- true
		}()
		err3 := dht2.Listen()
		assert.Equal(t, "closed", err3.Error())
		done <- true
	}()

	err = dht1.Listen()
	assert.Equal(t, "closed", err.Error())

	assert.Equal(t, 1, dht1.NumNodes(getDefaultCtx(dht1)))
	assert.Equal(t, 1, dht2.NumNodes(getDefaultCtx(dht2)))
	<-done
	<-done
}

// Creates three DHTs, bootstrap B using A, bootstrap C using B. A should know
// about both B and C
func TestBootstrapThreeNodes(t *testing.T) {
	done := make(chan bool)

	id1, _ := node.NewIDs(1)
	st1, s1, tp1, r1, err := realDhtParams(id1, "127.0.0.1:3000")
	assert.NoError(t, err)
	dht1, _ := NewDHT(st1, s1, tp1, r1, &Options{}, relay.CreateProxy())

	id2, _ := node.NewIDs(1)
	st2, s2, tp2, r2, err := realDhtParams(id2, "127.0.0.1:3001")
	assert.NoError(t, err)
	dht2, _ := NewDHT(st2, s2, tp2, r2, &Options{
		BootstrapNodes: []*node.Node{
			{
				ID:      id1[0],
				Address: dht1.origin.Address,
			},
		},
	},
		relay.CreateProxy())

	st3, s3, tp3, r3, err := realDhtParams(nil, "127.0.0.1:3002")
	assert.NoError(t, err)
	dht3, _ := NewDHT(st3, s3, tp3, r3, &Options{
		BootstrapNodes: []*node.Node{
			{
				ID:      id2[0],
				Address: dht2.origin.Address,
			},
		},
	},
		relay.CreateProxy())

	assert.Equal(t, 0, dht1.NumNodes(getDefaultCtx(dht1)))
	assert.Equal(t, 0, dht2.NumNodes(getDefaultCtx(dht2)))
	assert.Equal(t, 0, dht3.NumNodes(getDefaultCtx(dht3)))

	go func(dht1 *DHT, dht2 *DHT, dht3 *DHT) {
		go func(dht1 *DHT, dht2 *DHT, dht3 *DHT) {
			err2 := dht2.Bootstrap()
			assert.NoError(t, err2)

			go func(dht1 *DHT, dht2 *DHT, dht3 *DHT) {
				err3 := dht3.Bootstrap()
				assert.NoError(t, err3)

				time.Sleep(500 * time.Millisecond)

				dht1.Disconnect()

				time.Sleep(100 * time.Millisecond)

				dht2.Disconnect()

				dht3.Disconnect()
				done <- true
			}(dht1, dht2, dht3)

			err4 := dht3.Listen()
			assert.Equal(t, "closed", err4.Error())
			done <- true
		}(dht1, dht2, dht3)
		err5 := dht2.Listen()
		assert.Equal(t, "closed", err5.Error())
		done <- true
	}(dht1, dht2, dht3)

	err = dht1.Listen()
	assert.Equal(t, "closed", err.Error())

	assert.Equal(t, 2, dht1.NumNodes(getDefaultCtx(dht1)))
	assert.Equal(t, 2, dht2.NumNodes(getDefaultCtx(dht2)))
	assert.Equal(t, 2, dht3.NumNodes(getDefaultCtx(dht3)))

	<-done
	<-done
	<-done
}

// Creates two DHTs and bootstraps using only IP:Port. Connecting node should
// ping the first node to find its RequestID
func TestBootstrapNoID(t *testing.T) {
	done := make(chan bool)

	id1, _ := node.NewIDs(1)
	st1, s1, tp1, r1, err := realDhtParams(id1, "0.0.0.0:3000")
	assert.NoError(t, err)
	dht1, _ := NewDHT(st1, s1, tp1, r1, &Options{}, relay.CreateProxy())

	st2, s2, tp2, r2, err := realDhtParams(nil, "0.0.0.0:3001")
	assert.NoError(t, err)
	dht2, _ := NewDHT(st2, s2, tp2, r2, &Options{
		BootstrapNodes: []*node.Node{
			{
				Address: dht1.origin.Address,
			},
		},
	},
		relay.CreateProxy())

	assert.Equal(t, 0, dht1.NumNodes(getDefaultCtx(dht1)))
	assert.Equal(t, 0, dht2.NumNodes(getDefaultCtx(dht2)))

	go func() {
		go func() {
			err2 := dht2.Bootstrap()
			assert.NoError(t, err2)

			time.Sleep(50 * time.Millisecond)

			dht2.Disconnect()
			dht1.Disconnect()
			done <- true
		}()
		err3 := dht2.Listen()
		assert.Equal(t, "closed", err3.Error())
		done <- true
	}()

	err = dht1.Listen()
	assert.Equal(t, "closed", err.Error())

	assert.Equal(t, 1, dht1.NumNodes(getDefaultCtx(dht1)))
	assert.Equal(t, 1, dht2.NumNodes(getDefaultCtx(dht2)))

	<-done
	<-done
}

// Create two DHTs have them connect and bootstrap, then disconnect. Repeat
// 100 times to ensure that we can use the same IP and port without EADDRINUSE
// errors.
func TestReconnect(t *testing.T) {
	for i := 0; i < 100; i++ {
		done := make(chan bool)

		id1, _ := node.NewIDs(1)
		st1, s1, tp1, r1, err := realDhtParams(id1, "127.0.0.1:3000")
		assert.NoError(t, err)
		dht1, _ := NewDHT(st1, s1, tp1, r1, &Options{}, relay.CreateProxy())

		st2, s2, tp2, r2, err := realDhtParams(nil, "127.0.0.1:3001")
		assert.NoError(t, err)
		dht2, _ := NewDHT(st2, s2, tp2, r2, &Options{
			BootstrapNodes: []*node.Node{
				{
					ID:      id1[0],
					Address: dht1.origin.Address,
				},
			},
		},
			relay.CreateProxy())

		assert.Equal(t, 0, dht1.NumNodes(getDefaultCtx(dht1)))

		go func() {
			go func() {
				err2 := dht2.Bootstrap()
				assert.NoError(t, err2)

				dht2.Disconnect()
				dht1.Disconnect()

				done <- true
			}()
			err3 := dht2.Listen()
			assert.Equal(t, "closed", err3.Error())
			done <- true

		}()

		err = dht1.Listen()
		assert.Equal(t, "closed", err.Error())

		assert.Equal(t, 1, dht1.NumNodes(getDefaultCtx(dht1)))
		assert.Equal(t, 1, dht2.NumNodes(getDefaultCtx(dht2)))

		<-done
		<-done
	}
}

// Create two DHTs and have them connect. Send a store message with 100mb
// payload from one node to another. Ensure that the other node now has
// this data in its store.
func TestStoreAndFindLargeValue(t *testing.T) {
	done := make(chan bool)

	id1, _ := node.NewIDs(1)
	st1, s1, tp1, r1, _ := realDhtParams(id1, "127.0.0.1:3000")
	dht1, _ := NewDHT(st1, s1, tp1, r1, &Options{}, relay.CreateProxy())

	st2, s2, tp2, r2, _ := realDhtParams(nil, "127.0.0.1:3001")
	dht2, _ := NewDHT(st2, s2, tp2, r2, &Options{
		BootstrapNodes: []*node.Node{
			{
				ID:      id1[0],
				Address: dht1.origin.Address,
			},
		},
	}, relay.CreateProxy())

	go func() {
		err := dht1.Listen()
		assert.Equal(t, "closed", err.Error())
		done <- true
	}()

	go func() {
		err := dht2.Listen()
		assert.Equal(t, "closed", err.Error())
		done <- true
	}()

	time.Sleep(1 * time.Second)

	dht2.Bootstrap()

	payload := [1000000]byte{}

	key, err := dht1.Store(getDefaultCtx(dht1), payload[:])
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)

	value, exists, err := dht2.Get(getDefaultCtx(dht1), key)
	assert.NoError(t, err)
	assert.Equal(t, true, exists)
	assert.Equal(t, 0, bytes.Compare(payload[:], value))

	dht1.Disconnect()
	dht2.Disconnect()

	<-done
	<-done
}

// Tests sending a message which results in an error when attempting to
// send over uTP
func TestNetworkingSendError(t *testing.T) {
	id := getIDWithValues(0)
	done := make(chan int)

	bootstrapAddr, _ := node.NewAddress("0.0.0.0:3001")
	st, s, tp, r, err := dhtParams([]node.ID{id}, "0.0.0.0:3000")
	assert.NoError(t, err)

	dht, _ := NewDHT(st, s, tp, r, &Options{
		BootstrapNodes: []*node.Node{{
			ID:      getZerodIDWithNthByte(1, byte(255)),
			Address: bootstrapAddr,
		}},
	},
		relay.CreateProxy())
	mockTp := tp.(*mockTransport)

	go func() {
		dht.Listen()
	}()

	go func() {
		v := <-mockTp.recv
		assert.Nil(t, v)
		close(done)
	}()

	mockTp.failNextSendMessage()

	dht.Bootstrap()

	dht.Disconnect()

	<-done
}

// Tests sending a message which results in a successful send, but the node
// never responds
func TestNodeResponseSendError(t *testing.T) {
	id := getIDWithValues(0)
	done := make(chan int)

	bootstrapAddr, _ := node.NewAddress("0.0.0.0:3001")
	st, s, tp, r, err := dhtParams([]node.ID{id}, "0.0.0.0:3000")
	assert.NoError(t, err)

	dht, _ := NewDHT(st, s, tp, r, &Options{
		BootstrapNodes: []*node.Node{{
			ID:      getZerodIDWithNthByte(1, byte(255)),
			Address: bootstrapAddr,
		}},
	},
		relay.CreateProxy())
	mockTp := tp.(*mockTransport)

	queries := 0

	go func() {
		dht.Listen()
	}()

	go func() {
		for {
			request := <-mockTp.recv
			if request == nil {
				return
			}
			if queries == 1 {
				// Don't respond
				close(done)
			} else {
				queries++
				res := mockFindNodeResponse(request, getZerodIDWithNthByte(2, byte(255)))
				mockTp.send <- res
			}
		}
	}()

	dht.Bootstrap()

	assert.Equal(t, 1, dht.tables[0].TotalNodes())

	dht.Disconnect()

	<-done
}

// Tests a bucket refresh by setting a very low RefreshTime value, adding a single
// node to a bucket, and waiting for the refresh message for the bucket
func TestBucketRefresh(t *testing.T) {
	id := getIDWithValues(0)
	done := make(chan int)
	refresh := make(chan int)

	bootstrapAddr, _ := node.NewAddress("0.0.0.0:3001")
	st, s, tp, r, err := dhtParams([]node.ID{id}, "0.0.0.0:3000")
	assert.NoError(t, err)

	dht, _ := NewDHT(st, s, tp, r, &Options{
		RefreshTime: time.Second * 2,
		BootstrapNodes: []*node.Node{{
			ID:      getZerodIDWithNthByte(1, byte(255)),
			Address: bootstrapAddr,
		}},
	},
		relay.CreateProxy())
	mockTp := tp.(*mockTransport)

	queries := 0

	go func() {
		dht.Listen()
	}()

	go func() {
		for {
			request := <-mockTp.recv
			if request == nil {
				close(done)
				return
			}
			queries++

			res := mockFindNodeResponseEmpty(request)
			mockTp.send <- res

			if queries == 2 {
				close(refresh)
			}
		}
	}()

	dht.Bootstrap()

	assert.Equal(t, 1, dht.tables[0].TotalNodes())

	<-refresh

	dht.Disconnect()

	<-done
}

// Tets store replication by setting the ReplicateTime time to a very small value.
// Stores some data, and then expects another store message in ReplicateTime time
func TestStoreReplication(t *testing.T) {
	id := getIDWithValues(0)
	done := make(chan int)
	replicate := make(chan int)

	bootstrapAddr, _ := node.NewAddress("0.0.0.0:3001")
	st, s, tp, r, err := dhtParams([]node.ID{id}, "0.0.0.0:3000")
	assert.NoError(t, err)

	dht, _ := NewDHT(st, s, tp, r, &Options{
		ReplicateTime: time.Second * 2,
		BootstrapNodes: []*node.Node{{
			ID:      getZerodIDWithNthByte(1, byte(255)),
			Address: bootstrapAddr,
		}},
	},
		relay.CreateProxy())
	mockTp := tp.(*mockTransport)

	go func() {
		dht.Listen()
	}()

	stores := 0

	go func() {
		for {
			request := <-mockTp.recv
			if request == nil {
				close(done)
				return
			}

			switch request.Type {
			case message.TypeFindNode:
				res := mockFindNodeResponseEmpty(request)
				mockTp.send <- res
			case message.TypeStore:
				stores++
				d := request.Data.(*message.RequestDataStore)
				assert.Equal(t, []byte("foo"), d.Data)
				if stores >= 2 {
					close(replicate)
				}
			}
		}
	}()

	dht.Bootstrap()

	dht.Store(getDefaultCtx(dht), []byte("foo"))

	<-replicate

	dht.Disconnect()

	<-done
}

// Test Expiration by setting ExpirationTime to a very low value. Store a value,
// and then wait longer than ExpirationTime. The value should no longer exist in
// the store.
func TestStoreExpiration(t *testing.T) {
	id := getIDWithValues(0)

	st, s, tp, r, err := realDhtParams([]node.ID{id}, "0.0.0.0:3000")
	assert.NoError(t, err)

	dht, _ := NewDHT(st, s, tp, r, &Options{
		ExpirationTime: time.Second,
	},
		relay.CreateProxy())

	go func() {
		dht.Listen()
	}()

	k, _ := dht.Store(getDefaultCtx(dht), []byte("foo"))

	v, exists, _ := dht.Get(getDefaultCtx(dht), k)
	assert.Equal(t, true, exists)

	assert.Equal(t, []byte("foo"), v)

	<-time.After(time.Second * 3)

	_, exists, _ = dht.Get(getDefaultCtx(dht), k)

	assert.Equal(t, false, exists)

	dht.Disconnect()
}

// Create a new node and bootstrap it. All nodes in the network know of a
// single node closer to the original node. This continues until every MaxContactsInBucket bucket
// is occupied.
func TestFindNodeAllBuckets(t *testing.T) {
	id := getIDWithValues(0)

	bootstrapAddr, _ := node.NewAddress("0.0.0.0:3001")
	st, s, tp, r, err := dhtParams([]node.ID{id}, "0.0.0.0:3000")
	assert.NoError(t, err)

	dht, _ := NewDHT(st, s, tp, r, &Options{
		BootstrapNodes: []*node.Node{{
			ID:      getZerodIDWithNthByte(0, byte(math.Pow(2, 7))),
			Address: bootstrapAddr,
		}},
	},
		relay.CreateProxy())
	mockTp := tp.(*mockTransport)

	go func() {
		dht.Listen()
	}()

	var k = 0
	var i = 6

	go func() {
		for {
			request := <-mockTp.recv
			if request == nil {
				return
			}

			res := mockFindNodeResponse(request, getZerodIDWithNthByte(k, byte(math.Pow(2, float64(i)))))

			i--
			if i < 0 {
				i = 7
				k++
			}
			if k > 19 {
				k = 19
			}

			mockTp.send <- res
		}
	}()

	dht.Bootstrap()

	for _, v := range dht.tables[0].RoutingTable {
		assert.Equal(t, 1, len(v))
	}

	dht.Disconnect()
}

// Tests timing out of nodes in a bucket. DHT bootstraps networks and learns
// about 20 subsequent nodes in the same bucket. Upon attempting to add the 21st
// node to the now full bucket, we should receive a ping to the very first node
// added in order to determine if it is still alive.
func TestAddNodeTimeout(t *testing.T) {
	id := getIDWithValues(0)
	done := make(chan int)
	pinged := make(chan int)

	bootstrapAddr, _ := node.NewAddress("0.0.0.0:3001")
	st, s, tp, r, err := dhtParams([]node.ID{id}, "0.0.0.0:3000")
	assert.NoError(t, err)

	dht, _ := NewDHT(st, s, tp, r, &Options{
		BootstrapNodes: []*node.Node{{
			ID:      getZerodIDWithNthByte(1, byte(255)),
			Address: bootstrapAddr,
		}},
	},
		relay.CreateProxy())
	mockTp := tp.(*mockTransport)

	go func() {
		dht.Listen()
	}()

	var nodesAdded = 1
	var firstNode []byte
	var lastNode []byte

	go func() {
		for {
			request := <-mockTp.recv
			if request == nil {
				return
			}
			switch request.Type {
			case message.TypeFindNode:
				id := getIDWithValues(0)
				if nodesAdded > routing.MaxContactsInBucket+1 {
					close(done)
					return
				}

				if nodesAdded == 1 {
					firstNode = id
				}

				if nodesAdded == routing.MaxContactsInBucket {
					lastNode = id
				}

				id[1] = byte(255 - nodesAdded)
				nodesAdded++

				res := mockFindNodeResponse(request, id)
				mockTp.send <- res
			case message.TypePing:
				assert.Equal(t, message.TypePing, request.Type)
				assert.Equal(t, getZerodIDWithNthByte(1, byte(255)), request.Receiver.ID)
				close(pinged)
			}
		}
	}()

	dht.Bootstrap()

	// ensure the first node in the table is the second node contacted, and the
	// last is the last node contacted
	assert.Equal(t, 0, bytes.Compare(dht.tables[0].RoutingTable[routing.KeyBitSize-9][0].ID, firstNode))
	assert.Equal(t, 0, bytes.Compare(dht.tables[0].RoutingTable[routing.KeyBitSize-9][19].ID, lastNode))

	<-done
	<-pinged

	dht.Disconnect()
}

func TestGetRandomIDFromBucket(t *testing.T) {
	id := getIDWithValues(0)
	st, s, tp, r, err := realDhtParams([]node.ID{id}, "0.0.0.0:3000")
	assert.NoError(t, err)

	dht, _ := NewDHT(st, s, tp, r, &Options{}, relay.CreateProxy())

	go func() {
		dht.Listen()
	}()

	// Bytes should be equal up to the bucket index that the random RequestID was
	// generated for, and random afterwards
	for i := 0; i < routing.KeyBitSize/8; i++ {
		r := dht.tables[0].GetRandomIDFromBucket(i * 8)
		for j := 0; j < i; j++ {
			assert.Equal(t, byte(0), r[j])
		}
	}

	dht.Disconnect()
}

func getZerodIDWithNthByte(n int, v byte) node.ID {
	id := getIDWithValues(0)
	id[n] = v
	return id
}

func getIDWithValues(b byte) node.ID {
	return []byte{b, b, b, b, b, b, b, b, b, b, b, b, b, b, b, b, b, b, b, b}
}

func TestDHT_RelayRequest(t *testing.T) {
	count := 20
	done := make(chan bool)
	port := 3000
	var dhts []*DHT
	idx := make(map[int]string, count)

	for i := 0; i < count; i++ {
		id, _ := node.NewIDs(1)
		idx[i] = id[0].String()
		st, s, tp, r, _ := realDhtParams(id, "127.0.0.1:"+strconv.Itoa(port))
		address, _ := node.NewAddress("127.0.0.1:" + strconv.Itoa(port-1))
		bootstrapNode := node.NewNode(address)
		dht, err := NewDHT(st, s, tp, r, &Options{
			BootstrapNodes: []*node.Node{
				bootstrapNode,
			},
		},
			relay.CreateProxy())
		port++
		dhts = append(dhts, dht)
		assert.NoError(t, err)
	}

	for _, dht := range dhts {
		ctx, _ := NewContextBuilder(dht).SetDefaultNode().Build()
		assert.Equal(t, 0, dht.NumNodes(ctx))
		go func(dht *DHT) {
			err := dht.Listen()
			assert.Equal(t, "closed", err.Error())
			done <- true
		}(dht)
		go func(dht *DHT) {
			err := dht.Bootstrap()
			assert.NoError(t, err)
		}(dht)
		time.Sleep(time.Millisecond * 200)
	}

	time.Sleep(time.Millisecond * 2000)

	index := 0
	for i := range dhts {
		if (i + 1) == count {
			index = 0
		} else {
			index++
		}

		ctx, _ := NewContextBuilder(dhts[i]).SetDefaultNode().Build()
		err := dhts[index].RelayRequest(ctx, "start", idx[i])
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Millisecond * 200)

		assert.Equal(t, 1, dhts[i].relay.ClientsCount())
		assert.Equal(t, true, dhts[index].proxy.ProxyNodesCount() > 0)
	}

	time.Sleep(time.Millisecond * 2000)

	for _, dht := range dhts {
		assert.Equal(t, count-1, dht.NumNodes(getDefaultCtx(dht)))
		dht.Disconnect()
		<-done
	}
}
