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

package host

import (
	"bytes"
	"errors"
	"math"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/insolar/insolar/network/host/connection"
	"github.com/insolar/insolar/network/host/id"
	"github.com/insolar/insolar/network/host/node"
	"github.com/insolar/insolar/network/host/packet"
	"github.com/insolar/insolar/network/host/relay"
	"github.com/insolar/insolar/network/host/routing"
	"github.com/insolar/insolar/network/host/rpc"
	"github.com/insolar/insolar/network/host/store"
	"github.com/insolar/insolar/network/host/transport"

	"github.com/stretchr/testify/assert"
)

const closedPacket = "closed" // "broken pipe" for kcpTransport

func getDefaultCtx(dht *DHT) Context {
	ctx, _ := NewContextBuilder(dht).SetDefaultNode().Build()
	return ctx
}

type mockFuture struct {
	result    chan *packet.Packet
	actor     *node.Node
	request   *packet.Packet
	requestID packet.RequestID
}

func (f *mockFuture) ID() packet.RequestID {
	return f.requestID
}

func (f *mockFuture) Actor() *node.Node {
	return f.actor
}

func (f *mockFuture) Request() *packet.Packet {
	return f.request
}

func (f *mockFuture) Result() <-chan *packet.Packet {
	return f.result
}

func (f *mockFuture) SetResult(msg *packet.Packet) {
	f.result <- msg
}

func (f *mockFuture) Cancel() {}

type mockTransport struct {
	recv     chan *packet.Packet
	send     chan *packet.Packet
	dc       chan bool
	msgChan  chan *packet.Packet
	failNext bool
	sequence *uint64
}

