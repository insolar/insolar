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

package hostnetwork

import (
	"bytes"
	"context"
	"errors"
	"log"
	"math"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/huandu/xstrings"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/id"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/relay"
	"github.com/insolar/insolar/network/hostnetwork/routing"
	"github.com/insolar/insolar/network/hostnetwork/rpc"
	"github.com/insolar/insolar/network/hostnetwork/store"
	"github.com/insolar/insolar/network/hostnetwork/transport"
	"github.com/jbenet/go-base58"
)

// RPC is remote procedure call interface
type RPC interface {
	RemoteProcedureCall(ctx Context, target string, method string, args [][]byte) (result []byte, err error)
	RemoteProcedureRegister(name string, method core.RemoteProcedure)
}

// DHT represents the state of the local host in the distributed hash table.
type DHT struct {
	tables  []*routing.HashTable
	options *Options

	origin *host.Origin

	transport transport.Transport
	store     store.Store
	rpc       rpc.RPC
	relay     relay.Relay
	proxy     relay.Proxy
	auth      AuthInfo
	subnet    Subnet
}

// AuthInfo collects some information about authentication.
type AuthInfo struct {
	// Sent/received unique auth keys.
	SentKeys     map[string][]byte
	ReceivedKeys map[string][]byte

	authenticatedHosts map[string]bool

	mut sync.Mutex
}

// Subnet collects some information about self network part
type Subnet struct {
	SubnetIDs        map[string][]string // key - ip, value - id
	HomeSubnetKey    string              // key of home subnet fo SubnetIDs
	PossibleRelayIDs []string
	PossibleProxyIDs []string
	HighKnownHosts   HighKnownOuterHostsHost
}

// HighKnownOuterHostsHost collects an information about host in home subnet which have a more known outer hosts.
type HighKnownOuterHostsHost struct {
	ID                  string
	OuterHosts          int // high known outer hosts by ID host
	SelfKnownOuterHosts int
}

// Options contains configuration options for the local host.
type Options struct {
	// The hosts being used to bootstrap the network. Without a bootstrap
	// host there is no way to connect to the network. NetworkHosts can be
	// initialized via host.NewHost().
	BootstrapHosts []*host.Host

	// The time after which a key/value pair expires;
	// this is a time-to-live (TTL) from the original publication date.
	ExpirationTime time.Duration

	// Seconds after which an otherwise unaccessed bucket must be refreshed.
	RefreshTime time.Duration

	// The interval between Kademlia replication events, when a host is
	// required to publish its entire database.
	ReplicateTime time.Duration

	// The time after which the original publisher must
	// republish a key/value pair. Currently not implemented.
	RepublishTime time.Duration

	// The maximum time to wait for a response from a host before discarding
	// it from the bucket.
	PingTimeout time.Duration

	// The maximum time to wait for a response to any packet.
	PacketTimeout time.Duration
}

// NewDHT initializes a new DHT host.
func NewDHT(store store.Store, origin *host.Origin, transport transport.Transport, rpc rpc.RPC, options *Options, proxy relay.Proxy) (dht *DHT, err error) {
	tables, err := newTables(origin)
	if err != nil {
		return nil, err
	}

	rel := relay.NewRelay()

	dht = &DHT{
		options:   options,
		origin:    origin,
		rpc:       rpc,
		transport: transport,
		tables:    tables,
		store:     store,
		relay:     rel,
		proxy:     proxy,
	}

	if options.ExpirationTime == 0 {
		options.ExpirationTime = time.Second * 86410
	}

	if options.RefreshTime == 0 {
		options.RefreshTime = time.Second * 3600
	}

	if options.ReplicateTime == 0 {
		options.ReplicateTime = time.Second * 3600
	}

	if options.RepublishTime == 0 {
		options.RepublishTime = time.Second * 86400
	}

	if options.PingTimeout == 0 {
		options.PingTimeout = time.Second * 1
	}

	if options.PacketTimeout == 0 {
		options.PacketTimeout = time.Second * 10
	}

	dht.auth.authenticatedHosts = make(map[string]bool)
	dht.auth.SentKeys = make(map[string][]byte)
	dht.auth.ReceivedKeys = make(map[string][]byte)

	dht.subnet.SubnetIDs = make(map[string][]string)

	return dht, nil
}

func newTables(origin *host.Origin) ([]*routing.HashTable, error) {
	tables := make([]*routing.HashTable, len(origin.IDs))

	for i, id1 := range origin.IDs {
		ht, err := routing.NewHashTable(id1, origin.Address)
		if err != nil {
			return nil, err
		}

		tables[i] = ht
	}

	return tables, nil
}

func (dht *DHT) getExpirationTime(ctx context.Context, key []byte) time.Time {
	ht := dht.htFromCtx(ctx)

	bucket := routing.GetBucketIndexFromDifferingBit(key, ht.Origin.ID.GetKey())
	var total int
	for i := 0; i < bucket; i++ {
		total += ht.GetTotalHostsInBucket(i)
	}
	closer := ht.GetAllHostsInBucketCloserThan(bucket, key)
	score := total + len(closer)

	if score == 0 {
		score = 1
	}

	if score > routing.MaxContactsInBucket {
		return time.Now().Add(dht.options.ExpirationTime)
	}

	day := dht.options.ExpirationTime
	seconds := day.Nanoseconds() * int64(math.Exp(float64(routing.MaxContactsInBucket/score)))
	dur := time.Second * time.Duration(seconds)
	return time.Now().Add(dur)
}

