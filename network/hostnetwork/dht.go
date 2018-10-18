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

package hostnetwork

import (
	"bytes"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/huandu/xstrings"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/dns"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/hosthandler"
	"github.com/insolar/insolar/network/hostnetwork/id"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/relay"
	"github.com/insolar/insolar/network/hostnetwork/routing"
	"github.com/insolar/insolar/network/hostnetwork/store"
	"github.com/insolar/insolar/network/hostnetwork/transport"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/jbenet/go-base58"
	"github.com/pkg/errors"
)

// DHT represents the state of the local host in the distributed hash table.
type DHT struct {
	tables  []*routing.HashTable
	options *Options

	origin *host.Origin

	transport         transport.Transport
	store             store.Store
	ncf               hosthandler.NetworkCommonFacade
	relay             relay.Relay
	proxy             relay.Proxy
	auth              AuthInfo
	subnet            Subnet
	timeout           int // bootstrap reconnect timeout
	infinityBootstrap bool
	nodeID            *core.RecordRef
	activeNodeKeeper  consensus.NodeKeeper
	majorityRule      int
}

// AuthInfo collects some information about authentication.
type AuthInfo struct {
	// Sent/received unique auth keys.
	SentKeys     map[string][]byte
	ReceivedKeys map[string][]byte

	AuthenticatedHosts map[string]bool
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
func NewDHT(
	store store.Store,
	origin *host.Origin,
	transport transport.Transport,
	ncf hosthandler.NetworkCommonFacade,
	options *Options,
	proxy relay.Proxy,
	timeout int,
	infbootstrap bool,
	nodeID *core.RecordRef,
	majorityRule int,
) (dht *DHT, err error) {
	tables, err := newTables(origin)
	if err != nil {
		return nil, err
	}

	rel := relay.NewRelay()

	dht = &DHT{
		options:           options,
		origin:            origin,
		ncf:               ncf,
		transport:         transport,
		tables:            tables,
		store:             store,
		relay:             rel,
		proxy:             proxy,
		timeout:           timeout,
		infinityBootstrap: infbootstrap,
		nodeID:            nodeID,
		majorityRule:      majorityRule,
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

	dht.auth.AuthenticatedHosts = make(map[string]bool)
	dht.auth.SentKeys = make(map[string][]byte)
	dht.auth.ReceivedKeys = make(map[string][]byte)

	dht.subnet.SubnetIDs = make(map[string][]string)

	return dht, nil
}

func (dht *DHT) SetNodeKeeper(keeper consensus.NodeKeeper) {
	dht.activeNodeKeeper = keeper
	if dht.GetNetworkCommonFacade().GetConsensus() == nil {
		log.Warn("consensus is nil")
		return
	}
	dht.GetNetworkCommonFacade().GetConsensus().SetNodeKeeper(keeper)
}

func newTables(origin *host.Origin) ([]*routing.HashTable, error) {
	tables := make([]*routing.HashTable, len(origin.IDs))

	for i, id1 := range origin.IDs {
		ht, err := routing.NewHashTable(id1, origin.Address)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create HashTable")
		}

		tables[i] = ht
	}

	return tables, nil
}

// StoreData stores data on the network. This will trigger an iterateStore loop.
// The base58 encoded identifier will be returned if the store is successful.
func (dht *DHT) StoreData(ctx hosthandler.Context, data []byte) (id string, err error) {
	key := store.NewKey(data)
	expiration := dht.GetExpirationTime(ctx, key)
	replication := time.Now().Add(dht.options.ReplicateTime)
	err = dht.store.Store(key, data, replication, expiration, true)
	if err != nil {
		return "", errors.Wrap(err, "Failed to store data")
	}
	_, _, err = dht.iterate(ctx, routing.IterateStore, key, data)
	if err != nil {
		return "", errors.Wrap(err, "Failed to iterate")
	}
	str := base58.Encode(key)
	return str, nil
}

// Get retrieves data from the transport using key. Key is the base58 encoded
// identifier of the data.
func (dht *DHT) Get(ctx hosthandler.Context, key string) ([]byte, bool, error) {
	keyBytes := base58.Decode(key)
	if len(keyBytes) != routing.MaxContactsInBucket {
		return nil, false, errors.New("invalid key")
	}

	value, exists := dht.store.Retrieve(keyBytes)
	if !exists {
		var err error
		value, _, err = dht.iterate(ctx, routing.IterateFindValue, keyBytes, nil)
		if err != nil {
			return nil, false, errors.Wrap(err, "Failed to iterate")
		}
		if value != nil {
			exists = true
		}
	}

	return value, exists, nil
}

// NumHosts returns the total number of hosts stored in the local routing table.
func (dht *DHT) NumHosts(ctx hosthandler.Context) int {
	ht := dht.HtFromCtx(ctx)
	return ht.TotalHosts()
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
		log.Info("empty bootstrap hosts")
		return nil
	}
	dht.checkBootstrapHostsDomains(dht.options.BootstrapHosts)
	cb := NewContextBuilder(dht)

	for _, ht := range dht.tables {
		dht.iterateBootstrapHosts(ht, cb)
	}

	return dht.iterateHt(cb)
}

