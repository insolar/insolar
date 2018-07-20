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
	"context"
	"errors"
	"log"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/insolar/insolar/network/host/message"
	"github.com/insolar/insolar/network/host/node"
	"github.com/insolar/insolar/network/host/relay"
	"github.com/insolar/insolar/network/host/routing"
	"github.com/insolar/insolar/network/host/rpc"
	"github.com/insolar/insolar/network/host/store"
	"github.com/insolar/insolar/network/host/transport"
	"github.com/jbenet/go-base58"
)

// DHT represents the state of the local node in the distributed hash table.
type DHT struct {
	tables  []*routing.HashTable
	options *Options

	origin *node.Origin

	transport transport.Transport
	store     store.Store
	rpc       rpc.RPC
	relay     relay.Relay
	proxy     relay.Proxy
}

// Options contains configuration options for the local node.
type Options struct {
	// The nodes being used to bootstrap the insolar. Without a bootstrap
	// node there is no way to connect to the insolar. NetworkNodes can be
	// initialized via insolar.NewNode().
	BootstrapNodes []*node.Node

	// The time after which a key/value pair expires;
	// this is a time-to-live (TTL) from the original publication date.
	ExpirationTime time.Duration

	// Seconds after which an otherwise unaccessed bucket must be refreshed.
	RefreshTime time.Duration

	// The interval between Kademlia replication events, when a node is
	// required to publish its entire database.
	ReplicateTime time.Duration

	// The time after which the original publisher must
	// republish a key/value pair. Currently not implemented.
	RepublishTime time.Duration

	// The maximum time to wait for a response from a node before discarding
	// it from the bucket.
	PingTimeout time.Duration

	// The maximum time to wait for a response to any message.
	MessageTimeout time.Duration
}

// NewDHT initializes a new DHT node.
func NewDHT(store store.Store, origin *node.Origin, transport transport.Transport, rpc rpc.RPC, options *Options, proxy relay.Proxy) (dht *DHT, err error) {
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

	if options.MessageTimeout == 0 {
		options.MessageTimeout = time.Second * 10
	}

	return dht, nil
}

func newTables(origin *node.Origin) ([]*routing.HashTable, error) {
	tables := make([]*routing.HashTable, len(origin.IDs))

	for i, id := range origin.IDs {
		ht, err := routing.NewHashTable(id, origin.Address)
		if err != nil {
			return nil, err
		}

		tables[i] = ht
	}

	return tables, nil
}