// Store stores data on the network. This will trigger an iterateStore loop.
// The base58 encoded identifier will be returned if the store is successful.
func (dht *DHT) Store(ctx Context, data []byte) (id string, err error) {
	key := store.NewKey(data)
	expiration := dht.getExpirationTime(ctx, key)
	replication := time.Now().Add(dht.options.ReplicateTime)
	err = dht.store.Store(key, data, replication, expiration, true)
	if err != nil {
		return "", err
	}
	_, _, err = dht.iterate(ctx, routing.IterateStore, key, data)
	if err != nil {
		return "", err
	}
	str := base58.Encode(key)
	return str, nil
}

// Get retrieves data from the transport using key. Key is the base58 encoded
// identifier of the data.
func (dht *DHT) Get(ctx Context, key string) ([]byte, bool, error) {
	keyBytes := base58.Decode(key)
	if len(keyBytes) != routing.MaxContactsInBucket {
		return nil, false, errors.New("invalid key")
	}

	value, exists := dht.store.Retrieve(keyBytes)
	if !exists {
		var err error
		value, _, err = dht.iterate(ctx, routing.IterateFindValue, keyBytes, nil)
		if err != nil {
			return nil, false, err
		}
		if value != nil {
			exists = true
		}
	}

	return value, exists, nil
}

// FindHost returns target host's real network address.
func (dht *DHT) FindHost(ctx Context, key string) (*host.Host, bool, error) {
	keyBytes := base58.Decode(key)
	if len(keyBytes) != routing.MaxContactsInBucket {
		return nil, false, errors.New("invalid key")
	}
	ht := dht.htFromCtx(ctx)

	if ht.Origin.ID.KeyEqual(keyBytes) {
		return ht.Origin, true, nil
	}

	var targetHost *host.Host
	var exists = false
	routeSet := ht.GetClosestContacts(1, keyBytes, nil)

	if routeSet.Len() > 0 && routeSet.FirstHost().ID.KeyEqual(keyBytes) {
		targetHost = routeSet.FirstHost()
		exists = true
	} else if dht.proxy.ProxyHostsCount() > 0 {
		address, _ := host.NewAddress(dht.proxy.GetNextProxyAddress())
		// TODO: current key insertion
		id1, _ := id.NewID()
		targetHost = &host.Host{ID: id1, Address: address}
		return targetHost, true, nil
	} else {
		log.Println("Host not found in routing table. Iterating through network...")
		_, closest, err := dht.iterate(ctx, routing.IterateFindHost, keyBytes, nil)
		if err != nil {
			return nil, false, err
		}
		for i := range closest {
			if closest[i].ID.KeyEqual(keyBytes) {
				targetHost = closest[i]
				exists = true
			}
		}
	}

	return targetHost, exists, nil
}

// NumHosts returns the total number of hosts stored in the local routing table.
func (dht *DHT) NumHosts(ctx Context) int {
	ht := dht.htFromCtx(ctx)
	return ht.TotalHosts()
}

// GetOriginHost returns the local host.
func (dht *DHT) GetOriginHost(ctx Context) *host.Host {
	ht := dht.htFromCtx(ctx)
	return ht.Origin
}

// Listen begins listening on the socket for incoming Packets.
func (dht *DHT) Listen() error {
	start := make(chan bool)
	stop := make(chan bool)

	go dht.handleDisconnect(start, stop)
	go dht.handlePackets(start, stop)
	go dht.handleStoreTimers(start, stop)

	return dht.transport.Start()
}

// Bootstrap attempts to bootstrap the network using the BootstrapHosts provided
// to the Options struct. This will trigger an iterateBootstrap to the provided
// BootstrapHosts.
func (dht *DHT) Bootstrap() error {
	if len(dht.options.BootstrapHosts) == 0 {
		return nil
	}
	var futures []transport.Future
	wg := &sync.WaitGroup{}
	cb := NewContextBuilder(dht)

	for _, ht := range dht.tables {
		futures = dht.iterateBootstrapHosts(ht, cb, wg, futures)
	}

	for _, f := range futures {
		go func(future transport.Future) {
			select {
			case result := <-future.Result():
				// If result is nil, channel was closed
				if result != nil {
					ctx, err := cb.SetHostByID(result.Receiver.ID).Build()
					// TODO: must return error here
					if err != nil {
						log.Fatal(err)
					}
					dht.addHost(ctx, routing.NewRouteHost(result.Sender))
				}
				wg.Done()
				return
			case <-time.After(dht.options.PacketTimeout):
				log.Println("bootstrap response timeout")
				future.Cancel()
				wg.Done()
				return
			}
		}(f)
	}

	wg.Wait()
	return dht.iterateHt(cb)
}

func (dht *DHT) iterateHt(cb ContextBuilder) error {
	for _, ht := range dht.tables {
		ctx, err := cb.SetHostByID(ht.Origin.ID).Build()
		if err != nil {
			return err
		}

		if dht.NumHosts(ctx) > 0 {
			_, _, err = dht.iterate(ctx, routing.IterateBootstrap, ht.Origin.ID.GetKey(), nil)
			return err
		}
	}
	return nil
}