func (dht *DHT) checkBootstrapHostsDomains(hosts []*host.Host) {
	for _, hst := range hosts {
		ip, err := dns.GetIpFromDomain(hst.Address.String())
		if err != nil {
			log.Warn(err)
		}
		hst.Address, err = host.NewAddress(ip)
		if err != nil {
			log.Warn(err)
		}
	}
}

func (dht *DHT) GetHostsFromBootstrap() {
	cb := NewContextBuilder(dht)
	if len(dht.options.BootstrapHosts) == 0 {
		return
	}
	for _, ht := range dht.tables {
		dht.iterateHtGetNearestHosts(ht, cb)
	}
}

func (dht *DHT) iterateHtGetNearestHosts(ht *routing.HashTable, cb ContextBuilder) {
	ctx, err := cb.SetHostByID(ht.Origin.ID).Build()
	if err != nil {
		log.Errorf("Error sending GetNearestHosts packet: %s", err.Error())
		return
	}

	futures := make([]transport.Future, 0)

	for _, host := range dht.options.BootstrapHosts {
		p := packet.NewBuilder().Type(packet.TypeFindHost).Sender(ht.Origin).Receiver(host).
			Request(&packet.RequestDataFindHost{Target: ht.Origin.ID}).Build()
		f, err := dht.transport.SendRequest(p)
		if err != nil {
			log.Errorf("Error sending GetNearestHosts packet to host: %s", host.String())
			continue
		}
		futures = append(futures, f)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(futures))
	for _, f := range futures {
		go func(f transport.Future) {
			defer wg.Done()
			result, err := f.GetResult(dht.options.PacketTimeout)
			if err != nil {
				log.Errorln("Error getting nearest hosts:", err.Error())
				return
			}
			data := result.Data.(*packet.ResponseDataFindHost)
			for _, host := range data.Closest {
				dht.AddHost(ctx, routing.NewRouteHost(host))
				log.Debugf("Added host to DHT routing table: %s %s", host.ID, host.Address)
			}
		}(f)
	}
	wg.Wait()
}

func (dht *DHT) iterateHt(cb ContextBuilder) error {
	for _, ht := range dht.tables {
		ctx, err := cb.SetHostByID(ht.Origin.ID).Build()
		if err != nil {
			return errors.Wrap(err, "Failed to SetHostByID")
		}

		if dht.NumHosts(ctx) > 0 {
			_, _, err = dht.iterate(ctx, routing.IterateBootstrap, ht.Origin.ID.Bytes(), nil)
			return errors.Wrap(err, "Failed to iterate")
		}
	}
	return nil
}