func newMockTransport() *mockTransport {
	net := &mockTransport{
		recv:     make(chan *packet.Packet),
		send:     make(chan *packet.Packet),
		dc:       make(chan bool),
		msgChan:  make(chan *packet.Packet),
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

func (t *mockTransport) Stopped() <-chan bool {
	return t.dc
}

func (t *mockTransport) Packets() <-chan *packet.Packet {
	return t.msgChan
}

func (t *mockTransport) failNextSendPacket() {
	t.failNext = true
}

func (t *mockTransport) SendRequest(q *packet.Packet) (transport.Future, error) {
	sequenceNumber := transport.AtomicLoadAndIncrementUint64(t.sequence)

	if t.failNext {
		t.failNext = false
		return nil, errors.New("MockNetworking Error")
	}
	t.recv <- q

	return &mockFuture{result: t.send, request: q, actor: q.Receiver, requestID: packet.RequestID(sequenceNumber)}, nil
}

func (t *mockTransport) SendResponse(requestID packet.RequestID, q *packet.Packet) error {
	if t.failNext {
		t.failNext = false
		return errors.New("MockNetworking Error")
	}
	return nil
}

func mockFindNodeResponse(request *packet.Packet, nextID []byte) *packet.Packet {
	r := &packet.Packet{}
	n := &node.Node{}
	n.ID = request.Sender.ID
	n.Address = request.Sender.Address
	r.Receiver = n
	netAddr, _ := node.NewAddress("0.0.0.0:13001")
	r.Sender = &node.Node{ID: request.Receiver.ID, Address: netAddr}
	r.Type = request.Type
	r.IsResponse = true
	responseData := &packet.ResponseDataFindNode{}
	id1, _ := id.NewID(id.GetRandomKey())
	id1.SetHash(nextID)
	responseData.Closest = []*node.Node{{ID: id1, Address: netAddr}}
	r.Data = responseData
	return r
}

func mockFindNodeResponseEmpty(request *packet.Packet) *packet.Packet {
	r := &packet.Packet{}
	n := &node.Node{}
	n.ID = request.Sender.ID
	n.Address = request.Sender.Address
	r.Receiver = n
	netAddr, _ := node.NewAddress("0.0.0.0:14001")
	r.Sender = &node.Node{ID: request.Receiver.ID, Address: netAddr}
	r.Type = request.Type
	r.IsResponse = true
	responseData := &packet.ResponseDataFindNode{}
	responseData.Closest = []*node.Node{}
	r.Data = responseData
	return r
}

func dhtParams(ids []id.ID, address string) (store.Store, *node.Origin, transport.Transport, rpc.RPC, error) {
	st := store.NewMemoryStore()
	addr, _ := node.NewAddress(address)
	origin, err := node.NewOrigin(ids, addr)
	tp := newMockTransport()
	r := rpc.NewRPC()
	return st, origin, tp, r, err
}

func realDhtParams(ids []id.ID, address string) (store.Store, *node.Origin, transport.Transport, rpc.RPC, error) {
	st := store.NewMemoryStore()
	addr, _ := node.NewAddress(address)
	origin, _ := node.NewOrigin(ids, addr)
	conn, _ := connection.NewConnectionFactory().Create(address)
	tp, err := transport.NewUTPTransport(conn, relay.NewProxy())
	r := rpc.NewRPC()
	return st, origin, tp, r, err
}

// Creates twenty DHTs and bootstraps each with the previous
// at the end all should know about each other
func TestBootstrapTwentyNodes(t *testing.T) {
	done := make(chan bool)
	port := 15000
	var dhts []*DHT
	count := 10

	for i := 0; i < count; i++ {
		ids := make([]id.ID, 0)
		id1, _ := id.NewID(id.GetRandomKey())
		ids = append(ids, id1)
		st, s, tp, r, err := realDhtParams(ids, "127.0.0.1:"+strconv.Itoa(port))
		assert.NoError(t, err)
		address, err := node.NewAddress("127.0.0.1:" + strconv.Itoa(port-1))
		assert.NoError(t, err)
		bootstrapNode := node.NewNode(address)
		dht, err := NewDHT(st, s, tp, r, &Options{
			BootstrapNodes: []*node.Node{
				bootstrapNode,
			},
		},
			relay.NewProxy())
		port++
		assert.NoError(t, err)
		dhts = append(dhts, dht)
	}

	for _, dht := range dhts {
		ctx, _ := NewContextBuilder(dht).SetDefaultNode().Build()
		assert.Equal(t, 0, dht.NumNodes(ctx))
		go func(dht *DHT) {
			err := dht.Listen()
			assert.Equal(t, closedPacket, err.Error())
			done <- true
		}(dht)
		go func(dht *DHT) {
			err := dht.Bootstrap()
			assert.NoError(t, err)
			time.Sleep(time.Millisecond * 200)
		}(dht)
		time.Sleep(time.Millisecond * 200)
	}

	time.Sleep(time.Millisecond * 3000)

	for _, dht := range dhts {
		assert.Equal(t, count-1, dht.NumNodes(getDefaultCtx(dht)))
		dht.Disconnect()
		<-done
	}
}

// Creates two DHTs, bootstrap one using the other, ensure that they both know
// about each other afterwards.
func TestBootstrapTwoNodes(t *testing.T) {
	done := make(chan bool)

	ids1 := make([]id.ID, 0)
	id1, _ := id.NewID(id.GetRandomKey())
	ids1 = append(ids1, id1)
	st, s, tp, r, err := realDhtParams(ids1, "127.0.0.1:16000")
	dht1, _ := NewDHT(st, s, tp, r, &Options{}, relay.NewProxy())
	assert.NoError(t, err)

	bootstrapAddr2, _ := node.NewAddress("127.0.0.1:16000")
	st2, s2, tp2, r2, err := realDhtParams(nil, "127.0.0.1:16001")
	dht2, _ := NewDHT(st2, s2, tp2, r2, &Options{
		BootstrapNodes: []*node.Node{
			{
				ID:      ids1[0],
				Address: bootstrapAddr2,
			},
		},
	},
		relay.NewProxy())

	assert.NoError(t, err)
	assert.Equal(t, 0, dht1.NumNodes(getDefaultCtx(dht1)))
	assert.Equal(t, 0, dht2.NumNodes(getDefaultCtx(dht2)))

	go func() {
		go func() {
			err2 := dht2.Bootstrap()
			assert.NoError(t, err2)

			time.Sleep(500 * time.Millisecond)

			dht2.Disconnect()
			dht1.Disconnect()
			done <- true
		}()
		err3 := dht2.Listen()
		assert.Equal(t, closedPacket, err3.Error())
		done <- true
	}()

	err = dht1.Listen()
	assert.Equal(t, closedPacket, err.Error())

	assert.Equal(t, 1, dht1.NumNodes(getDefaultCtx(dht1)))
	assert.Equal(t, 1, dht2.NumNodes(getDefaultCtx(dht2)))
	<-done
	<-done
}

// Creates three DHTs, bootstrap B using A, bootstrap C using B. A should know
// about both B and C
func TestBootstrapThreeNodes(t *testing.T) {
	done := make(chan bool)

	ids1 := make([]id.ID, 0)
	id1, _ := id.NewID(id.GetRandomKey())
	ids1 = append(ids1, id1)
	st1, s1, tp1, r1, err := realDhtParams(ids1, "127.0.0.1:17000")
	assert.NoError(t, err)
	dht1, _ := NewDHT(st1, s1, tp1, r1, &Options{}, relay.NewProxy())

	ids2 := make([]id.ID, 0)
	id2, _ := id.NewID(id.GetRandomKey())
	ids2 = append(ids2, id2)
	st2, s2, tp2, r2, err := realDhtParams(ids2, "127.0.0.1:17001")
	assert.NoError(t, err)
	dht2, _ := NewDHT(st2, s2, tp2, r2, &Options{
		BootstrapNodes: []*node.Node{
			{
				ID:      ids1[0],
				Address: dht1.origin.Address,
			},
		},
	},
		relay.NewProxy())

	st3, s3, tp3, r3, err := realDhtParams(nil, "127.0.0.1:17002")
	assert.NoError(t, err)
	dht3, _ := NewDHT(st3, s3, tp3, r3, &Options{
		BootstrapNodes: []*node.Node{
			{
				ID:      ids2[0],
				Address: dht2.origin.Address,
			},
		},
	},
		relay.NewProxy())

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
			assert.Equal(t, closedPacket, err4.Error())
			done <- true
		}(dht1, dht2, dht3)
		err5 := dht2.Listen()
		assert.Equal(t, closedPacket, err5.Error())
		done <- true
	}(dht1, dht2, dht3)

	err = dht1.Listen()
	assert.Equal(t, closedPacket, err.Error())

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

	ids1 := make([]id.ID, 0)
	id1, _ := id.NewID(id.GetRandomKey())
	ids1 = append(ids1, id1)
	st1, s1, tp1, r1, err := realDhtParams(ids1, "0.0.0.0:18000")
	assert.NoError(t, err)
	dht1, _ := NewDHT(st1, s1, tp1, r1, &Options{}, relay.NewProxy())

	st2, s2, tp2, r2, err := realDhtParams(nil, "0.0.0.0:18001")
	assert.NoError(t, err)
	dht2, _ := NewDHT(st2, s2, tp2, r2, &Options{
		BootstrapNodes: []*node.Node{
			{
				Address: dht1.origin.Address,
			},
		},
	},
		relay.NewProxy())

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
		assert.Equal(t, closedPacket, err3.Error())
		done <- true
	}()

	err = dht1.Listen()
	assert.Equal(t, closedPacket, err.Error())

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

		ids1 := make([]id.ID, 0)
		id1, _ := id.NewID(id.GetRandomKey())
		ids1 = append(ids1, id1)
		st1, s1, tp1, r1, err := realDhtParams(ids1, "127.0.0.1:19000")
		assert.NoError(t, err)
		dht1, _ := NewDHT(st1, s1, tp1, r1, &Options{}, relay.NewProxy())

		st2, s2, tp2, r2, err := realDhtParams(nil, "127.0.0.1:19001")
		assert.NoError(t, err)
		dht2, _ := NewDHT(st2, s2, tp2, r2, &Options{
			BootstrapNodes: []*node.Node{
				{
					ID:      ids1[0],
					Address: dht1.origin.Address,
				},
			},
		},
			relay.NewProxy())

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
			assert.Equal(t, closedPacket, err3.Error())
			done <- true

		}()

		err = dht1.Listen()
		assert.Equal(t, closedPacket, err.Error())

		assert.Equal(t, 1, dht1.NumNodes(getDefaultCtx(dht1)))
		assert.Equal(t, 1, dht2.NumNodes(getDefaultCtx(dht2)))

		<-done
		<-done

		time.Sleep(time.Millisecond * 50)
	}
}