func (dht *DHT) iterateBootstrapHosts(
	ht *routing.HashTable,
	cb ContextBuilder,
	wg *sync.WaitGroup,
	futures []transport.Future,
) []transport.Future {
	ctx, err := cb.SetHostByID(ht.Origin.ID).Build()
	if err != nil {
		return futures
	}
	for _, bn := range dht.options.BootstrapHosts {
		request := packet.NewPingPacket(ht.Origin, bn)

		if bn.ID.GetKey() == nil {
			res, err := dht.transport.SendRequest(request)
			if err != nil {
				continue
			}
			log.Println("sending ping to " + bn.Address.String())
			wg.Add(1)
			futures = append(futures, res)
		} else {
			routeHost := routing.NewRouteHost(bn)
			dht.addHost(ctx, routeHost)
		}
	}
	return futures
}

// Disconnect will trigger a Stop from the network.
func (dht *DHT) Disconnect() {
	dht.transport.Stop()
}

// Iterate does an iterative search through the network. This can be done
// for multiple reasons. These reasons include:
//     iterateStore - Used to store new information in the network.
//     iterateFindHost - Used to find host in the network given host abstract address.
//     iterateFindValue - Used to find a value among the network given a key.
//     iterateBootstrap - Used to bootstrap the network.
func (dht *DHT) iterate(ctx Context, t routing.IterateType, target []byte, data []byte) (value []byte, closest []*host.Host, err error) {
	ht := dht.htFromCtx(ctx)
	routeSet := ht.GetClosestContacts(routing.ParallelCalls, target, []*host.Host{})

	// We keep track of hosts contacted so far. We don't contact the same host
	// twice.
	var contacted = make(map[string]bool)

	// According to the Kademlia white paper, after a round of FIND_NODE RPCs
	// fails to provide a host closer than closestHost, we should send a
	// FIND_NODE RPC to all remaining hosts in the route set that have not
	// yet been contacted.
	queryRest := false

	// We keep a reference to the closestHost. If after performing a search
	// we do not find a closer host, we stop searching.
	if routeSet.Len() == 0 {
		return nil, nil, nil
	}

	closestHost := routeSet.FirstHost()

	checkAndRefreshTimeForBucket(t, ht, target)

	var removeFromRouteSet []*host.Host

	for {
		var futures []transport.Future
		var futuresCount int

		futures, removeFromRouteSet = dht.sendPacketToAlphaHosts(routeSet, queryRest, t, ht, contacted, target, futures, removeFromRouteSet)

		routeSet.RemoveMany(routing.RouteHostsFrom(removeFromRouteSet))

		futuresCount = len(futures)

		resultChan := make(chan *packet.Packet)
		dht.setUpResultChan(futures, ctx, resultChan)

		value, closest, err = dht.checkFuturesCountAndGo(t, &queryRest, routeSet, futuresCount, resultChan, target, closest)
		if (err == nil) || ((err != nil) && (err.Error() != "do nothing")) {
			return value, closest, err
		}

		sort.Sort(routeSet)

		var tmpValue []byte
		var tmpClosest []*host.Host
		var tmpHost *host.Host
		tmpValue, tmpClosest, tmpHost, err = dht.iterateIsDone(t, &queryRest, routeSet, data, ht, closestHost)
		if err == nil {
			return tmpValue, tmpClosest, err
		} else if tmpHost != nil {
			closestHost = tmpHost
		}
	}
}

func (dht *DHT) iterateIsDone(
	t routing.IterateType,
	queryRest *bool,
	routeSet *routing.RouteSet,
	data []byte,
	ht *routing.HashTable,
	closestHost *host.Host,
) (value []byte, closest []*host.Host, close *host.Host, err error) {

	if routeSet.FirstHost().ID.KeyEqual(closestHost.ID.GetKey()) || *(queryRest) {
		switch t {
		case routing.IterateBootstrap:
			if !(*queryRest) {
				*queryRest = true
				err = errors.New("do nothing")
				return nil, nil, nil, err
			}
			return nil, routeSet.Hosts(), nil, nil
		case routing.IterateFindHost, routing.IterateFindValue:
			return nil, routeSet.Hosts(), nil, nil
		case routing.IterateStore:
			for i, receiver := range routeSet.Hosts() {
				if i >= routing.MaxContactsInBucket {
					return nil, nil, nil, nil
				}

				msg := packet.NewBuilder().Sender(ht.Origin).Receiver(receiver).Type(packet.TypeStore).Request(
					&packet.RequestDataStore{
						Data: data,
					}).Build()

				future, _ := dht.transport.SendRequest(msg)
				// We do not need to handle result of this packet
				future.Cancel()
			}
			return nil, nil, nil, nil
		}
	} else {
		err = errors.New("do nothing")
		return nil, nil, routeSet.FirstHost(), err
	}
	err = errors.New("do nothing")
	return nil, nil, nil, err
}