func (dht *DHT) iterateBootstrapHosts(
	ht *routing.HashTable,
	cb ContextBuilder,
) {
	localwg := &sync.WaitGroup{}
	log.Info("bootstrapping to each known hosts.")
	for _, bh := range dht.options.BootstrapHosts {
		localwg.Add(1)
		go func(cb ContextBuilder, dht *DHT, bh *host.Host, ht *routing.HashTable, localwg *sync.WaitGroup) {
			counter := 1
			if dht.infinityBootstrap {
				log.Info("do infinity mode bootstrap.")
				for {
					if dht.gotBootstrap(ht, bh, cb) {
						localwg.Done()
						return
					}
					if counter < dht.timeout {
						counter = counter * 2
					}
					time.Sleep(time.Second * time.Duration(counter))
				}
			} else {
				log.Info("do one time mode bootstrap.")
				_ = dht.gotBootstrap(ht, bh, cb)
				localwg.Done()
			}
		}(cb, dht, bh, ht, localwg)
	}
	localwg.Wait()
}

func (dht *DHT) gotBootstrap(ht *routing.HashTable, bh *host.Host, cb ContextBuilder) bool {
	request := packet.NewPingPacket(ht.Origin, bh)
	if bh.ID.Bytes() == nil {
		log.Info("sending ping request")
		res, err := dht.transport.SendRequest(request)
		if err != nil {
			log.Error(err)
			return false
		}
		result, err := res.GetResult(dht.options.PingTimeout)
		if err != nil {
			log.Warn("gotBootstrap:", err.Error())
			return false
		}
		dht.updateBootstrapHost(result.Sender.Address.String(), result.Sender.ID)
		log.Info("checking response")
		if result == nil {
			log.Warn("gotBootstrap: result is nil")
			return false
		}
		ctx, err := cb.SetHostByID(result.Receiver.ID).Build()
		if err != nil {
			log.Error(err)
		}
		dht.AddHost(ctx, routing.NewRouteHost(result.Sender))
	} else {
		log.Info("bootstrap host known. creating new route host.")
		routeHost := routing.NewRouteHost(bh)
		ctx, err := cb.SetHostByID(ht.Origin.ID).Build()
		if err != nil {
			log.Error("failed to create a context")
			return false
		}
		dht.AddHost(ctx, routeHost)
	}
	return true
}

func (dht *DHT) updateBootstrapHost(bootstrapAddress string, bootstrapID id.ID) {
	for _, target := range dht.options.BootstrapHosts {
		if target.Address.String() == bootstrapAddress {
			target.ID = bootstrapID
		}
	}
}

// StartAuthorize start authorize to discovery nodes.
func (dht *DHT) StartAuthorize() error {
	// hack for zeronet
	if len(dht.options.BootstrapHosts) == 0 {
		return nil
	}

	discoveryNodesCount := len(dht.options.BootstrapHosts)
	ch := make(chan []*core.ActiveNode, discoveryNodesCount)
	for _, h := range dht.options.BootstrapHosts {
		go func(ch chan []*core.ActiveNode, h *host.Host) {
			activeNodes, err := GetNonceRequest(dht, h.ID.String())
			if err != nil {
				log.Warnf("error authorizing on %s host: %s", h, err.Error())
			}
			ch <- activeNodes
		}(ch, h)
	}

	receivedResults := make([][]*core.ActiveNode, 0)
	i := 0
LOOP:
	for {
		select {
		case activeNodeList := <-ch:
			receivedResults = append(receivedResults, activeNodeList)
			i++
			if i == discoveryNodesCount {
				break LOOP
			}
		case <-time.After(time.Minute):
			log.Warn("StartAuthorize: timeout exceeded")
			break LOOP
		}
	}

	if len(receivedResults) == 0 {
		return errors.New("StartAuthorize: No answers received from discovery nodes")
	}

	atLeastOneResultIsFine := false
	for _, result := range receivedResults {
		err := dht.AddActiveNodes(result)
		if err != nil {
			log.Error(err.Error())
		} else {
			atLeastOneResultIsFine = true
		}
	}
	if !atLeastOneResultIsFine {
		return errors.New("StartAuthorize: received active nodes do not pass majority rule")
	}
	return nil
}