func (dht *DHT) getExpirationTime(ctx context.Context, key []byte) time.Time {
	ht := dht.htFromCtx(ctx)

	bucket := routing.GetBucketIndexFromDifferingBit(key, ht.Origin.ID)
	var total int
	for i := 0; i < bucket; i++ {
		total += ht.GetTotalNodesInBucket(i)
	}
	closer := ht.GetAllNodesInBucketCloserThan(bucket, key)
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

// Store stores data on the insolar. This will trigger an iterateStore loop.
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

// FindNode returns target node's real insolar address.
func (dht *DHT) FindNode(ctx Context, key string) (*node.Node, bool, error) {
	keyBytes := base58.Decode(key)
	if len(keyBytes) != routing.MaxContactsInBucket {
		return nil, false, errors.New("invalid key")
	}
	ht := dht.htFromCtx(ctx)

	if ht.Origin.ID.Equal(keyBytes) {
		return ht.Origin, true, nil
	}

	var targetNode *node.Node
	var exists = false
	routeSet := ht.GetClosestContacts(1, keyBytes, nil)

	if routeSet.Len() > 0 && routeSet.FirstNode().ID.Equal(keyBytes) {
		targetNode = routeSet.FirstNode()
		exists = true
	} else if dht.proxy.ProxyNodesCount() > 0 {
		address, _ := node.NewAddress(dht.proxy.GetNextProxyAddress())
		targetNode = &node.Node{ID: keyBytes, Address: address}
		return targetNode, true, nil
	} else {
		log.Println("Node not found in routing table. Iterating through insolar...")
		_, closest, err := dht.iterate(ctx, routing.IterateFindNode, keyBytes, nil)
		if err != nil {
			return nil, false, err
		}
		for i := range closest {
			if closest[i].ID.Equal(keyBytes) {
				targetNode = closest[i]
				exists = true
			}
		}
	}

	return targetNode, exists, nil
}

// NumNodes returns the total number of nodes stored in the local routing table.
func (dht *DHT) NumNodes(ctx Context) int {
	ht := dht.htFromCtx(ctx)
	return ht.TotalNodes()
}

// GetOriginID returns the base58 encoded identifier of the local node.
func (dht *DHT) GetOriginID(ctx Context) string {
	ht := dht.htFromCtx(ctx)
	return ht.Origin.ID.String()
}

// Listen begins listening on the socket for incoming Messages.
func (dht *DHT) Listen() error {
	start := make(chan bool)
	stop := make(chan bool)

	go dht.handleDisconnect(start, stop)
	go dht.handleMessages(start, stop)
	go dht.handleStoreTimers(start, stop)

	return dht.transport.Start()
}

// Bootstrap attempts to bootstrap the insolar using the BootstrapNodes provided
// to the Options struct. This will trigger an iterateBootstrap to the provided
// BootstrapNodes.
func (dht *DHT) Bootstrap() error {
	if len(dht.options.BootstrapNodes) == 0 {
		return nil
	}
	var futures []transport.Future
	wg := &sync.WaitGroup{}
	cb := NewContextBuilder(dht)

	for _, ht := range dht.tables {
		dht.iterateHtSendRequest(ht, cb, wg, &futures)
	}

	for _, f := range futures {
		go func(future transport.Future) {
			select {
			case result := <-future.Result():
				// If result is nil, channel was closed
				if result != nil {
					ctx, err := cb.SetNodeByID(result.Receiver.ID).Build()
					// TODO: must return error here
					if err != nil {
						log.Fatal(err)
					}
					dht.addNode(ctx, routing.NewRouteNode(result.Sender))
				}
				wg.Done()
				return
			case <-time.After(dht.options.MessageTimeout):
				future.Cancel()
				wg.Done()
				return
			}
		}(f)
	}

	wg.Wait()
	err := dht.iterateHt(cb)
	return err
}

func (dht *DHT) iterateHt(cb ContextBuilder) error {
	for _, ht := range dht.tables {
		ctx, err := cb.SetNodeByID(ht.Origin.ID).Build()
		if err != nil {
			return err
		}

		if dht.NumNodes(ctx) > 0 {
			_, _, err = dht.iterate(ctx, routing.IterateBootstrap, ht.Origin.ID, nil)
			return err
		}
	}
	return nil
}

func (dht *DHT) iterateHtSendRequest(ht *routing.HashTable, cb ContextBuilder, wg *sync.WaitGroup, futures *[]transport.Future) {
	ctx, err := cb.SetNodeByID(ht.Origin.ID).Build()
	if err != nil {
		return
	}
	for _, bn := range dht.options.BootstrapNodes {
		request := message.NewPingMessage(ht.Origin, bn)

		if bn.ID == nil {
			res, err := dht.transport.SendRequest(request)
			if err != nil {
				continue
			}
			wg.Add(1)
			*futures = append(*futures, res)
		} else {
			routeNode := routing.NewRouteNode(bn)
			dht.addNode(ctx, routeNode)
		}
	}
}

// Disconnect will trigger a Stop from the insolar.
func (dht *DHT) Disconnect() {
	dht.transport.Stop()
}

// Iterate does an iterative search through the insolar. This can be done
// for multiple reasons. These reasons include:
//     iterateStore - Used to store new information in the insolar.
//     iterateFindNode - Used to find node in the insolar given node abstract address.
//     iterateFindValue - Used to find a value among the insolar given a key.
//     iterateBootstrap - Used to bootstrap the insolar.
func (dht *DHT) iterate(ctx Context, t routing.IterateType, target []byte, data []byte) (value []byte, closest []*node.Node, err error) {
	ht := dht.htFromCtx(ctx)
	routeSet := ht.GetClosestContacts(routing.ParallelCalls, target, []*node.Node{})

	// We keep track of nodes contacted so far. We don't contact the same node
	// twice.
	var contacted = make(map[string]bool)

	// According to the Kademlia white paper, after a round of FIND_NODE RPCs
	// fails to provide a node closer than closestNode, we should send a
	// FIND_NODE RPC to all remaining nodes in the route set that have not
	// yet been contacted.
	queryRest := false

	// We keep a reference to the closestNode. If after performing a search
	// we do not find a closer node, we stop searching.
	if routeSet.Len() == 0 {
		return nil, nil, nil
	}

	closestNode := routeSet.FirstNode()

	checkAndRefreshTieForBucket(t, ht, target)

	var removeFromRouteSet []*node.Node

	for {
		var futures []transport.Future
		var futuresCount int

		dht.sendMessageToAlphaNodes(routeSet, queryRest, t, ht, &contacted, target, &futures, &removeFromRouteSet)

		removeNodesFromRouteSet(routeSet, &removeFromRouteSet)

		futuresCount = len(futures)

		resultChan := make(chan *message.Message)
		dht.setUpResultChan(&futures, ctx, resultChan)

		// ---
		// var results []*message.Message
		// if futuresCount > 0 {
		// Loop:
		// 	for {
		// 		if dht.selectResultChan(resultChan, &futuresCount, &results) {
		// 			break Loop
		// 		}
		// 	}
		//
		// 	value, closest, err = resultsIterate(t, &results, routeSet, target)
		// 	if closest != nil {
		// 		return nil, closest, nil
		// 	} else if value != nil {
		// 		return value, nil, nil
		// 	}
		// }
		//
		// if !queryRest && routeSet.Len() == 0 {
		// 	return nil, nil, nil
		// }
		value, closest, err = dht.checkFuturesCountAndGo(t, &queryRest, routeSet, futuresCount, resultChan, target, &closest)
		if (err != nil) && (err.Error() != "nothing") {
			return value, closest, err
		}
		// ----

		sort.Sort(routeSet)

		// value, closest, err = dht.iterateIsDone(t, &queryRest, routeSet, data, ht, closestNode)
		// if (err != nil) && (strings.Compare("continue", err.Error()) == 0) {
		// 	continue
		// } else if (err != nil) && (strings.Compare("do nothing", err.Error()) != 0) {
		// 	return value, closest, err
		// }

		// If closestNode is unchanged then we are done
		if routeSet.FirstNode().ID.Equal(closestNode.ID) || queryRest {
			// We are done
			switch t {
			case routing.IterateBootstrap:
				if !queryRest {
					queryRest = true
					continue
				}
				return nil, routeSet.Nodes(), nil
			case routing.IterateFindNode, routing.IterateFindValue:
				return nil, routeSet.Nodes(), nil
			case routing.IterateStore:
				for i, receiver := range routeSet.Nodes() {
					if i >= routing.MaxContactsInBucket {
						return nil, nil, nil
					}

					msg := message.NewBuilder().Sender(ht.Origin).Receiver(receiver).Type(message.TypeStore).Request(
						&message.RequestDataStore{
							Data: data,
						}).Build()

					future, _ := dht.transport.SendRequest(msg)
					// We do not need to handle result of this message
					future.Cancel()
				}
				return nil, nil, nil
			}
		} else {
			closestNode = routeSet.FirstNode()
		}
	}
}

func (dht *DHT) checkFuturesCountAndGo(
	t routing.IterateType,
	queryRest *bool,
	routeSet *routing.RouteSet,
	futuresCount int,
	resultChan chan *message.Message,
	target []byte,
	close *[]*node.Node,
) (value []byte, closest []*node.Node, err error) {

	var results []*message.Message
	if futuresCount > 0 {
	Loop:
		for {
			if dht.selectResultChan(resultChan, &futuresCount, &results) {
				break Loop
			}
		}

		value, *close, err = resultsIterate(t, &results, routeSet, target)
		if *close != nil {
			return nil, *close, nil
		} else if value != nil {
			return value, nil, nil
		}
	}

	if !*queryRest && routeSet.Len() == 0 {
		return nil, nil, nil
	}
	err = errors.New("nothing")
	return nil, nil, err
}

// func (dht *DHT) iterateIsDone(
// 	t routing.IterateType,
// 	queryRest *bool,
// 	routeSet *routing.RouteSet,
// 	data []byte,
// 	ht *routing.HashTable,
// 	closestNode *node.Node) (value []byte, closest []*node.Node, err error) {
//
// 	if routeSet.FirstNode().ID.Equal(closestNode.ID) || *queryRest {
// 		switch t {
// 		case routing.IterateBootstrap:
// 			if !*queryRest {
// 				*queryRest = true
// 				err = errors.New("continue")
// 				return nil, nil, err
// 			}
// 			return nil, routeSet.Nodes(), nil
// 		case routing.IterateFindNode, routing.IterateFindValue:
// 			return nil, routeSet.Nodes(), nil
// 		case routing.IterateStore:
// 			for i, receiver := range routeSet.Nodes() {
// 				if i >= routing.MaxContactsInBucket {
// 					return nil, nil, nil
// 				}
//
// 				msg := message.NewBuilder().Sender(ht.Origin).Receiver(receiver).Type(message.TypeStore).Request(
// 					&message.RequestDataStore{
// 						Data: data,
// 					}).Build()
//
// 				future, _ := dht.transport.SendRequest(msg)
// 				// We do not need to handle result of this message
// 				future.Cancel()
// 			}
// 			return nil, nil, nil
// 		}
// 	} else {
// 		*closestNode = *routeSet.FirstNode()
// 	}
// 	err = errors.New("do nothing")
// 	return nil, nil, err
// }

func resultsIterate(
	t routing.IterateType,
	results *[]*message.Message,
	routeSet *routing.RouteSet,
	target []byte,
) (value []byte, closest []*node.Node, err error) {

	for _, result := range *results {
		if result.Error != nil {
			routeSet.Remove(routing.NewRouteNode(result.Sender))
			continue
		}
		switch t {
		case routing.IterateBootstrap, routing.IterateFindNode, routing.IterateStore:
			responseData := result.Data.(*message.ResponseDataFindNode)
			if len(responseData.Closest) > 0 && responseData.Closest[0].ID.Equal(target) {
				return nil, responseData.Closest, nil
			}
			routeSet.Extend(routing.RouteNodesFrom(responseData.Closest))
		case routing.IterateFindValue:
			responseData := result.Data.(*message.ResponseDataFindValue)
			routeSet.Extend(routing.RouteNodesFrom(responseData.Closest))
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

func checkAndRefreshTieForBucket(t routing.IterateType, ht *routing.HashTable, target []byte) {
	if t == routing.IterateBootstrap {
		bucket := routing.GetBucketIndexFromDifferingBit(target, ht.Origin.ID)
		ht.ResetRefreshTimeForBucket(bucket)
	}
}

func removeNodesFromRouteSet(routeSet *routing.RouteSet, removeFromRouteSet *[]*node.Node) {
	for _, n := range *removeFromRouteSet {
		routeSet.Remove(routing.NewRouteNode(n))
	}
}

func (dht *DHT) selectResultChan(resultChan chan *message.Message, futuresCount *int, results *[]*message.Message) bool {
	select {
	case result := <-resultChan:
		if result != nil {
			*results = append(*results, result)
		} else {
			*futuresCount--
		}
		if len(*results) == *futuresCount {
			close(resultChan)
			return true
		}
	case <-time.After(dht.options.MessageTimeout):
		close(resultChan)
		return true
	}
	return false
}

func (dht *DHT) setUpResultChan(futures *[]transport.Future, ctx Context, resultChan chan *message.Message) {
	for _, f := range *futures {
		go func(future transport.Future) {
			select {
			case result := <-future.Result():
				if result == nil {
					// Channel was closed
					return
				}
				dht.addNode(ctx, routing.NewRouteNode(result.Sender))
				resultChan <- result
				return
			case <-time.After(dht.options.MessageTimeout):
				future.Cancel()
				return
			}
		}(f)
	}
}

func (dht *DHT) sendMessageToAlphaNodes(
	routeSet *routing.RouteSet,
	queryRest bool,
	t routing.IterateType,
	ht *routing.HashTable,
	contacted *map[string]bool,
	target []byte,
	futures *[]transport.Future,
	removeFromRouteSet *[]*node.Node,
) {
	// Next we send Messages to the first (closest) alpha nodes in the
	// route set and wait for a response

	for i, receiver := range routeSet.Nodes() {
		// Contact only alpha nodes
		if i >= routing.ParallelCalls && !queryRest {
			break
		}

		// Don't contact nodes already contacted
		if (*contacted)[string(receiver.ID)] {
			continue
		}

		(*contacted)[string(receiver.ID)] = true

		messageBuilder := message.NewBuilder().Sender(ht.Origin).Receiver(receiver)

		checkRoutingIterateType(t, &messageBuilder, target)

		msg := messageBuilder.Build()

		// Send the async queries and wait for a response
		res, err := dht.transport.SendRequest(msg)
		if err != nil {
			// Node was unreachable for some reason. We will have to remove
			// it from the route set, but we will keep it in our routing
			// table in hopes that it might come back online in the f.
			*removeFromRouteSet = append(*removeFromRouteSet, msg.Receiver)
			continue
		}

		*futures = append(*futures, res)
	}
}

func checkRoutingIterateType(t routing.IterateType, messageBuilder *message.Builder, target []byte) {
	switch t {
	case routing.IterateBootstrap, routing.IterateFindNode:
		*messageBuilder = messageBuilder.Type(message.TypeFindNode).Request(&message.RequestDataFindNode{
			Target: target,
		})
	case routing.IterateFindValue:
		*messageBuilder = messageBuilder.Type(message.TypeFindValue).Request(&message.RequestDataFindValue{
			Target: target,
		})
	case routing.IterateStore:
		*messageBuilder = messageBuilder.Type(message.TypeFindNode).Request(&message.RequestDataFindNode{
			Target: target,
		})
	default:
		panic("Unknown iterate type")
	}
}

// addNode adds a node into the appropriate k bucket
// we store these buckets in big-endian order so we look at the bits
// from right to left in order to find the appropriate bucket
func (dht *DHT) addNode(ctx Context, node *routing.RouteNode) {
	ht := dht.htFromCtx(ctx)
	index := routing.GetBucketIndexFromDifferingBit(ht.Origin.ID, node.ID)

	// Make sure node doesn't already exist
	// If it does, mark it as seen
	if ht.DoesNodeExistInBucket(index, node.ID) {
		ht.MarkNodeAsSeen(node.ID)
		return
	}

	ht.Lock()
	defer ht.Unlock()

	bucket := ht.RoutingTable[index]

	if len(bucket) == routing.MaxContactsInBucket {
		// If the bucket is full we need to ping the first node to find out
		// if it responds back in a reasonable amount of time. If not -
		// we may remove it
		n := bucket[0].Node
		request := message.NewPingMessage(ht.Origin, n)
		future, err := dht.transport.SendRequest(request)
		if err != nil {
			bucket = append(bucket, node)
			bucket = bucket[1:]
		} else {
			select {
			case <-future.Result():
				return
			case <-time.After(dht.options.PingTimeout):
				bucket = bucket[1:]
				bucket = append(bucket, node)
			}
		}
	} else {
		bucket = append(bucket, node)
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
		dht.selectTicker(ticker, &cb, &stop)
	}
}

func (dht *DHT) selectTicker(ticker *time.Ticker, cb *ContextBuilder, stop *chan bool) {
	select {
	case <-ticker.C:
		keys := dht.store.GetKeysReadyToReplicate()
		for _, ht := range dht.tables {
			ctx, err := cb.SetNodeByID(ht.Origin.ID).Build()
			// TODO: do something sane with error
			if err != nil {
				log.Fatal(err)
			}
			// Refresh
			for i := 0; i < routing.KeyBitSize; i++ {
				if time.Since(ht.GetRefreshTimeForBucket(i)) > dht.options.RefreshTime {
					id := ht.GetRandomIDFromBucket(routing.MaxContactsInBucket)
					_, _, err = dht.iterate(ctx, routing.IterateBootstrap, id, nil)
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
	case <-(*stop):
		ticker.Stop()
		return
	}
}

func (dht *DHT) handleMessages(start, stop chan bool) {
	start <- true

	cb := NewContextBuilder(dht)
	for {
		select {
		case msg := <-dht.transport.Messages():
			if msg == nil || !msg.IsForMe(*dht.origin) {
				continue
			}

			var ctx Context
			ctx = checkRecvID(cb, msg)
			ht := dht.htFromCtx(ctx)

			if ht.Origin.ID.Equal(msg.Receiver.ID) || !dht.relay.NeedToRelay(msg.Sender.Address.String()) {
				dht.switchType(ctx, msg, ht)
			} else {
				targetNode, exist, err := dht.FindNode(ctx, msg.Receiver.ID.String())
				if err != nil {
					log.Println(err)
				} else if !exist {
					log.Printf("Target node addr: %s, ID: %s not found", msg.Receiver.Address.String(), msg.Receiver.ID.String())
				} else {
					// need to relay incoming message
					request := &message.Message{Sender: &node.Node{Address: dht.origin.Address, ID: msg.Sender.ID},
						Receiver:  &node.Node{ID: msg.Receiver.ID, Address: targetNode.Address},
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

func (dht *DHT) sendRelayedRequest(request *message.Message, ctx Context) {
	future, err := dht.transport.SendRequest(request)
	if err != nil {
		log.Println(err)
	}
	select {
	case rsp := <-future.Result():
		if rsp == nil {
			// Channel was closed
			log.Println("chanel closed unexpectedly")
		}
		dht.addNode(ctx, routing.NewRouteNode(rsp.Sender))

		response := rsp.Data.(*message.ResponseDataRPC)
		if response.Success {
			log.Println(response.Result)
		}
		log.Println(response.Error)
	case <-time.After(dht.options.MessageTimeout):
		future.Cancel()
		log.Println("timeout")
	}
}

func checkRecvID(cb ContextBuilder, msg *message.Message) Context {
	var ctx Context
	var err error
	if msg.Receiver.ID == nil {
		ctx, err = cb.SetDefaultNode().Build()
	} else {
		ctx, err = cb.SetNodeByID(msg.Receiver.ID).Build()
	}
	if err != nil {
		// TODO: Do something sane with error!
		log.Println(err) // don't return this error cuz don't know what to do with
	}
	return ctx
}

func (dht *DHT) switchType(ctx Context, msg *message.Message, ht *routing.HashTable) {
	messageBuilder := message.NewBuilder().Sender(ht.Origin).Receiver(msg.Sender).Type(msg.Type)
	switch msg.Type {
	case message.TypeFindNode:
		dht.processFindNode(ctx, msg, messageBuilder)
	case message.TypeFindValue:
		dht.processFindValue(ctx, msg, messageBuilder)
	case message.TypeStore:
		dht.processStore(ctx, msg, messageBuilder)
	case message.TypePing:
		dht.processPing(ctx, msg, messageBuilder)
	case message.TypeRPC:
		dht.processRPC(ctx, msg, messageBuilder)
	case message.TypeRelay:
		dht.processRelay(ctx, msg, messageBuilder)
	}
}

func (dht *DHT) processFindNode(ctx Context, msg *message.Message, messageBuilder message.Builder) {
	ht := dht.htFromCtx(ctx)
	data := msg.Data.(*message.RequestDataFindNode)
	dht.addNode(ctx, routing.NewRouteNode(msg.Sender))
	closest := ht.GetClosestContacts(routing.MaxContactsInBucket, data.Target, []*node.Node{msg.Sender})
	response := &message.ResponseDataFindNode{
		Closest: closest.Nodes(),
	}
	err := dht.transport.SendResponse(msg.RequestID, messageBuilder.Response(response).Build())
	if err != nil {
		log.Println("Failed to send response:", err.Error())
	}
}

func (dht *DHT) processFindValue(ctx Context, msg *message.Message, messageBuilder message.Builder) {
	ht := dht.htFromCtx(ctx)
	data := msg.Data.(*message.RequestDataFindValue)
	dht.addNode(ctx, routing.NewRouteNode(msg.Sender))
	value, exists := dht.store.Retrieve(data.Target)
	response := &message.ResponseDataFindValue{}
	if exists {
		response.Value = value
	} else {
		closest := ht.GetClosestContacts(routing.MaxContactsInBucket, data.Target, []*node.Node{msg.Sender})
		response.Closest = closest.Nodes()
	}
	err := dht.transport.SendResponse(msg.RequestID, messageBuilder.Response(response).Build())
	if err != nil {
		log.Println("Failed to send response:", err.Error())
	}
}

func (dht *DHT) processStore(ctx Context, msg *message.Message, messageBuilder message.Builder) {
	data := msg.Data.(*message.RequestDataStore)
	dht.addNode(ctx, routing.NewRouteNode(msg.Sender))
	key := store.NewKey(data.Data)
	expiration := dht.getExpirationTime(ctx, key)
	replication := time.Now().Add(dht.options.ReplicateTime)
	err := dht.store.Store(key, data.Data, replication, expiration, false)
	if err != nil {
		log.Println("Failed to store data:", err.Error())
	}
}

func (dht *DHT) processPing(ctx Context, msg *message.Message, messageBuilder message.Builder) {
	err := dht.transport.SendResponse(msg.RequestID, messageBuilder.Response(nil).Build())
	if err != nil {
		log.Println("Failed to send response:", err.Error())
	}
}

func (dht *DHT) processRPC(ctx Context, msg *message.Message, messageBuilder message.Builder) {
	data := msg.Data.(*message.RequestDataRPC)
	dht.addNode(ctx, routing.NewRouteNode(msg.Sender))
	result, err := dht.rpc.Invoke(msg.Sender, data.Method, data.Args)
	response := &message.ResponseDataRPC{
		Success: true,
		Result:  result,
		Error:   "",
	}
	if err != nil {
		response.Success = false
		response.Error = err.Error()
	}
	err = dht.transport.SendResponse(msg.RequestID, messageBuilder.Response(response).Build())
	if err != nil {
		log.Println("Failed to send response:", err.Error())
	}
}

// Precess relay request.
func (dht *DHT) processRelay(ctx Context, msg *message.Message, messageBuilder message.Builder) {
	data := msg.Data.(*message.RequestRelay)
	dht.addNode(ctx, routing.NewRouteNode(msg.Sender))

	var success bool
	var state relay.State
	var err error

	switch data.Command {
	case message.Start:
		err = dht.relay.AddClient(msg.Sender)
		success = true
		state = relay.Started
	case message.Stop:
		err = dht.relay.RemoveClient(msg.Sender)
		success = true
		state = relay.Stopped
	default:
		success = false
		state = relay.Unknown
	}

	if err != nil {
		success = false
		state = relay.Error
	}

	response := &message.ResponseRelay{
		Success: success,
		State:   state,
	}

	err = dht.transport.SendResponse(msg.RequestID, messageBuilder.Response(response).Build())
	if err != nil {
		log.Println("Failed to send response:", err.Error())
	}
}

// RelayRequest sends relay request to target.
func (dht *DHT) RelayRequest(ctx Context, command, target string) error {
	var typedCommand message.CommandType
	targetNode, exist, err := dht.FindNode(ctx, target)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("target for relay request not found")
		return err
	}

	switch command {
	case "start":
		typedCommand = message.Start
	case "stop":
		typedCommand = message.Stop
	default:
		err = errors.New("unknown command")
		return err
	}
	request := message.NewRelayMessage(typedCommand, dht.htFromCtx(ctx).Origin, targetNode)
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

		response := rsp.Data.(*message.ResponseRelay)
		dht.handleRelayResponse(response, target)

	case <-time.After(dht.options.MessageTimeout):
		future.Cancel()
		err = errors.New("timeout")
		return err
	}

	return nil
}

func (dht *DHT) handleRelayResponse(response *message.ResponseRelay, target string) {
	if response.Success {
		switch response.State {
		case relay.Stopped:
			// stop use this address as relay
			dht.proxy.RemoveProxyNode(target)
		case relay.Started:
			// start use this address as relay
			dht.proxy.AddProxyNode(target)
		default:
			// unknown state/failed to change state
			log.Println("unknown response state:")
		}
	} else {
		log.Printf("Unable to execute relay command on %s", target)
	}
}

// RemoteProcedureCall calls remote procedure on target node.
func (dht *DHT) RemoteProcedureCall(ctx Context, target string, method string, args [][]byte) (result []byte, err error) {
	targetNode, exists, err := dht.FindNode(ctx, target)
	ht := dht.htFromCtx(ctx)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("targetNode not found")
	}

	request := &message.Message{
		Sender:   ht.Origin,
		Receiver: targetNode,
		Type:     message.TypeRPC,
		Data: &message.RequestDataRPC{
			Method: method,
			Args:   args,
		},
	}

	if target == dht.GetOriginID(ctx) {
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
		dht.addNode(ctx, routing.NewRouteNode(rsp.Sender))

		response := rsp.Data.(*message.ResponseDataRPC)
		if response.Success {
			return response.Result, nil
		}
		return nil, errors.New(response.Error)
	case <-time.After(dht.options.MessageTimeout):
		future.Cancel()
		return nil, errors.New("timeout")
	}
}

func (dht *DHT) htFromCtx(ctx Context) *routing.HashTable {
	htIdx := ctx.Value(ctxTableIndex).(int)
	return dht.tables[htIdx]
}