func (dht *DHT) checkFuturesCountAndGo(
	t routing.IterateType,
	queryRest *bool,
	routeSet *routing.RouteSet,
	futuresCount int,
	resultChan chan *packet.Packet,
	target []byte,
	close []*host.Host,
) ([]byte, []*host.Host, error) {

	var err error
	var results []*packet.Packet
	var selected bool
	if futuresCount > 0 {
	Loop:
		for {
			results, selected = dht.selectResultChan(resultChan, &futuresCount, results)
			if selected {
				break Loop
			}
		}

		_, close, err = resultsIterate(t, results, routeSet, target)
		if close != nil {
			return nil, close, err
		}
	}

	if !*queryRest && routeSet.Len() == 0 {
		return nil, close, nil
	}
	err = errors.New("do nothing")
	return nil, close, err
}

func resultsIterate(
	t routing.IterateType,
	results []*packet.Packet,
	routeSet *routing.RouteSet,
	target []byte,
) (value []byte, closest []*host.Host, err error) {

	for _, result := range results {
		if result.Error != nil {
			routeSet.Remove(routing.NewRouteHost(result.Sender))
			continue
		}
		switch t {
		case routing.IterateBootstrap, routing.IterateFindHost, routing.IterateStore:
			responseData := result.Data.(*packet.ResponseDataFindHost)
			if len(responseData.Closest) > 0 && responseData.Closest[0].ID.KeyEqual(target) {
				return nil, responseData.Closest, nil
			}
			routeSet.AppendMany(routing.RouteHostsFrom(responseData.Closest))
		case routing.IterateFindValue:
			responseData := result.Data.(*packet.ResponseDataFindValue)
			routeSet.AppendMany(routing.RouteHostsFrom(responseData.Closest))
			if responseData.Value != nil {
				// TODO When an iterateFindValue succeeds, the initiator must
				// store the key/value pair at the closest receiver seen which did
				// not return the value.
				return responseData.Value, nil, nil
			}
		}
	}
	return nil, nil, nil
}

func checkAndRefreshTimeForBucket(t routing.IterateType, ht *routing.HashTable, target []byte) {
	if t == routing.IterateBootstrap {
		bucket := routing.GetBucketIndexFromDifferingBit(target, ht.Origin.ID.GetKey())
		ht.ResetRefreshTimeForBucket(bucket)
	}
}

func (dht *DHT) selectResultChan(
	resultChan chan *packet.Packet,
	futuresCount *int,
	results []*packet.Packet,
) ([]*packet.Packet, bool) {
	select {
	case result := <-resultChan:
		if result != nil {
			results = append(results, result)
		} else {
			*futuresCount--
		}
		if len(results) == *futuresCount {
			close(resultChan)
			return results, true
		}
	case <-time.After(dht.options.PacketTimeout):
		close(resultChan)
		return results, true
	}
	return results, false
}

func (dht *DHT) setUpResultChan(futures []transport.Future, ctx Context, resultChan chan *packet.Packet) {
	for _, f := range futures {
		go func(future transport.Future) {
			select {
			case result := <-future.Result():
				if result == nil {
					// Channel was closed
					return
				}
				dht.addHost(ctx, routing.NewRouteHost(result.Sender))
				resultChan <- result
				return
			case <-time.After(dht.options.PacketTimeout):
				future.Cancel()
				return
			}
		}(f)
	}
}

func (dht *DHT) sendPacketToAlphaHosts(
	routeSet *routing.RouteSet,
	queryRest bool,
	t routing.IterateType,
	ht *routing.HashTable,
	contacted map[string]bool,
	target []byte,
	futures []transport.Future,
	removeFromRouteSet []*host.Host,
) (resultFutures []transport.Future, resultRouteSet []*host.Host) {
	// Next we send Packets to the first (closest) alpha hosts in the
	// route set and wait for a response

	for i, receiver := range routeSet.Hosts() {
		// Contact only alpha hosts
		if i >= routing.ParallelCalls && !queryRest {
			break
		}

		// Don't contact hosts already contacted
		if (contacted)[string(receiver.ID.GetKey())] {
			continue
		}

		(contacted)[string(receiver.ID.GetKey())] = true

		packetBuilder := packet.NewBuilder().Sender(ht.Origin).Receiver(receiver)
		packetBuilder = getPacketBuilder(t, packetBuilder, target)
		msg := packetBuilder.Build()

		// Send the async queries and wait for a response
		res, err := dht.transport.SendRequest(msg)
		if err != nil {
			// Host was unreachable for some reason. We will have to remove
			// it from the route set, but we will keep it in our routing
			// table in hopes that it might come back online in the f.
			removeFromRouteSet = append(removeFromRouteSet, msg.Receiver)
			continue
		}

		futures = append(futures, res)
	}
	return futures, removeFromRouteSet
}

func getPacketBuilder(t routing.IterateType, packetBuilder packet.Builder, target []byte) packet.Builder {
	switch t {
	case routing.IterateBootstrap, routing.IterateFindHost:
		return packetBuilder.Type(packet.TypeFindHost).Request(&packet.RequestDataFindHost{Target: target})
	case routing.IterateFindValue:
		return packetBuilder.Type(packet.TypeFindValue).Request(&packet.RequestDataFindValue{Target: target})
	case routing.IterateStore:
		return packetBuilder.Type(packet.TypeFindHost).Request(&packet.RequestDataFindHost{Target: target})
	default:
		panic("Unknown iterate type")
	}
}

