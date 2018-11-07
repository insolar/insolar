/*
 *    Copyright 2018 Insolar
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

package dhtnetwork

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"strconv"
	"testing"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/network/cascade"
	"github.com/insolar/insolar/network/dhtnetwork/signhandler"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/insolar/insolar/testutils"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/network/dhtnetwork/hosthandler"
	"github.com/insolar/insolar/network/dhtnetwork/routing"
	"github.com/insolar/insolar/network/dhtnetwork/rpc"
	"github.com/insolar/insolar/network/dhtnetwork/store"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/id"
	"github.com/insolar/insolar/network/transport/packet"

	"github.com/insolar/insolar/network/transport/relay"
	"github.com/stretchr/testify/assert"
)

const closedPacket = "closed" // "broken pipe" for kcpTransport

type mockFuture struct {
	result    chan *packet.Packet
	actor     *host.Host
	request   *packet.Packet
	requestID packet.RequestID
}

func (f *mockFuture) ID() packet.RequestID {
	return f.requestID
}

func (f *mockFuture) Actor() *host.Host {
	return f.actor
}

func (f *mockFuture) Request() *packet.Packet {
	return f.request
}

func (f *mockFuture) Result() <-chan *packet.Packet {
	return f.result
}

func (f *mockFuture) GetResult(duration time.Duration) (*packet.Packet, error) {
	result, ok := <-f.Result()
	if !ok || result == nil {
		return nil, errors.New("channel is closed")
	}
	return result, nil
}

func (f *mockFuture) SetResult(msg *packet.Packet) {
	f.result <- msg
}

func (f *mockFuture) Cancel() {}

type mockTransport struct {
	recv          chan *packet.Packet
	send          chan *packet.Packet
	dc            chan bool
	msgChan       chan *packet.Packet
	failNext      bool
	sequence      *uint64
	publicAddress string
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

func (t *mockTransport) PublicAddress() string {
	return t.publicAddress
}

func mockFindHostResponse(request *packet.Packet) *packet.Packet {
	r := &packet.Packet{}
	n := &host.Host{}
	n.ID = request.Sender.ID
	n.Address = request.Sender.Address
	r.Receiver = n
	netAddr, _ := host.NewAddress("0.0.0.0:13001")
	r.Sender = &host.Host{ID: request.Receiver.ID, Address: netAddr}
	r.Type = request.Type
	r.IsResponse = true
	responseData := &packet.ResponseDataFindHost{}
	id1, _ := id.NewID()
	responseData.Closest = []*host.Host{{ID: id1, Address: netAddr}}
	r.Data = responseData
	return r
}

func mockFindHostResponseEmpty(request *packet.Packet) *packet.Packet {
	r := &packet.Packet{}
	n := &host.Host{}
	n.ID = request.Sender.ID
	n.Address = request.Sender.Address
	r.Receiver = n
	netAddr, _ := host.NewAddress("0.0.0.0:14001")
	r.Sender = &host.Host{ID: request.Receiver.ID, Address: netAddr}
	r.Type = request.Type
	r.IsResponse = true
	responseData := &packet.ResponseDataFindHost{}
	responseData.Closest = []*host.Host{}
	r.Data = responseData
	return r
}

func newRealDHT(t *testing.T, bootstrap []*host.Host, port string) *DHT {
	ids1 := make([]id.ID, 0)
	id1, _ := id.NewID()
	ids1 = append(ids1, id1)
	st1, s1, tp1, r1, err := realDhtParams(ids1, "0.0.0.0:"+port)
	assert.NoError(t, err)
	var dht *DHT
	if bootstrap == nil {
		dht, _ = NewDHT(
			st1, s1, tp1, r1,
			&Options{},
			relay.NewProxy(),
			4,
			false,
			testutils.RandomRef(),
			5,
			nil)
	} else {
		dht, _ = NewDHT(
			st1, s1, tp1, r1,
			&Options{
				BootstrapHosts: bootstrap,
			},
			relay.NewProxy(),
			4,
			false,
			testutils.RandomRef(),
			5,
			nil)
	}
	return dht
}

func newDHT(t *testing.T, bootstrap []*host.Host, port string) (*DHT, transport.Transport) {
	zeroID := getIDWithValues()
	ids1 := make([]id.ID, 0)
	id1, _ := id.NewID()
	ids1 = append(ids1, id1)
	st1, s1, tp1, r1, err := dhtParams([]id.ID{zeroID}, "0.0.0.0:"+port)
	assert.NoError(t, err)
	var dht *DHT
	if bootstrap == nil {
		dht, _ = NewDHT(
			st1, s1, tp1, r1,
			&Options{
				RefreshTime:    time.Second * 2,
				ReplicateTime:  time.Second * 2,
				ExpirationTime: time.Second,
			},
			relay.NewProxy(),
			4,
			false,
			testutils.RandomRef(),
			5,
			nil)
	} else {
		dht, _ = NewDHT(
			st1, s1, tp1, r1,
			&Options{
				RefreshTime:    time.Second * 2,
				ReplicateTime:  time.Second * 2,
				ExpirationTime: time.Second,
				BootstrapHosts: bootstrap,
			},
			relay.NewProxy(),
			4,
			false,
			testutils.RandomRef(),
			5,
			nil)
	}
	return dht, tp1
}

func dhtParams(ids []id.ID, address string) (store.Store, *host.Origin, transport.Transport, hosthandler.NetworkCommonFacade, error) {
	st := store.NewMemoryStore()
	addr, _ := host.NewAddress(address)
	origin, err := host.NewOrigin(ids, addr)
	tp := newMockTransport()
	cascade1 := &cascade.Cascade{}
	sign := signhandler.NewSignHandler(nil)
	ncf := hosthandler.NewNetworkCommonFacade(rpc.NewRPCFactory(nil).Create(), cascade1, sign, func(core.Pulse) {})
	return st, origin, tp, ncf, err
}

func realDhtParams(ids []id.ID, address string) (store.Store, *host.Origin, transport.Transport, hosthandler.NetworkCommonFacade, error) {
	st := store.NewMemoryStore()
	addr, _ := host.NewAddress(address)
	origin, _ := host.NewOrigin(ids, addr)
	cfg := configuration.NewConfiguration().Host.Transport
	cfg.Address = address
	cfg.BehindNAT = false
	tp, err := transport.NewTransport(cfg, relay.NewProxy())
	cascade1 := &cascade.Cascade{}
	sign := signhandler.NewSignHandler(nil)
	ncf := hosthandler.NewNetworkCommonFacade(rpc.NewRPCFactory(nil).Create(), cascade1, sign, func(core.Pulse) {})
	return st, origin, tp, ncf, err
}

func realDhtParamsWithId(address string) (store.Store, *host.Origin, transport.Transport, hosthandler.NetworkCommonFacade, error) {
	ids := make([]id.ID, 1)
	ids[0], _ = id.NewID()
	return realDhtParams(ids, address)
}

// Creates twenty DHTs and bootstraps each with the previous
// at the end all should know about each other
func TestBootstrapManyHosts(t *testing.T) {
	port := 15000
	var dhts []*DHT
	hostCount := 10

	for i := 0; i < hostCount; i++ {
		address, err := host.NewAddress("127.0.0.1:" + strconv.Itoa(port-1))
		assert.NoError(t, err)
		bootstrapHost := host.NewHost(address)
		dht := newRealDHT(t, []*host.Host{bootstrapHost}, strconv.Itoa(port))
		port++
		assert.NoError(t, err)
		dhts = append(dhts, dht)
	}

	for _, dht := range dhts {
		ctx, _ := NewContextBuilder(dht).SetDefaultHost().Build()
		assert.Equal(t, 0, dht.NumHosts(ctx))
		go func(dht *DHT) {
			err := dht.Listen()
			assert.Equal(t, closedPacket, err.Error())
		}(dht)
	}

	for _, dht := range dhts {
		err := dht.Bootstrap()
		assert.NoError(t, err)
	}

	for _, dht := range dhts {
		assert.Equal(t, hostCount-1, dht.NumHosts(GetDefaultCtx(dht)))
		dht.Disconnect()
	}
}

// Creates two DHTs and bootstraps using only IP:Port. Connecting host should
// ping the first host to find its RequestID
func TestBootstrapNoID(t *testing.T) {
	done := make(chan bool)

	dht1 := newRealDHT(t, nil, "18000")
	dht2 := newRealDHT(t, []*host.Host{{ID: dht1.origin.IDs[0], Address: dht1.origin.Address}}, "18001")

	assert.Equal(t, 0, dht1.NumHosts(GetDefaultCtx(dht1)))
	assert.Equal(t, 0, dht2.NumHosts(GetDefaultCtx(dht2)))

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

	err := dht1.Listen()
	assert.Equal(t, closedPacket, err.Error())

	assert.Equal(t, 1, dht1.NumHosts(GetDefaultCtx(dht1)))
	assert.Equal(t, 1, dht2.NumHosts(GetDefaultCtx(dht2)))

	<-done
	<-done
}

// create two DHTs have them connect and bootstrap, then disconnect. Repeat
// 100 times to ensure that we can use the same IP and port without EADDRINUSE
// errors.
func TestReconnect(t *testing.T) {
	for i := 0; i < 100; i++ {
		done := make(chan bool)

		dht1 := newRealDHT(t, nil, "19000")
		dht2 := newRealDHT(t, []*host.Host{{Address: dht1.origin.Address}}, "19001")

		assert.Equal(t, 0, dht1.NumHosts(GetDefaultCtx(dht1)))

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

		err := dht1.Listen()
		assert.Equal(t, closedPacket, err.Error())

		assert.Equal(t, 1, dht1.NumHosts(GetDefaultCtx(dht1)))
		assert.Equal(t, 1, dht2.NumHosts(GetDefaultCtx(dht2)))

		<-done
		<-done
	}
}

// create two DHTs and have them connect. Send a store packet with 100mb
// payload from one host to another. Ensure that the other host now has
// this data in its store.
func TestStoreAndFindLargeValue(t *testing.T) {
	// this test is skipped cuz store data execution time is undefined.
	t.Skip()
	done := make(chan bool)

	dht1 := newRealDHT(t, nil, "20000")
	dht2 := newRealDHT(t, []*host.Host{{Address: dht1.origin.Address}}, "20001")

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

	dht2.Bootstrap()

	payload := make([]byte, 1000000)

	storeKey, err := dht1.StoreData(GetDefaultCtx(dht1), payload[:])
	assert.NoError(t, err)

	value, exists, err := dht2.Get(GetDefaultCtx(dht1), storeKey)
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
	done := make(chan int)

	bootstrapAddr, _ := host.NewAddress("0.0.0.0:21001")
	dht, _ := newDHT(t, []*host.Host{{Address: bootstrapAddr}}, "21000")
	mockTp := dht.transport.(*mockTransport)

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

// Tests sending a packet which results in a successful send, but the host
// never responds
func TestHostResponseSendError(t *testing.T) {
	done := make(chan int)

	bootstrapAddr, _ := host.NewAddress("0.0.0.0:22001")

	zeroID := getIDWithValues()
	st, s, tp, r, err := dhtParams([]id.ID{zeroID}, "0.0.0.0:22000")
	assert.NoError(t, err)

	dht, _ := NewDHT(st, s, tp, r, &Options{
		BootstrapHosts: []*host.Host{{
			ID:      getZerodIDWithNthByte(1, byte(255)),
			Address: bootstrapAddr,
		}},
	},
		relay.NewProxy(), 4, false, testutils.RandomRef(), 5, nil)

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
				res := mockFindHostResponse(request)
				mockTp.send <- res
			}
		}
	}()

	dht.Bootstrap()
	assert.Equal(t, 1, dht.tables[0].TotalHosts())
	dht.Disconnect()

	<-done
}

// Tests a bucket refresh by setting a very low RefreshTime value, adding a single
// host to a bucket, and waiting for the refresh packet for the bucket
func TestBucketRefresh(t *testing.T) {
	done := make(chan int)
	refresh := make(chan int)

	bootstrapAddr, _ := host.NewAddress("0.0.0.0:23001")
	dht, tp := newDHT(
		t,
		[]*host.Host{
			{
				Address: bootstrapAddr,
				ID:      getZerodIDWithNthByte(1, byte(255)),
			},
		},
		"23000")

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

			res := mockFindHostResponseEmpty(request)
			mockTp.send <- res

			if queries == 2 {
				close(refresh)
			}
		}
	}()

	dht.Bootstrap()

	assert.Equal(t, 1, dht.tables[0].TotalHosts())

	<-refresh

	dht.Disconnect()

	<-done
}

// Tets store replication by setting the ReplicateTime time to a very small value.
// Stores some data, and then expects another store packet in ReplicateTime time
func TestStoreReplication(t *testing.T) {
	done := make(chan int)
	replicate := make(chan int)

	bootstrapAddr, _ := host.NewAddress("0.0.0.0:24001")
	// TODO: try to do this
	// dht, tp := newDHT(t, []*host.Host{
	// 	{
	// 		ID:      getZerodIDWithNthByte(1, byte(255)),
	// 		Address: bootstrapAddr,
	// 	},
	// },
	// 	"24000",
	// )

	zeroID := getIDWithValues()
	st, s, tp, r, err := dhtParams([]id.ID{zeroID}, "0.0.0.0:24000")
	assert.NoError(t, err)

	dht, _ := NewDHT(st, s, tp, r, &Options{
		ReplicateTime: time.Second * 2,
		BootstrapHosts: []*host.Host{{
			ID:      getZerodIDWithNthByte(1, byte(255)),
			Address: bootstrapAddr,
		}},
	},
		relay.NewProxy(), 4, false, testutils.RandomRef(), 5, nil)

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
			case types.TypeFindHost:
				res := mockFindHostResponseEmpty(request)
				mockTp.send <- res
			case types.TypeStore:
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

	dht.StoreData(GetDefaultCtx(dht), []byte("foo"))

	<-replicate

	dht.Disconnect()

	<-done
}

// Test Expiration by setting ExpirationTime to a very low value. Store a value,
// and then wait longer than ExpirationTime. The value should no longer exist in
// the store.
func TestStoreExpiration(t *testing.T) {
	done := make(chan bool)

	dht, _ := newDHT(t, nil, "25000")

	go func() {
		dht.Listen()
		done <- true
	}()

	k, _ := dht.StoreData(GetDefaultCtx(dht), []byte("foo"))

	v, exists, _ := dht.Get(GetDefaultCtx(dht), k)
	assert.Equal(t, true, exists)

	assert.Equal(t, []byte("foo"), v)

	<-time.After(time.Second * 3)

	_, exists, _ = dht.Get(GetDefaultCtx(dht), k)

	assert.Equal(t, false, exists)

	dht.Disconnect()
	<-done
}

// create a new host and bootstrap it. All hosts in the network know of a
// single host closer to the original host. This continues until every MaxContactsInBucket bucket
// is occupied.
func TestFindHostAllBuckets(t *testing.T) {
	t.Skip()
	done := make(chan bool)

	bootstrapAddr, _ := host.NewAddress("127.0.0.1:26011")
	dht := newRealDHT(t, []*host.Host{{Address: bootstrapAddr}}, "26010")
	mockTp := dht.transport.(*mockTransport)

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

			res := mockFindHostResponse(request)

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
		assert.Equal(t, 0, len(v))
	}

	dht.Disconnect()
	<-done
}

func TestGetRandomIDFromBucket(t *testing.T) {
	done := make(chan bool)

	dht := newRealDHT(t, nil, "28000")

	go func() {
		dht.Listen()
		done <- true
	}()

	// Bytes should be equal up to the bucket index that the random RequestID was
	// generated for, and random afterwards
	for i := 0; i < routing.KeyBitSize/8; i++ {
		// r := dht.tables[0].GetRandomIDFromBucket(i * 8)
		for j := 0; j < i; j++ {
			// assert.Equal(t, byte(0), r[j])
		}
	}

	dht.Disconnect()
	<-done
}

func getZerodIDWithNthByte(n int, v byte) id.ID {
	id1 := getIDWithValues()
	id1.Bytes()[n] = v
	return id1
}

func getIDWithValues() id.ID {
	id1, _ := id.NewID()
	return id1
}

func TestDHT_Listen(t *testing.T) {
	count := 5
	port := 8000
	var dhts []*DHT
	done := make(chan bool)

	for i := 0; i < count; i++ {
		address, _ := host.NewAddress("127.0.0.1:" + strconv.Itoa(port-1))
		bootstrapHost := host.NewHost(address)
		dht := newRealDHT(t, []*host.Host{{Address: bootstrapHost.Address}}, strconv.Itoa(port))
		port++
		dhts = append(dhts, dht)
	}

	for _, dht := range dhts {
		ctx, _ := NewContextBuilder(dht).SetDefaultHost().Build()
		assert.Equal(t, 0, dht.NumHosts(ctx))
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
		address, _ := host.NewAddress("127.0.0.1:" + strconv.Itoa(port-1))
		bootstrapHost := host.NewHost(address)
		dht := newRealDHT(t, []*host.Host{{Address: bootstrapHost.Address}}, strconv.Itoa(port))
		port++
		dhts = append(dhts, dht)
	}

	for _, dht := range dhts {
		ctx, _ := NewContextBuilder(dht).SetDefaultHost().Build()
		assert.Equal(t, 0, dht.NumHosts(ctx))
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

func TestNewDHT(t *testing.T) {
	done := make(chan bool)
	port := 11000
	address, _ := host.NewAddress("127.0.0.1:" + strconv.Itoa(port-1))
	bootstrapHost := host.NewHost(address)
	dht := newRealDHT(t, []*host.Host{{Address: bootstrapHost.Address}}, strconv.Itoa(port))
	assert.NotEqual(t, nil, dht)

	go func(dht *DHT) {
		_ = dht.Listen()
		done <- true
	}(dht)

	dht.Disconnect()
	<-done
}

func TestDHT_AnalyzeNetwork(t *testing.T) {
	count := 2
	port := 48000
	var dhts []*DHT
	done := make(chan bool)

	for i := 0; i < count; i++ {
		address, _ := host.NewAddress("127.0.0.1:" + strconv.Itoa(port-1))
		bootstrapHost := host.NewHost(address)
		dht := newRealDHT(t, []*host.Host{{Address: bootstrapHost.Address}}, strconv.Itoa(port))
		port++
		dhts = append(dhts, dht)
	}

	for _, dht := range dhts {
		ctx, _ := NewContextBuilder(dht).SetDefaultHost().Build()
		assert.Equal(t, 0, dht.NumHosts(ctx))
		go func(dht *DHT) {
			err := dht.Listen()
			assert.Equal(t, "closed", err.Error())
			done <- true
		}(dht)
	}

	go func() {
		for _, dht := range dhts {
			dht.Bootstrap()
		}
	}()

	ctx, _ := NewContextBuilder(dhts[0]).SetDefaultHost().Build()

	err := dhts[0].ObtainIP()
	assert.NoError(t, err)

	err = dhts[0].AnalyzeNetwork(ctx)
	assert.NoError(t, err)

	for _, dht := range dhts {
		dht.Disconnect()
	}
	<-done
}

func TestDHT_StartCheckNodesRole(t *testing.T) {
	var dhts []*DHT

	done := make(chan bool)

	ids1 := make([]id.ID, 0)
	id1, _ := id.NewID()
	ids1 = append(ids1, id1)
	st, s, tp, r, err := realDhtParams(ids1, "127.0.0.1:16000")
	dht1, _ := NewDHT(st, s, tp, r, &Options{}, relay.NewProxy(), 4, false, testutils.RandomRef(), 5, nil)
	assert.NoError(t, err)

	bootstrapAddr2, _ := host.NewAddress("127.0.0.1:16000")
	st2, s2, tp2, r2, _ := realDhtParams(nil, "127.0.0.1:16001")
	dht2, _ := NewDHT(st2, s2, tp2, r2, &Options{
		BootstrapHosts: []*host.Host{
			{
				ID:      ids1[0],
				Address: bootstrapAddr2,
			},
		},
	},
		relay.NewProxy(), 4, false, testutils.RandomRef(), 5, nil)

	dhts = append(dhts, dht1)
	dhts = append(dhts, dht2)

	for _, dht := range dhts {
		ctx, _ := NewContextBuilder(dht).SetDefaultHost().Build()
		assert.Equal(t, 0, dht.NumHosts(ctx))
		go func(dht *DHT) {
			err := dht.Listen()
			assert.Equal(t, "closed", err.Error())
			done <- true
		}(dht)
	}

	for _, dht := range dhts {
		dht.Bootstrap()
	}

	err = dhts[1].CheckNodeRole("domain ID")
	assert.NoError(t, err)

	for _, dht := range dhts {
		dht.Disconnect()
	}

	<-done
}

func TestDHT_RemoteProcedureCall(t *testing.T) {
	bootstrapAddr, _ := host.NewAddress("127.0.0.1:23220")
	dht1 := newRealDHT(t, nil, "23220")
	dht2 := newRealDHT(t, []*host.Host{{Address: bootstrapAddr}}, "23221")

	go func(dht *DHT) {
		dht1.Listen()
	}(dht1)
	go func(dht *DHT) {
		dht2.Listen()
	}(dht2)
	dht1.Bootstrap()
	dht2.Bootstrap()

	msg := &message.CallMethod{
		ObjectRef: core.NewRefFromBase58("test"),
		Method:    "test",
		Arguments: []byte("test"),
	}

	reqBuff, _ := message.Serialize(msg)
	msg1, _ := ioutil.ReadAll(reqBuff)

	key1, _ := ecdsa.GeneratePrivateKey()
	key2, _ := ecdsa.GeneratePrivateKey()

	keeper1 := nodenetwork.NewNodeKeeper(nodenetwork.NewNode(dht1.nodeID, nil, nil, 0, 0, "", ""))
	keeper2 := nodenetwork.NewNodeKeeper(nodenetwork.NewNode(dht1.nodeID, nil, nil, 0, 0, "", ""))

	keeper1.AddActiveNodes([]core.Node{
		nodenetwork.NewNode(
			dht2.nodeID,
			[]core.NodeRole{core.RoleUnknown},
			&key2.PublicKey,
			5,
			2,
			"address",
			"",
		),
	})
	keeper2.AddActiveNodes([]core.Node{
		nodenetwork.NewNode(
			dht1.nodeID,
			[]core.NodeRole{core.RoleUnknown},
			&key1.PublicKey,
			5,
			2,
			"address",
			"",
		),
	})

	dht1.SetNodeKeeper(keeper1)
	dht2.SetNodeKeeper(keeper2)

	dht1.RemoteProcedureRegister("test", func(args [][]byte) ([]byte, error) {
		return nil, nil
	})

	dht2.RemoteProcedureCall(GetDefaultCtx(dht1), dht1.GetOriginHost().IDs[0].String(), "test", [][]byte{msg1})
}

func TestDHT_MessageSign(t *testing.T) {
	key, _ := ecdsa.GeneratePrivateKey()
	key2, _ := ecdsa.GeneratePrivateKey()
	ref := testutils.RandomRef()

	tmp := core.Message(&message.GenesisRequest{})
	msg, err := message.NewSignedMessage(context.TODO(), tmp, ref, key, 0)
	assert.NoError(t, err)
	assert.True(t, msg.IsValid(&key.PublicKey))
	assert.False(t, msg.IsValid(&key2.PublicKey))
}

func TestDHT_Getters(t *testing.T) {
	dht1 := newRealDHT(t, nil, "0")
	outerHostCount := 3

	relay1 := "127.0.0.1:123123"
	proxy1 := "127.0.0.1:123124"
	hostAddr, _ := host.NewAddress("127.0.0.1:50001")
	str1 := "4gU79K6woTZDvn4YUFHauNKfcHW69X42uyk8ZvRevCiMv3PLS24eM1vcA9mhKPv8b2jWj9J5RgGN9CB7PUzCtBsj"
	ref1 := core.NewRefFromBase58(str1)

	host1 := host.NewHost(hostAddr)
	//newKey, _ := ecdsa.ImportPrivateKey("asd")

	assert.False(t, dht1.HostIsAuthenticated(ref1.String()))
	_, check := dht1.KeyIsReceived(ref1.String())
	assert.False(t, check)

	dht1.AddPossibleRelayID(relay1)
	dht1.AddPossibleProxyID(proxy1)
	dht1.AddRelayClient(host1)
	dht1.AddReceivedKey(ref1.String(), host1.ID.Bytes())
	dht1.AddAuthSentKey(ref1.String(), host1.ID.Bytes())
	dht1.AddSubnetID("id", ref1.String())
	dht1.AddSubnetID("id2", host1.ID.String())
	dht1.AddProxyHost(host1.ID.String())
	dht1.SetAuthStatus(ref1.String(), true)
	dht1.SetOuterHostsCount(outerHostCount)
	dht1.SetHighKnownHostID(ref1.String())

	assert.True(t, dht1.HostIsAuthenticated(ref1.String()))
	_, check = dht1.KeyIsReceived(ref1.String())
	assert.True(t, check)
	assert.Equal(t, dht1.GetHighKnownHostID(), ref1.String())
	assert.Equal(t, dht1.GetSelfKnownOuterHosts(), 0)
	assert.Equal(t, dht1.GetOuterHostsCount(), outerHostCount)
	assert.Equal(t, dht1.GetProxyHostsCount(), 1)
	assert.True(t, dht1.EqualAuthSentKey(ref1.String(), host1.ID.Bytes()))

	dht1.RemoveRelayClient(host1)
	dht1.RemoveAuthSentKeys(ref1.String())
	dht1.RemovePossibleProxyID(ref1.String())
	dht1.RemoveProxyHost(ref1.String())
	dht1.RemoveAuthHost(ref1.String())
}

func TestDHT_GetHostsFromBootstrap(t *testing.T) {
	// prefix := "127.0.0.1:"
	port := 10000
	bootstrapAddresses := make([]string, 0)
	dhts := make([]*DHT, 0)
	b := 2
	n := 5

	for i := 0; i < b; i++ {
		dht := newRealDHT(t, nil, strconv.Itoa(port))
		bootstrapAddresses = append(bootstrapAddresses, "127.0.0.1:"+strconv.Itoa(port))
		dhts = append(dhts, dht)
		go dht.Listen()
		dht.Bootstrap()
		port++
	}

	bootstrapHosts := make([]*host.Host, len(bootstrapAddresses))
	for i, h := range bootstrapAddresses {
		address, _ := host.NewAddress(h)
		bootstrapHosts[i] = host.NewHost(address)
	}

	for i := 0; i < n; i++ {
		dht := newRealDHT(t, bootstrapHosts, strconv.Itoa(port))
		dhts = append(dhts, dht)
		go dht.Listen()
		dht.Bootstrap()
		dht.GetHostsFromBootstrap()
		port++
	}
	lastDht := dhts[len(dhts)-1]
	hostsCount := lastDht.HtFromCtx(GetDefaultCtx(lastDht)).TotalHosts()
	assert.Equal(t, b+n-1, hostsCount)

	for _, dht := range dhts {
		dht.Disconnect()
	}
}

func TestDHT_BootstrapInfinity(t *testing.T) {
	bootstrapAddress := "127.0.0.1:10000"
	bootstrapDht := newRealDHT(t, nil, "10000")

	go func() {
		time.Sleep(time.Second * 5)
		bootstrapDht.Bootstrap()
		bootstrapDht.Listen()
	}()

	bootstrapHosts := make([]*host.Host, 1)
	a, _ := host.NewAddress(bootstrapAddress)
	bootstrapHosts[0] = host.NewHost(a)
	dht := newRealDHT(t, bootstrapHosts, "10001")

	defer func() {
		dht.Disconnect()
		bootstrapDht.Disconnect()
	}()

	go dht.Listen()
	err := dht.Bootstrap()
	assert.NoError(t, err)
}