func (dht *DHT) AddUnsync(nodeID core.RecordRef, roles []core.NodeRole, address string /*, publicKey *ecdsa.PublicKey*/) (chan *core.ActiveNode, error) {
	// TODO: return nodekeeper from helper method in HostHandler and remove this func and GetActiveNodes
	return dht.activeNodeKeeper.AddUnsync(nodeID, roles, address /*, publicKey*/)
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
func (dht *DHT) iterate(ctx hosthandler.Context, t routing.IterateType, target []byte, data []byte) (value []byte, closest []*host.Host, err error) {
	ht := dht.HtFromCtx(ctx)
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

	if routeSet.FirstHost().ID.Equal(closestHost.ID.Bytes()) || *(queryRest) {
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
				future, err := dht.transport.SendRequest(msg)
				if err != nil {
					return nil, nil, nil, errors.Wrap(err, "Failed transport to SendRequest")
				}
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
			return nil, close, errors.Wrap(err, "Failed to resultsIterate")
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
			if len(responseData.Closest) > 0 && responseData.Closest[0].ID.Equal(target) {
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
		bucket := routing.GetBucketIndexFromDifferingBit(target, ht.Origin.ID.Bytes())
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

func (dht *DHT) setUpResultChan(futures []transport.Future, ctx hosthandler.Context, resultChan chan *packet.Packet) {
	for _, f := range futures {
		go func(future transport.Future, ctx hosthandler.Context, resultChan chan *packet.Packet) {
			result, err := future.GetResult(dht.options.PacketTimeout)
			if err != nil {
				log.Warn("setUpResultChan future error:", err.Error())
				return
			}
			dht.AddHost(ctx, routing.NewRouteHost(result.Sender))
			resultChan <- result
		}(f, ctx, resultChan)
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
		if (contacted)[string(receiver.ID.Bytes())] {
			continue
		}

		(contacted)[string(receiver.ID.Bytes())] = true

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

			go func(msg *packet.Packet) {
				if msg == nil || !msg.IsForMe(*dht.origin) {
					return
				}

				var ctx hosthandler.Context
				ctx = BuildContext(cb, msg)
				ht := dht.HtFromCtx(ctx)

				if ht.Origin.ID.Equal(msg.Receiver.ID.Bytes()) || !dht.relay.NeedToRelay(msg.Sender.Address.String()) {
					dht.dispatchPacketType(ctx, msg, ht)
				} else {
					targetHost, exist, err := dht.FindHost(ctx, msg.Receiver.ID.String())
					if err != nil {
						log.Errorln(err)
					} else if !exist {
						log.Warnln("Target host addr: %s, ID: %s not found", msg.Receiver.Address.String(), msg.Receiver.ID.String())
					} else {
						// need to relay incoming packet
						request := &packet.Packet{Sender: &host.Host{Address: dht.origin.Address, ID: msg.Sender.ID},
							Receiver:  &host.Host{ID: msg.Receiver.ID, Address: targetHost.Address},
							Type:      msg.Type,
							RequestID: msg.RequestID,
							Data:      msg.Data}
						sendRelayedRequest(dht, request)
					}
				}
			}(msg)
		case <-stop:
			return
		}
	}
}

func (dht *DHT) dispatchPacketType(ctx hosthandler.Context, msg *packet.Packet, ht *routing.HashTable) {
	packetBuilder := packet.NewBuilder().Sender(ht.Origin).Receiver(msg.Sender).Type(msg.Type)

	// TODO: fix sign and check sign logic
	if msg.Type == packet.TypeRPC {
		data := msg.Data.(*packet.RequestDataRPC)
		signedMsg, err := message.Deserialize(bytes.NewBuffer(data.Args[0]))
		if err != nil {
			log.Error(err, "failed to parse incoming RPC")
			return
		}
		activeNode := dht.activeNodeKeeper.GetActiveNode(*data.NodeID)
		if activeNode == nil {
			log.Warn("couldn't check a sign from non active node")
			return
		}
		if !message.SignIsCorrect(signedMsg, activeNode.PublicKey) {
			log.Warn("RPC message not signed")
			return
		}
	}

	response, err := ParseIncomingPacket(dht, ctx, msg, packetBuilder)
	if err != nil {
		log.Errorln(err)
	} else if response != nil {
		err = dht.transport.SendResponse(msg.RequestID, response)
		if err != nil {
			log.Errorln("Failed to send response:", err.Error())
		}
	}
}

// CheckNodeRole starting a check all known nodes.
func (dht *DHT) CheckNodeRole(domainID string) error {
	var err error
	// TODO: change or choose another auth host
	if len(dht.options.BootstrapHosts) > 0 {
		err = checkNodePrivRequest(dht, dht.options.BootstrapHosts[0].ID.String())
	} else {
		err = errors.New("bootstrap node not exist")
	}
	return err
}

// RemoteProcedureRegister registers procedure for remote call on this host
func (dht *DHT) RemoteProcedureRegister(name string, method core.RemoteProcedure) {
	rp := func(sender *host.Host, args [][]byte) ([]byte, error) {
		return method(args)
	}

	dht.ncf.GetRPC().RegisterMethod(name, rp)
}

// ObtainIP starts to self IP obtaining.
func (dht *DHT) ObtainIP() error {
	for _, table := range dht.tables {
		for i := range table.RoutingTable {
			for j := range table.RoutingTable[i] {
				err := ObtainIPRequest(dht, table.RoutingTable[i][j].ID.String())
				if err != nil {
					return errors.Wrap(err, "Failed to ObtainIPRequest")
				}
			}
		}
	}
	return nil
}

// GetNetworkCommonFacade returns a networkcommonfacade ptr.
func (dht *DHT) GetNetworkCommonFacade() hosthandler.NetworkCommonFacade {
	return dht.ncf
}

func (dht *DHT) getHomeSubnetKey(ctx hosthandler.Context) (string, error) {
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
				return "", errors.Wrap(err, "Failed to FindHost")
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
func (dht *DHT) AnalyzeNetwork(ctx hosthandler.Context) error {
	var err error
	dht.subnet.HomeSubnetKey, err = dht.getHomeSubnetKey(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to getHomeSubnetKey")
	}
	dht.countOuterHosts()
	dht.subnet.HighKnownHosts.OuterHosts = dht.subnet.HighKnownHosts.SelfKnownOuterHosts
	hosts := dht.subnet.SubnetIDs[dht.subnet.HomeSubnetKey]
	for _, ids := range hosts {
		err = knownOuterHostsRequest(dht, ids, dht.subnet.HighKnownHosts.OuterHosts)
		if err != nil {
			return errors.Wrap(err, "Failed to knownOuterHostsRequest")
		}
	}
	if len(dht.subnet.SubnetIDs) == 1 {
		if dht.subnet.HomeSubnetKey == "" { // current host have a static IP
			for _, subnetIDs := range dht.subnet.SubnetIDs {
				SendRelayOwnership(dht, subnetIDs)
			}
		}
	}

	return nil
}

// GetActiveNodesList returns an active nodes list.
func (dht *DHT) GetActiveNodesList() []*core.ActiveNode {
	return dht.activeNodeKeeper.GetActiveNodes()
}

// AddActiveNodes adds an active nodes slice.
func (dht *DHT) AddActiveNodes(activeNodes []*core.ActiveNode) error {
	err := dht.checkMajorityRule(activeNodes)
	if err != nil {
		return err
	}
	if len(dht.activeNodeKeeper.GetActiveNodes()) > 0 {
		currentHash, err := consensus.CalculateHash(dht.activeNodeKeeper.GetActiveNodes())
		if err != nil {
			return err
		}
		newHash, err := consensus.CalculateHash(activeNodes)
		if err != nil {
			return err
		}
		if !bytes.Equal(currentHash, newHash) {
			// TODO: disconnect from all or what?
			return errors.New("two or more active node lists are different but majority check has passed")
		}
	} else {
		dht.activeNodeKeeper.AddActiveNodes(activeNodes)
	}
	return nil
}

// HtFromCtx returns a routing hashtable known by ctx.
func (dht *DHT) HtFromCtx(ctx hosthandler.Context) *routing.HashTable {
	htIdx := ctx.Value(ctxTableIndex).(int)
	return dht.tables[htIdx]
}

// ConfirmNodeRole is a node role confirmation.
func (dht *DHT) ConfirmNodeRole(roleKey string) bool {
	// TODO implement this func
	return true
}

// CascadeSendMessage sends a message to the next cascade layer.
func (dht *DHT) CascadeSendMessage(data core.Cascade, targetID string, method string, args [][]byte) error {
	return CascadeSendMessage(dht, data, targetID, method, args)
}

// FindHost returns target host's real network address.
func (dht *DHT) FindHost(ctx hosthandler.Context, key string) (*host.Host, bool, error) {
	keyBytes := base58.Decode(key)
	if len(keyBytes) != routing.MaxContactsInBucket {
		return nil, false, errors.New("invalid key")
	}
	ht := dht.HtFromCtx(ctx)

	if ht.Origin.ID.Equal(keyBytes) {
		return ht.Origin, true, nil
	}

	var targetHost *host.Host
	var exists = false
	routeSet := ht.GetClosestContacts(1, keyBytes, nil)

	if routeSet.Len() > 0 && routeSet.FirstHost().ID.Equal(keyBytes) {
		targetHost = routeSet.FirstHost()
		exists = true
	} else if dht.proxy.ProxyHostsCount() > 0 {
		address, err := host.NewAddress(dht.proxy.GetNextProxyAddress())
		if err != nil {
			return nil, false, errors.Wrap(err, "Failed to parse host address")
		}
		// TODO: current key insertion
		id1, err := id.NewID()
		if err != nil {
			return nil, false, errors.Wrap(err, "Failed to create host ID")
		}
		targetHost = &host.Host{ID: id1, Address: address}
		return targetHost, true, nil
	} else {
		log.Infoln("Host not found in routing table. Iterating through network...")
		_, closest, err := dht.iterate(ctx, routing.IterateFindHost, keyBytes, nil)
		if err != nil {
			return nil, false, errors.Wrap(err, "Failed to iterate")
		}
		for i := range closest {
			if closest[i].ID.Equal(keyBytes) {
				targetHost = closest[i]
				exists = true
			}
		}
	}

	return targetHost, exists, nil
}

// InvokeRPC - invoke a method to rpc.
func (dht *DHT) InvokeRPC(sender *host.Host, method string, args [][]byte) ([]byte, error) {
	return dht.ncf.GetRPC().Invoke(sender, method, args)
}

// AddHost adds a host into the appropriate k bucket
// we store these buckets in big-endian order so we look at the bits
// from right to left in order to find the appropriate bucket
func (dht *DHT) AddHost(ctx hosthandler.Context, host *routing.RouteHost) {
	ht := dht.HtFromCtx(ctx)
	index := routing.GetBucketIndexFromDifferingBit(ht.Origin.ID.Bytes(), host.ID.Bytes())

	// Make sure host doesn't already exist
	// If it does, mark it as seen
	if ht.DoesHostExistInBucket(index, host.ID.Bytes()) {
		ht.MarkHostAsSeen(host.ID.Bytes())
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
			_, err := future.GetResult(dht.options.PingTimeout)
			if err == nil {
				return
			}
			bucket = bucket[1:]
			bucket = append(bucket, host)
		}
	} else {
		bucket = append(bucket, host)
	}

	ht.RoutingTable[index] = bucket
}

// GetReplicationTime returns a interval between Kademlia replication events.
func (dht *DHT) GetReplicationTime() time.Duration {
	return dht.options.ReplicateTime
}

// GetExpirationTime returns a expiration time after which a key/value pair expires.
func (dht *DHT) GetExpirationTime(ctx hosthandler.Context, key []byte) time.Time {
	ht := dht.HtFromCtx(ctx)

	bucket := routing.GetBucketIndexFromDifferingBit(key, ht.Origin.ID.Bytes())
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

func (dht *DHT) checkMajorityRule(nodes []*core.ActiveNode) error {
	if len(nodes) < dht.majorityRule {
		return errors.New("failed majority role check")
	}

	count := 0
	for _, activeNode := range nodes {
		for _, bootstrapNode := range dht.options.BootstrapHosts {
			if strings.EqualFold(bootstrapNode.ID.String(), nodenetwork.ResolveHostID(&activeNode.NodeID)) {
				count++
			}
		}
	}

	if count < dht.majorityRule {
		return errors.New("discovery nodes count < majority number")
	}

	return nil
}

// StoreRetrieve should return the local key/value if it exists.
func (dht *DHT) StoreRetrieve(key store.Key) ([]byte, bool) {
	return dht.store.Retrieve(key)
}

// EqualAuthSentKey returns true if a given key equals to targetID key.
func (dht *DHT) EqualAuthSentKey(targetID string, key []byte) bool {
	return bytes.Equal(dht.auth.SentKeys[targetID], key)
}

// SendRequest sends a packet.
func (dht *DHT) SendRequest(packet *packet.Packet) (transport.Future, error) {
	metrics.NetworkPacketSentTotal.WithLabelValues(packet.Type.String()).Inc()
	return dht.transport.SendRequest(packet)
}

// Store should store a key/value pair for the local host with the
// given replication and expiration times.
func (dht *DHT) Store(key store.Key, data []byte, replication time.Time, expiration time.Time, publisher bool) error {
	return dht.store.Store(key, data, replication, expiration, publisher)
}

// AddPossibleProxyID adds an id which could be a proxy.
func (dht *DHT) AddPossibleProxyID(id string) {
	dht.subnet.PossibleProxyIDs = append(dht.subnet.PossibleProxyIDs, id)
}

// AddProxyHost adds a proxy host.
func (dht *DHT) AddProxyHost(targetID string) {
	dht.proxy.AddProxyHost(targetID)
}

// AddSubnetID adds a subnet ID.
func (dht *DHT) AddSubnetID(ip, targetID string) {
	dht.subnet.SubnetIDs[ip] = append(dht.subnet.SubnetIDs[ip], targetID)
}

// AddAuthSentKey adds a sent key to auth.
func (dht *DHT) AddAuthSentKey(id string, key []byte) {
	dht.auth.SentKeys[id] = key
}

// AddRelayClient adds a new relay client.
func (dht *DHT) AddRelayClient(host *host.Host) error {
	return dht.relay.AddClient(host)
}

// RemoteProcedureCall calls remote procedure on target host.
func (dht *DHT) RemoteProcedureCall(ctx hosthandler.Context, targetID string, method string, args [][]byte) (result []byte, err error) {
	targetHost, exists, err := dht.FindHost(ctx, targetID)
	ht := dht.HtFromCtx(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to FindHost")
	}

	if !exists {
		return nil, errors.New("targetHost not found")
	}

	request := &packet.Packet{
		Sender:   ht.Origin,
		Receiver: targetHost,
		Type:     packet.TypeRPC,
		Data: &packet.RequestDataRPC{
			NodeID: dht.nodeID,
			Method: method,
			Args:   args,
		},
	}

	if targetID == dht.GetOriginHost().IDs[0].String() {
		return dht.ncf.GetRPC().Invoke(request.Sender, method, args)
	}

	// Send the async queries and wait for a future
	future, err := dht.transport.SendRequest(request)
	if err != nil {
		return nil, errors.Wrap(err, "Failed transport to send request")
	}

	rsp, err := future.GetResult(dht.options.PacketTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "RemoteProcedureCall error")
	}
	dht.AddHost(ctx, routing.NewRouteHost(rsp.Sender))

	response := rsp.Data.(*packet.ResponseDataRPC)
	if response.Success {
		return response.Result, nil
	}
	return nil, errors.New(response.Error)
}

// AddReceivedKey adds a new received key from target.
func (dht *DHT) AddReceivedKey(target string, key []byte) {
	dht.auth.ReceivedKeys[target] = key
}

// RemoveAuthHost removes a host from auth.
func (dht *DHT) RemoveAuthHost(key string) {
	delete(dht.auth.AuthenticatedHosts, key)
}

// RemoveProxyHost removes host from proxy list.
func (dht *DHT) RemoveProxyHost(targetID string) {
	dht.proxy.RemoveProxyHost(targetID)
}

// RemovePossibleProxyID removes if from possible proxy ids list.
func (dht *DHT) RemovePossibleProxyID(id string) {
	for i, proxy := range dht.subnet.PossibleProxyIDs {
		if id == proxy {
			dht.subnet.PossibleProxyIDs = append(dht.subnet.PossibleProxyIDs[:i], dht.subnet.PossibleProxyIDs[i+1:]...)
			return
		}
	}
}

// RemoveAuthSentKeys removes a targetID from sent keys.
func (dht *DHT) RemoveAuthSentKeys(targetID string) {
	delete(dht.auth.SentKeys, targetID)
}

// RemoveRelayClient removes a client from relay list.
func (dht *DHT) RemoveRelayClient(host *host.Host) error {
	return dht.relay.RemoveClient(host)
}

func (dht *DHT) SetNodeID(nodeID *core.RecordRef) {
	dht.nodeID = nodeID
}

// SetHighKnownHostID sets a new high known host ID.
func (dht *DHT) SetHighKnownHostID(id string) {
	dht.subnet.HighKnownHosts.ID = id
}

// SetOuterHostsCount sets a new value to outer hosts count.
func (dht *DHT) SetOuterHostsCount(hosts int) {
	dht.subnet.HighKnownHosts.OuterHosts = hosts
}

// SetAuthStatus sets a new auth status to targetID.
func (dht *DHT) SetAuthStatus(targetID string, status bool) {
	dht.auth.AuthenticatedHosts[targetID] = status
}

// GetProxyHostsCount returns a proxy hosts count.
func (dht *DHT) GetProxyHostsCount() int {
	return dht.proxy.ProxyHostsCount()
}

// GetOuterHostsCount returns a outer hosts count.
func (dht *DHT) GetOuterHostsCount() int {
	return dht.subnet.HighKnownHosts.OuterHosts
}

// GetSelfKnownOuterHosts return a self known hosts count.
func (dht *DHT) GetSelfKnownOuterHosts() int {
	return dht.subnet.HighKnownHosts.SelfKnownOuterHosts
}

// GetHighKnownHostID returns a high known host ID.
func (dht *DHT) GetHighKnownHostID() string {
	return dht.subnet.HighKnownHosts.ID
}

// GetPacketTimeout returns the maximum time to wait for a response to any packet.
func (dht *DHT) GetPacketTimeout() time.Duration {
	return dht.options.PacketTimeout
}

// GetNodeID returns a node ID.
func (dht *DHT) GetNodeID() *core.RecordRef {
	return dht.nodeID
}

// KeyIsReceived returns true and a key from targetID if exist.
func (dht *DHT) KeyIsReceived(key string) ([]byte, bool) {
	if key, ok := dht.auth.ReceivedKeys[key]; ok {
		return key, ok
	}
	return nil, false
}

// HostIsAuthenticated returns true if target ID is authenticated host.
func (dht *DHT) HostIsAuthenticated(targetID string) bool {
	if _, ok := dht.auth.AuthenticatedHosts[targetID]; ok {
		return true
	}
	return false
}

// AddPossibleRelayID add a host id which can be a relay.
func (dht *DHT) AddPossibleRelayID(id string) {
	dht.subnet.PossibleRelayIDs = append(dht.subnet.PossibleRelayIDs, id)
}

// GetOriginHost returns the local host.
func (dht *DHT) GetOriginHost() *host.Origin {
	return dht.origin
}