// addHost adds a host into the appropriate k bucket
// we store these buckets in big-endian order so we look at the bits
// from right to left in order to find the appropriate bucket
func (dht *DHT) addHost(ctx Context, host *routing.RouteHost) {
	ht := dht.htFromCtx(ctx)
	index := routing.GetBucketIndexFromDifferingBit(ht.Origin.ID.GetKey(), host.ID.GetKey())

	// Make sure host doesn't already exist
	// If it does, mark it as seen
	if ht.DoesHostExistInBucket(index, host.ID.GetKey()) {
		ht.MarkHostAsSeen(host.ID.GetKey())
		return
	}

	ht.Lock()
	defer ht.Unlock()

	bucket := ht.RoutingTable[index]

	if len(bucket) == routing.MaxContactsInBucket {
		// If the bucket is full we need to ping the first host to find out
		// if it responds back in a reasonable amount of time. If not -
		// we may remove it
		n := bucket[0].Host
		request := packet.NewPingPacket(ht.Origin, n)
		future, err := dht.transport.SendRequest(request)
		if err != nil {
			bucket = append(bucket, host)
			bucket = bucket[1:]
		} else {
			select {
			case <-future.Result():
				return
			case <-time.After(dht.options.PingTimeout):
				bucket = bucket[1:]
				bucket = append(bucket, host)
			}
		}
	} else {
		bucket = append(bucket, host)
	}

	ht.RoutingTable[index] = bucket
}

func (dht *DHT) handleDisconnect(start, stop chan bool) {
	multiplexCount := 0

	for {
		select {
		case <-start:
			multiplexCount++
		case <-dht.transport.Stopped():
			for i := 0; i < multiplexCount; i++ {
				stop <- true
			}
			dht.transport.Close()
			return
		}
	}
}

func (dht *DHT) handleStoreTimers(start, stop chan bool) {
	start <- true

	ticker := time.NewTicker(time.Second)
	cb := NewContextBuilder(dht)
	for {
		dht.selectTicker(ticker, &cb, stop)
	}
}

func (dht *DHT) selectTicker(ticker *time.Ticker, cb *ContextBuilder, stop chan bool) {
	select {
	case <-ticker.C:
		keys := dht.store.GetKeysReadyToReplicate()
		for _, ht := range dht.tables {
			ctx, err := cb.SetHostByID(ht.Origin.ID).Build()
			// TODO: do something sane with error
			if err != nil {
				log.Fatal(err)
			}
			// Refresh
			for i := 0; i < routing.KeyBitSize; i++ {
				if time.Since(ht.GetRefreshTimeForBucket(i)) > dht.options.RefreshTime {
					id1 := ht.GetRandomIDFromBucket(routing.MaxContactsInBucket)
					_, _, err = dht.iterate(ctx, routing.IterateBootstrap, id1, nil)
					if err != nil {
						continue
					}
				}
			}

			// Replication
			for _, key := range keys {
				value, _ := dht.store.Retrieve(key)
				_, _, err2 := dht.iterate(ctx, routing.IterateStore, key, value)
				if err2 != nil {
					continue
				}
			}
		}

		// Expiration
		dht.store.ExpireKeys()
	case <-stop:
		ticker.Stop()
		return
	}
}

func (dht *DHT) handlePackets(start, stop chan bool) {
	start <- true

	cb := NewContextBuilder(dht)
	for {
		select {
		case msg := <-dht.transport.Packets():
			if msg == nil || !msg.IsForMe(*dht.origin) {
				continue
			}

			var ctx Context
			ctx = buildContext(cb, msg)
			ht := dht.htFromCtx(ctx)

			if ht.Origin.ID.KeyEqual(msg.Receiver.ID.GetKey()) || !dht.relay.NeedToRelay(msg.Sender.Address.String()) {
				dht.dispatchPacketType(ctx, msg, ht)
			} else {
				targetHost, exist, err := dht.FindHost(ctx, msg.Receiver.ID.KeyString())
				if err != nil {
					log.Println(err)
				} else if !exist {
					log.Printf("Target host addr: %s, ID: %s not found", msg.Receiver.Address.String(), msg.Receiver.ID.KeyString())
				} else {
					// need to relay incoming packet
					request := &packet.Packet{Sender: &host.Host{Address: dht.origin.Address, ID: msg.Sender.ID},
						Receiver:  &host.Host{ID: msg.Receiver.ID, Address: targetHost.Address},
						Type:      msg.Type,
						RequestID: msg.RequestID,
						Data:      msg.Data}
					dht.sendRelayedRequest(request, ctx)
				}
			}
		case <-stop:
			return
		}
	}
}

func (dht *DHT) sendRelayedRequest(request *packet.Packet, ctx Context) {
	_, err := dht.transport.SendRequest(request)
	if err != nil {
		log.Println(err)
	}
}

func buildContext(cb ContextBuilder, msg *packet.Packet) Context {
	var ctx Context
	var err error
	if msg.Receiver.ID.GetKey() == nil {
		ctx, err = cb.SetDefaultHost().Build()
	} else {
		ctx, err = cb.SetHostByID(msg.Receiver.ID).Build()
	}
	if err != nil {
		// TODO: Do something sane with error!
		log.Println(err) // don't return this error cuz don't know what to do with
	}
	return ctx
}