// Create two DHTs and have them connect. Send a store packet with 100mb
// payload from one node to another. Ensure that the other node now has
// this data in its store.
func TestStoreAndFindLargeValue(t *testing.T) {
	done := make(chan bool)

	ids1 := make([]id.ID, 0)
	id1, _ := id.NewID(id.GetRandomKey())
	ids1 = append(ids1, id1)
	st1, s1, tp1, r1, _ := realDhtParams(ids1, "127.0.0.1:20000")
	dht1, _ := NewDHT(st1, s1, tp1, r1, &Options{}, relay.NewProxy())

	st2, s2, tp2, r2, _ := realDhtParams(nil, "127.0.0.1:20001")
	dht2, _ := NewDHT(st2, s2, tp2, r2, &Options{
		BootstrapNodes: []*node.Node{
			{
				ID:      ids1[0],
				Address: dht1.origin.Address,
			},
		},
	}, relay.NewProxy())

	go func() {
		err := dht1.Listen()
		assert.Equal(t, closedPacket, err.Error())
		done <- true
	}()

	go func() {
		err := dht2.Listen()
		assert.Equal(t, closedPacket, err.Error())
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

// Tests sending a packet which results in an error when attempting to
// send over uTP
func TestNetworkingSendError(t *testing.T) {
	zeroId := getIDWithValues(0)
	done := make(chan int)

	bootstrapAddr, _ := node.NewAddress("0.0.0.0:21001")
	st, s, tp, r, err := dhtParams([]id.ID{zeroId}, "0.0.0.0:21000")
	assert.NoError(t, err)

	dht, _ := NewDHT(st, s, tp, r, &Options{
		BootstrapNodes: []*node.Node{{
			ID:      getZerodIDWithNthByte(1, byte(255)),
			Address: bootstrapAddr,
		}},
	},
		relay.NewProxy())
	mockTp := tp.(*mockTransport)

	go func() {
		dht.Listen()
	}()

	go func() {
		v := <-mockTp.recv
		assert.Nil(t, v)
		close(done)
	}()

	mockTp.failNextSendPacket()

	dht.Bootstrap()

	dht.Disconnect()

	<-done
}

// Tests sending a packet which results in a successful send, but the node
// never responds
func TestNodeResponseSendError(t *testing.T) {
	zeroID := getIDWithValues(0)
	done := make(chan int)

	bootstrapAddr, _ := node.NewAddress("0.0.0.0:22001")
	st, s, tp, r, err := dhtParams([]id.ID{zeroID}, "0.0.0.0:22000")
	assert.NoError(t, err)

	dht, _ := NewDHT(st, s, tp, r, &Options{
		BootstrapNodes: []*node.Node{{
			ID:      getZerodIDWithNthByte(1, byte(255)),
			Address: bootstrapAddr,
		}},
	},
		relay.NewProxy())
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
				res := mockFindNodeResponse(request, getZerodIDWithNthByte(2, byte(255)).GetHash())
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
// node to a bucket, and waiting for the refresh packet for the bucket
func TestBucketRefresh(t *testing.T) {
	zeroID := getIDWithValues(0)
	done := make(chan int)
	refresh := make(chan int)

	bootstrapAddr, _ := node.NewAddress("0.0.0.0:23001")
	st, s, tp, r, err := dhtParams([]id.ID{zeroID}, "0.0.0.0:23000")
	assert.NoError(t, err)

	dht, _ := NewDHT(st, s, tp, r, &Options{
		RefreshTime: time.Second * 2,
		BootstrapNodes: []*node.Node{{
			ID:      getZerodIDWithNthByte(1, byte(255)),
			Address: bootstrapAddr,
		}},
	},
		relay.NewProxy())
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
// Stores some data, and then expects another store packet in ReplicateTime time
func TestStoreReplication(t *testing.T) {
	zeroID := getIDWithValues(0)
	done := make(chan int)
	replicate := make(chan int)

	bootstrapAddr, _ := node.NewAddress("0.0.0.0:24001")
	st, s, tp, r, err := dhtParams([]id.ID{zeroID}, "0.0.0.0:24000")
	assert.NoError(t, err)

	dht, _ := NewDHT(st, s, tp, r, &Options{
		ReplicateTime: time.Second * 2,
		BootstrapNodes: []*node.Node{{
			ID:      getZerodIDWithNthByte(1, byte(255)),
			Address: bootstrapAddr,
		}},
	},
		relay.NewProxy())
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
			case packet.TypeFindNode:
				res := mockFindNodeResponseEmpty(request)
				mockTp.send <- res
			case packet.TypeStore:
				stores++
				d := request.Data.(*packet.RequestDataStore)
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
	done := make(chan bool)
	zeroID := getIDWithValues(0)

	st, s, tp, r, err := realDhtParams([]id.ID{zeroID}, "0.0.0.0:25000")
	assert.NoError(t, err)

	dht, _ := NewDHT(st, s, tp, r, &Options{
		ExpirationTime: time.Second,
	},
		relay.NewProxy())

	go func() {
		dht.Listen()
		done <- true
	}()

	k, _ := dht.Store(getDefaultCtx(dht), []byte("foo"))

	v, exists, _ := dht.Get(getDefaultCtx(dht), k)
	assert.Equal(t, true, exists)

	assert.Equal(t, []byte("foo"), v)

	<-time.After(time.Second * 3)

	_, exists, _ = dht.Get(getDefaultCtx(dht), k)

	assert.Equal(t, false, exists)

	dht.Disconnect()
	<-done
}

// Create a new node and bootstrap it. All nodes in the network know of a
// single node closer to the original node. This continues until every MaxContactsInBucket bucket
// is occupied.
func TestFindNodeAllBuckets(t *testing.T) {
	done := make(chan bool)
	zeroID := getIDWithValues(0)

	bootstrapAddr, _ := node.NewAddress("0.0.0.0:26001")
	st, s, tp, r, err := dhtParams([]id.ID{zeroID}, "0.0.0.0:26000")
	assert.NoError(t, err)

	dht, _ := NewDHT(st, s, tp, r, &Options{
		BootstrapNodes: []*node.Node{{
			ID:      getZerodIDWithNthByte(0, byte(math.Pow(2, 7))),
			Address: bootstrapAddr,
		}},
	},
		relay.NewProxy())
	mockTp := tp.(*mockTransport)

	go func() {
		dht.Listen()
		done <- true
	}()

	var k = 0
	var i = 6

	go func() {
		for {
			request := <-mockTp.recv
			if request == nil {
				return
			}

			res := mockFindNodeResponse(request, getZerodIDWithNthByte(k, byte(math.Pow(2, float64(i)))).GetHash())

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
	<-done
}

// Tests timing out of nodes in a bucket. DHT bootstraps networks and learns
// about 20 subsequent nodes in the same bucket. Upon attempting to add the 21st
// node to the now full bucket, we should receive a ping to the very first node
// added in order to determine if it is still alive.
func TestAddNodeTimeout(t *testing.T) {
	zeroID := getIDWithValues(0)
	done := make(chan int)
	pinged := make(chan int)

	bootstrapAddr, _ := node.NewAddress("0.0.0.0:27001")
	st, s, tp, r, err := dhtParams([]id.ID{zeroID}, "0.0.0.0:27000")
	assert.NoError(t, err)

	dht, _ := NewDHT(st, s, tp, r, &Options{
		BootstrapNodes: []*node.Node{{
			ID:      getZerodIDWithNthByte(1, byte(255)),
			Address: bootstrapAddr,
		}},
	},
		relay.NewProxy())
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
			case packet.TypeFindNode:
				id1 := getIDWithValues(0)
				if nodesAdded > routing.MaxContactsInBucket+1 {
					close(done)
					return
				}

				if nodesAdded == 1 {
					firstNode = id1.GetHash()
				}

				if nodesAdded == routing.MaxContactsInBucket {
					lastNode = id1.GetHash()
				}

				id1.GetHash()[1] = byte(255 - nodesAdded)
				nodesAdded++

				res := mockFindNodeResponse(request, id1.GetHash())
				mockTp.send <- res
			case packet.TypePing:
				assert.Equal(t, packet.TypePing, request.Type)
				assert.Equal(t, getZerodIDWithNthByte(1, byte(255)), request.Receiver.ID)
				close(pinged)
			}
		}
	}()

	dht.Bootstrap()

	// ensure the first node in the table is the second node contacted, and the
	// last is the last node contacted
	assert.Equal(t, 0, bytes.Compare(dht.tables[0].RoutingTable[routing.KeyBitSize-9][0].ID.GetHash(), firstNode))
	assert.Equal(t, 0, bytes.Compare(dht.tables[0].RoutingTable[routing.KeyBitSize-9][19].ID.GetHash(), lastNode))

	dht.Disconnect()

	<-done
	<-pinged
}

func TestGetRandomIDFromBucket(t *testing.T) {
	zeroID := getIDWithValues(0)
	st, s, tp, r, err := realDhtParams([]id.ID{zeroID}, "0.0.0.0:28000")
	assert.NoError(t, err)
	done := make(chan bool)

	dht, _ := NewDHT(st, s, tp, r, &Options{}, relay.NewProxy())

	go func() {
		dht.Listen()
		done <- true
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
	<-done
}

func getZerodIDWithNthByte(n int, v byte) id.ID {
	id1 := getIDWithValues(0)
	id1.GetHash()[n] = v
	return id1
}

func getIDWithValues(b byte) id.ID {
	id1, _ := id.NewID(nil)
	id1.SetHash([]byte{b, b, b, b, b, b, b, b, b, b, b, b, b, b, b, b, b, b, b, b})
	return id1
}

func TestDHT_FindNode(t *testing.T) {
	count := 5
	port := 6000
	var dhts []*DHT
	idx := make(map[int]string, count)
	done := make(chan bool)

	for i := 0; i < count; i++ {
		ids1 := make([]id.ID, 0)
		id1, _ := id.NewID(id.GetRandomKey())
		ids1 = append(ids1, id1)
		idx[i] = ids1[0].HashString()
		st, s, tp, r, _ := realDhtParams(ids1, "127.0.0.1:"+strconv.Itoa(port))
		address, _ := node.NewAddress("127.0.0.1:" + strconv.Itoa(port-1))
		bootstrapNode := node.NewNode(address)
		dht, err := NewDHT(st, s, tp, r, &Options{
			BootstrapNodes: []*node.Node{
				bootstrapNode,
			},
		},
			relay.NewProxy())
		port++
		dhts = append(dhts, dht)
		assert.NoError(t, err)
	}

	for _, dht := range dhts {
		ctx, _ := NewContextBuilder(dht).SetDefaultNode().Build()
		assert.Equal(t, 0, dht.NumNodes(ctx))
		go func(dht *DHT) {
			err := dht.Listen()
			assert.Equal(t, closedPacket, err.Error())
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
		for i := range idx {
			ctx, _ := NewContextBuilder(dht).SetDefaultNode().Build()
			_, exist, err := dht.FindNode(ctx, idx[i])
			assert.NoError(t, err)
			assert.Equal(t, true, exist)
			time.Sleep(time.Millisecond * 200)
		}
	}

	for _, dht := range dhts {
		assert.Equal(t, count-1, dht.NumNodes(getDefaultCtx(dht)))
		dht.Disconnect()
	}
	<-done
}

func TestDHT_Listen(t *testing.T) {
	count := 5
	port := 8000
	var dhts []*DHT
	done := make(chan bool)

	for i := 0; i < count; i++ {
		ids1 := make([]id.ID, 0)
		id1, _ := id.NewID(id.GetRandomKey())
		ids1 = append(ids1, id1)
		st, s, tp, r, _ := realDhtParams(ids1, "127.0.0.1:"+strconv.Itoa(port))
		address, _ := node.NewAddress("127.0.0.1:" + strconv.Itoa(port-1))
		bootstrapNode := node.NewNode(address)
		dht, err := NewDHT(st, s, tp, r, &Options{
			BootstrapNodes: []*node.Node{
				bootstrapNode,
			},
		},
			relay.NewProxy())
		port++
		dhts = append(dhts, dht)
		assert.NoError(t, err)
	}

	for _, dht := range dhts {
		ctx, _ := NewContextBuilder(dht).SetDefaultNode().Build()
		assert.Equal(t, 0, dht.NumNodes(ctx))
		go func(dht *DHT) {
			err := dht.Listen()
			assert.Equal(t, closedPacket, err.Error())
			done <- true
		}(dht)
	}

	for _, dht := range dhts {
		dht.Disconnect()
	}
	<-done
}

func TestDHT_Disconnect(t *testing.T) {
	count := 5
	port := 9000
	var dhts []*DHT
	done := make(chan bool)

	for i := 0; i < count; i++ {
		ids1 := make([]id.ID, 0)
		id1, _ := id.NewID(id.GetRandomKey())
		ids1 = append(ids1, id1)
		st, s, tp, r, _ := realDhtParams(ids1, "127.0.0.1:"+strconv.Itoa(port))
		address, _ := node.NewAddress("127.0.0.1:" + strconv.Itoa(port-1))
		bootstrapNode := node.NewNode(address)
		dht, err := NewDHT(st, s, tp, r, &Options{
			BootstrapNodes: []*node.Node{
				bootstrapNode,
			},
		},
			relay.NewProxy())
		port++
		dhts = append(dhts, dht)
		assert.NoError(t, err)
	}

	for _, dht := range dhts {
		ctx, _ := NewContextBuilder(dht).SetDefaultNode().Build()
		assert.Equal(t, 0, dht.NumNodes(ctx))
		go func(dht *DHT) {
			err := dht.Listen()
			assert.Equal(t, closedPacket, err.Error())
			done <- true
		}(dht)
		time.Sleep(time.Millisecond * 200)
	}

	for _, dht := range dhts {
		dht.Disconnect()
	}
	<-done
}

func TestNewDHT(t *testing.T) {
	done := make(chan bool)
	port := 11000
	ids1 := make([]id.ID, 0)
	id1, _ := id.NewID(id.GetRandomKey())
	ids1 = append(ids1, id1)
	st, s, tp, r, _ := realDhtParams(ids1, "127.0.0.1:"+strconv.Itoa(port))
	address, _ := node.NewAddress("127.0.0.1:" + strconv.Itoa(port-1))
	bootstrapNode := node.NewNode(address)
	dht, err := NewDHT(st, s, tp, r,
		&Options{BootstrapNodes: []*node.Node{bootstrapNode}},
		relay.NewProxy())
	assert.NoError(t, err)
	assert.NotEqual(t, nil, dht)

	go func(dht *DHT) {
		_ = dht.Listen()
		done <- true
	}(dht)

	dht.Disconnect()
	<-done
}

func TestDHT_AuthenticationRequest(t *testing.T) {
	count := 2
	port := 46000
	var dhts []*DHT
	done := make(chan bool)
	var ids []string
	args := []struct {
		first  string
		second error
		name   string
	}{
		{"wrong command", errors.New("unknown command"), "undefined command"},
		{"begin", nil, "begin auth"},
		{"begin", errors.New("authentication unsuccessful"), "begin auth second time"},
		{"revoke", nil, "revoke auth"},
		{"begin", nil, "begin auth"},
	}

	for i := 0; i < count; i++ {
		ids1 := make([]id.ID, 0)
		id1, _ := id.NewID(id.GetRandomKey())
		ids1 = append(ids1, id1)
		ids = append(ids, ids1[0].HashString())
		st, s, tp, r, _ := realDhtParams(ids1, "127.0.0.1:"+strconv.Itoa(port))
		address, _ := node.NewAddress("127.0.0.1:" + strconv.Itoa(port-1))
		bootstrapNode := node.NewNode(address)
		dht, err := NewDHT(st, s, tp, r, &Options{
			BootstrapNodes: []*node.Node{
				bootstrapNode,
			},
		},
			relay.NewProxy())
		port++
		dhts = append(dhts, dht)
		assert.NoError(t, err)
	}

	for _, dht := range dhts {
		ctx, _ := NewContextBuilder(dht).SetDefaultNode().Build()
		assert.Equal(t, 0, dht.NumNodes(ctx))
		go func(dht *DHT) {
			err := dht.Listen()
			assert.Equal(t, closedPacket, err.Error())
			done <- true
		}(dht)
	}

	for _, dht := range dhts {
		dht.Bootstrap()
	}

	ctx, _ := NewContextBuilder(dhts[0]).SetDefaultNode().Build()

	for _, arg := range args {
		t.Run(arg.name, func(t *testing.T) {
			err := dhts[0].AuthenticationRequest(ctx, arg.first, ids[1])
			assert.Equal(t, err, arg.second)
			time.Sleep(time.Millisecond * 200)
		})
	}

	for _, dht := range dhts {
		dht.Disconnect()
	}
	<-done
}

func TestDHT_RelayRequest(t *testing.T) {
	count := 2
	port := 45000
	var dhts []*DHT
	done := make(chan bool)
	var ids []string
	args := []struct {
		first  string
		second error
		name   string
	}{
		{"start", errors.New("unable to execute relay because this node not authenticated"), "start relay w/o auth"},
		{"begin", nil, "begin auth"},
		{"wrong command", errors.New("unknown command"), "undefined command"},
		{"start", nil, "start relay"},
		{"start", errors.New("relay request error"), "start relay twice"},
		{"stop", nil, "stop relay"},
		{"start", nil, "start relay"},
	}

	for i := 0; i < count; i++ {
		ids1 := make([]id.ID, 0)
		id1, _ := id.NewID(id.GetRandomKey())
		ids1 = append(ids1, id1)
		ids = append(ids, ids1[0].HashString())
		st, s, tp, r, _ := realDhtParams(ids1, "127.0.0.1:"+strconv.Itoa(port))
		address, _ := node.NewAddress("127.0.0.1:" + strconv.Itoa(port-1))
		bootstrapNode := node.NewNode(address)
		dht, err := NewDHT(st, s, tp, r, &Options{
			BootstrapNodes: []*node.Node{
				bootstrapNode,
			},
		},
			relay.NewProxy())
		port++
		dhts = append(dhts, dht)
		assert.NoError(t, err)
	}

	for _, dht := range dhts {
		ctx, _ := NewContextBuilder(dht).SetDefaultNode().Build()
		assert.Equal(t, 0, dht.NumNodes(ctx))
		go func(dht *DHT) {
			err := dht.Listen()
			assert.Equal(t, closedPacket, err.Error())
			done <- true
		}(dht)
	}

	for _, dht := range dhts {
		dht.Bootstrap()
	}

	ctx, _ := NewContextBuilder(dhts[0]).SetDefaultNode().Build()

	var err error
	for _, arg := range args {
		t.Run(arg.name, func(t *testing.T) {
			if strings.Contains(arg.name, "begin auth") {
				err = dhts[0].AuthenticationRequest(ctx, arg.first, ids[1])
			} else {

				err = dhts[0].RelayRequest(ctx, arg.first, ids[1])
			}
			assert.Equal(t, err, arg.second)
			time.Sleep(time.Millisecond * 200)
		})
	}

	for _, dht := range dhts {
		dht.Disconnect()
	}
	<-done
}

func TestDHT_ObtainIP(t *testing.T) {
	count := 2
	port := 47000
	var dhts []*DHT
	done := make(chan bool)
	var ids []string

	for i := 0; i < count; i++ {
		ids1 := make([]id.ID, 0)
		id1, _ := id.NewID(id.GetRandomKey())
		ids1 = append(ids1, id1)
		ids = append(ids, ids1[0].HashString())
		st, s, tp, r, _ := realDhtParams(ids1, "127.0.0.1:"+strconv.Itoa(port))
		address, _ := node.NewAddress("127.0.0.1:" + strconv.Itoa(port-1))
		bootstrapNode := node.NewNode(address)
		dht, err := NewDHT(st, s, tp, r, &Options{
			BootstrapNodes: []*node.Node{
				bootstrapNode,
			},
		},
			relay.NewProxy())
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
	}

	for _, dht := range dhts {
		dht.Bootstrap()
	}

	ctx, _ := NewContextBuilder(dhts[0]).SetDefaultNode().Build()

	err := dhts[0].ObtainIP(ctx)
	assert.NoError(t, err)

	for _, dht := range dhts {
		dht.Disconnect()
	}
	<-done
}

func TestDHT_AnalyzeNetwork(t *testing.T) {
	count := 2
	port := 48000
	var dhts []*DHT
	done := make(chan bool)
	var ids []string

	for i := 0; i < count; i++ {
		ids1 := make([]id.ID, 0)
		id1, _ := id.NewID(id.GetRandomKey())
		ids1 = append(ids1, id1)
		ids = append(ids, ids1[0].HashString())
		st, s, tp, r, _ := realDhtParams(ids1, "127.0.0.1:"+strconv.Itoa(port))
		address, _ := node.NewAddress("127.0.0.1:" + strconv.Itoa(port-1))
		bootstrapNode := node.NewNode(address)
		dht, err := NewDHT(st, s, tp, r, &Options{
			BootstrapNodes: []*node.Node{
				bootstrapNode,
			},
		},
			relay.NewProxy())
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
	}

	for _, dht := range dhts {
		dht.Bootstrap()
	}

	ctx, _ := NewContextBuilder(dhts[0]).SetDefaultNode().Build()

	err := dhts[0].ObtainIP(ctx)
	assert.NoError(t, err)

	err = dhts[0].AnalyzeNetwork(ctx)
	assert.NoError(t, err)

	for _, dht := range dhts {
		dht.Disconnect()
	}
	<-done
}