func (dht *DHT) confirmNodeRole(roleKey string) bool {
	// TODO implement this func
	return true
}

// CheckNodeRole starting a check all known nodes.
func (dht *DHT) CheckNodeRole(ctx Context, domainID string) error {
	var err error
	// TODO: change or choose another auth host
	if len(dht.options.BootstrapHosts) > 0 {
		err = dht.checkNodePrivRequest(ctx, dht.options.BootstrapHosts[0].ID.KeyString(), domainID)
	} else {
		err = errors.New("bootstrap node not exist")
	}
	return err
}

// RelayRequest sends relay request to target.
func (dht *DHT) RelayRequest(ctx Context, command, targetID string) error {
	var typedCommand packet.CommandType
	targetHost, exist, err := dht.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("target for relay request not found")
		return err
	}

	switch command {
	case "start":
		typedCommand = packet.StartRelay
	case "stop":
		typedCommand = packet.StopRelay
	default:
		err = errors.New("unknown command")
		return err
	}
	request := packet.NewRelayPacket(typedCommand, dht.htFromCtx(ctx).Origin, targetHost)
	future, err := dht.transport.SendRequest(request)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	select {
	case rsp := <-future.Result():
		if rsp == nil {
			err = errors.New("chanel closed unexpectedly")
			return err
		}

		response := rsp.Data.(*packet.ResponseRelay)
		err = dht.handleRelayResponse(ctx, response, targetID)
		if err != nil {
			return err
		}

	case <-time.After(dht.options.PacketTimeout):
		future.Cancel()
		err = errors.New("timeout")
		return err
	}

	return nil
}

func (dht *DHT) handleRelayResponse(ctx Context, response *packet.ResponseRelay, targetID string) error {
	var err error
	switch response.State {
	case relay.Stopped:
		// stop use this address as relay
		dht.proxy.RemoveProxyHost(targetID)
		err = nil
	case relay.Started:
		// start use this address as relay
		dht.proxy.AddProxyHost(targetID)
		err = nil
	case relay.NoAuth:
		err = errors.New("unable to execute relay because this host not authenticated")
	case relay.Unknown:
		err = errors.New("unknown relay command")
	case relay.Error:
		err = errors.New("relay request error")
	default:
		// unknown state/failed to change state
		err = errors.New("unknown response state")
	}
	return err
}

func (dht *DHT) handleCheckOriginResponse(response *packet.ResponseCheckOrigin, targetID string) {
	if bytes.Equal(response.AuthUniqueKey, dht.auth.SentKeys[targetID]) {
		delete(dht.auth.SentKeys, targetID)
		dht.auth.authenticatedHosts[targetID] = true
	}
}

// CheckOriginRequest send a request to check target host originality
func (dht *DHT) CheckOriginRequest(ctx Context, targetID string) error {
	targetHost, exist, err := dht.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("target for relay request not found")
		return err
	}

	request := packet.NewCheckOriginPacket(dht.htFromCtx(ctx).Origin, targetHost)
	future, err := dht.transport.SendRequest(request)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	select {
	case rsp := <-future.Result():
		if rsp == nil {
			err = errors.New("chanel closed unexpectedly")
			return err
		}

		response := rsp.Data.(*packet.ResponseCheckOrigin)
		dht.handleCheckOriginResponse(response, targetID)

	case <-time.After(dht.options.PacketTimeout):
		future.Cancel()
		err = errors.New("timeout")
		return err
	}

	return nil
}

// AuthenticationRequest sends an authentication request.
func (dht *DHT) AuthenticationRequest(ctx Context, command, targetID string) error {
	targetHost, exist, err := dht.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("target for auth request not found")
		return err
	}

	origin := dht.htFromCtx(ctx).Origin
	var authCommand packet.CommandType
	switch command {
	case "begin":
		authCommand = packet.BeginAuth
	case "revoke":
		authCommand = packet.RevokeAuth
	default:
		err = errors.New("unknown command")
		return err
	}
	request := packet.NewAuthPacket(authCommand, origin, targetHost)
	future, err := dht.transport.SendRequest(request)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	select {
	case rsp := <-future.Result():
		if rsp == nil {
			err = errors.New("chanel closed unexpectedly")
			return err
		}

		response := rsp.Data.(*packet.ResponseAuth)
		err = dht.handleAuthResponse(response, targetHost.ID.KeyString())
		if err != nil {
			return err
		}

	case <-time.After(dht.options.PacketTimeout):
		future.Cancel()
		err = errors.New("timeout")
		return err
	}

	return nil
}

func (dht *DHT) checkNodePrivRequest(ctx Context, targetID string, roleKey string) error {
	targetHost, exist, err := dht.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("target for check node privileges request not found")
		return err
	}

	origin := dht.htFromCtx(ctx).Origin
	request := packet.NewCheckNodePrivPacket(origin, targetHost, roleKey)
	future, err := dht.transport.SendRequest(request)

	if err != nil {
		return err
	}

	select {
	case rsp := <-future.Result():
		if rsp == nil {
			err = errors.New("chanel closed unexpectedly")
			return err
		}

		response := rsp.Data.(*packet.ResponseCheckNodePriv)
		err = dht.handleCheckNodePrivResponse(response, roleKey)
		if err != nil {
			return err
		}

	case <-time.After(dht.options.PacketTimeout):
		future.Cancel()
		err = errors.New("timeout")
		return err
	}

	return nil
}

func (dht *DHT) handleCheckNodePrivResponse(response *packet.ResponseCheckNodePriv, toleKey string) error {
	switch response.State {
	case packet.Error:
		return errors.New(response.Error)
	case packet.Confirmed:
		return nil
	case packet.Declined:
		// TODO: set default unconfirmed role
		// dht.nodesMap[dht.nodesIDMap[roleKey]].SetNodeRole("Unconfirmed default role")
	}
	return nil
}

func (dht *DHT) handleAuthResponse(response *packet.ResponseAuth, target string) error {
	var err error
	if (len(response.AuthUniqueKey) != 0) && response.Success {
		dht.auth.mut.Lock()
		defer dht.auth.mut.Unlock()
		dht.auth.ReceivedKeys[target] = response.AuthUniqueKey
		err = nil
	} else {
		if response.Success && (len(response.AuthUniqueKey) == 0) { // revoke success
			return err
		}
		if !response.Success {
			err = errors.New("authentication unsuccessful")
		} else if len(response.AuthUniqueKey) == 0 {
			err = errors.New("wrong auth unique key received")
		}
	}
	return err
}

// ObtainIPRequest is request to self IP obtaining.
func (dht *DHT) ObtainIPRequest(ctx Context, targetID string) error {
	targetHost, exist, err := dht.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("target for relay request not found")
		return err
	}

	origin := dht.htFromCtx(ctx).Origin
	request := packet.NewObtainIPPacket(origin, targetHost)

	future, err := dht.transport.SendRequest(request)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	select {
	case rsp := <-future.Result():
		if rsp == nil {
			err = errors.New("chanel closed unexpectedly")
			return err
		}

		response := rsp.Data.(*packet.ResponseObtainIP)
		err = dht.handleObtainIPResponse(response, targetHost.ID.KeyString())
		if err != nil {
			return err
		}

	case <-time.After(dht.options.PacketTimeout):
		future.Cancel()
		err = errors.New("timeout")
		return err
	}

	return nil
}

func (dht *DHT) handleObtainIPResponse(response *packet.ResponseObtainIP, targetID string) error {
	if response.IP != "" {
		dht.subnet.SubnetIDs[response.IP] = append(dht.subnet.SubnetIDs[response.IP], targetID)
	} else {
		return errors.New("received empty IP")
	}
	return nil
}

// RemoteProcedureCall calls remote procedure on target host.
func (dht *DHT) RemoteProcedureCall(ctx Context, targetID string, method string, args [][]byte) (result []byte, err error) {
	targetHost, exists, err := dht.FindHost(ctx, targetID)
	ht := dht.htFromCtx(ctx)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("targetHost not found")
	}

	request := &packet.Packet{
		Sender:   ht.Origin,
		Receiver: targetHost,
		Type:     packet.TypeRPC,
		Data: &packet.RequestDataRPC{
			Method: method,
			Args:   args,
		},
	}

	if targetID == dht.GetOriginHost(ctx).ID.KeyString() {
		return dht.rpc.Invoke(request.Sender, method, args)
	}

	// Send the async queries and wait for a future
	future, err := dht.transport.SendRequest(request)
	if err != nil {
		return nil, err
	}

	select {
	case rsp := <-future.Result():
		if rsp == nil {
			// Channel was closed
			return nil, errors.New("chanel closed unexpectedly")
		}
		dht.addHost(ctx, routing.NewRouteHost(rsp.Sender))

		response := rsp.Data.(*packet.ResponseDataRPC)
		if response.Success {
			return response.Result, nil
		}
		return nil, errors.New(response.Error)
	case <-time.After(dht.options.PacketTimeout):
		future.Cancel()
		return nil, errors.New("timeout")
	}
}

// RemoteProcedureRegister registers procedure for remote call on this host
func (dht *DHT) RemoteProcedureRegister(name string, method core.RemoteProcedure) {
	rp := func(sender *host.Host, args [][]byte) ([]byte, error) {
		return method(args)
	}

	dht.rpc.RegisterMethod(name, rp)
}

// ObtainIP starts to self IP obtaining.
func (dht *DHT) ObtainIP(ctx Context) error {
	for _, table := range dht.tables {
		for i := range table.RoutingTable {
			for j := range table.RoutingTable[i] {
				err := dht.ObtainIPRequest(ctx, table.RoutingTable[i][j].ID.KeyString())
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// GetDistance returns a distance between id1 and id2.
func (dht *DHT) GetDistance(id1, id2 []byte) *big.Int {
	buf1 := new(big.Int).SetBytes(id1)
	buf2 := new(big.Int).SetBytes(id2)
	return new(big.Int).Xor(buf1, buf2)
}

func (dht *DHT) getHomeSubnetKey(ctx Context) (string, error) {
	var result string
	for key, subnet := range dht.subnet.SubnetIDs {
		first := key
		first = xstrings.Reverse(first)
		first = strings.SplitAfterN(first, ".", 2)[1] // remove X.X.X.this byte
		first = strings.SplitAfterN(first, ".", 2)[1] // remove X.X.this byte
		first = xstrings.Reverse(first)
		for _, id1 := range subnet {
			target, exist, err := dht.FindHost(ctx, id1)
			if err != nil {
				return "", err
			} else if !exist {
				return "", errors.New("couldn't find a host")
			}
			if !strings.Contains(target.Address.IP.String(), first) {
				result = ""
				break
			} else {
				result = key
			}
		}
	}
	return result, nil
}

func (dht *DHT) countOuterHosts() {
	if len(dht.subnet.SubnetIDs) > 1 {
		for key, hosts := range dht.subnet.SubnetIDs {
			if key == dht.subnet.HomeSubnetKey {
				continue
			}
			dht.subnet.HighKnownHosts.SelfKnownOuterHosts += len(hosts)
		}
	}
}

// AnalyzeNetwork is func to analyze the network after IP obtaining.
func (dht *DHT) AnalyzeNetwork(ctx Context) error {
	var err error
	dht.subnet.HomeSubnetKey, err = dht.getHomeSubnetKey(ctx)
	if err != nil {
		return err
	}
	dht.countOuterHosts()
	dht.subnet.HighKnownHosts.OuterHosts = dht.subnet.HighKnownHosts.SelfKnownOuterHosts
	hosts := dht.subnet.SubnetIDs[dht.subnet.HomeSubnetKey]
	for _, ids := range hosts {
		err = dht.knownOuterHostsRequest(ids, dht.subnet.HighKnownHosts.OuterHosts)
		if err != nil {
			return err
		}
	}
	if len(dht.subnet.SubnetIDs) == 1 {
		if dht.subnet.HomeSubnetKey == "" { // current host have a static IP
			for _, subnetIDs := range dht.subnet.SubnetIDs {
				dht.sendRelayOwnership(subnetIDs)
			}
		}
	}

	return nil
}

func (dht *DHT) sendRelayOwnership(subnetIDs []string) {
	for _, id1 := range subnetIDs {
		err := dht.relayOwnershipRequest(id1, true)
		log.Println(err.Error())
	}
}

func (dht *DHT) handleRelayOwnership(response *packet.ResponseRelayOwnership, target string) {
	if response.Accepted {
		dht.subnet.PossibleRelayIDs = append(dht.subnet.PossibleRelayIDs, target)
	}
}

func (dht *DHT) relayOwnershipRequest(targetID string, ready bool) error {
	ctx, err := NewContextBuilder(dht).SetDefaultHost().Build()
	if err != nil {
		return err
	}
	targetHost, exist, err := dht.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("target for relay request not found")
		return err
	}

	request := packet.NewRelayOwnershipPacket(dht.htFromCtx(ctx).Origin, targetHost, true)
	future, err := dht.transport.SendRequest(request)

	if err != nil {
		return err
	}

	select {
	case rsp := <-future.Result():
		if rsp == nil {
			return err
		}

		response := rsp.Data.(*packet.ResponseRelayOwnership)
		dht.handleRelayOwnership(response, targetID)

	case <-time.After(dht.options.PacketTimeout):
		future.Cancel()
		err = errors.New("timeout")
		return err
	}

	return nil
}

func (dht *DHT) knownOuterHostsRequest(targetID string, hosts int) error {
	ctx, err := NewContextBuilder(dht).SetDefaultHost().Build()
	if err != nil {
		return err
	}
	targetHost, exist, err := dht.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("target for relay request not found")
		return err
	}

	request := packet.NewKnownOuterHostsPacket(dht.htFromCtx(ctx).Origin, targetHost, hosts)
	future, err := dht.transport.SendRequest(request)

	if err != nil {
		return err
	}

	select {
	case rsp := <-future.Result():
		if rsp == nil {
			return err
		}

		response := rsp.Data.(*packet.ResponseKnownOuterHosts)
		err = dht.handleKnownOuterHosts(ctx, response, targetID)
		if err != nil {
			return err
		}

	case <-time.After(dht.options.PacketTimeout):
		future.Cancel()
		err = errors.New("timeout")
		return err
	}

	return nil
}

func (dht *DHT) handleKnownOuterHosts(ctx Context, response *packet.ResponseKnownOuterHosts, targetID string) error {
	var err error
	if response.OuterHosts > dht.subnet.HighKnownHosts.OuterHosts { // update data
		dht.subnet.HighKnownHosts.OuterHosts = response.OuterHosts
		dht.subnet.HighKnownHosts.ID = response.ID
	}
	if (response.OuterHosts > dht.subnet.HighKnownHosts.SelfKnownOuterHosts) &&
		(dht.proxy.ProxyHostsCount() == 0) {
		err = dht.AuthenticationRequest(ctx, "begin", targetID)
		if err != nil {
			return err
		}
		err = dht.RelayRequest(ctx, "start", targetID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dht *DHT) htFromCtx(ctx Context) *routing.HashTable {
	htIdx := ctx.Value(ctxTableIndex).(int)
	return dht.tables[htIdx]
}
