// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package tests

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	"github.com/insolar/insolar/network/consensus/common/endpoints"
)

type NetStrategy interface {
	GetLinkStrategy(hostAddress endpoints.Name) LinkStrategy
}

type PacketFunc func(packet *Packet)

type LinkStrategy interface {
	BeforeSend(packet *Packet, out PacketFunc)
	BeforeReceive(packet *Packet, out PacketFunc)
}

type Packet struct {
	Payload interface{}
	Host    endpoints.Name
}

type EmuRoute struct {
	host     endpoints.Name
	network  *EmuNetwork
	strategy LinkStrategy
	toHost   chan<- Packet
	fromHost <-chan Packet
}

type EmuNetwork struct {
	hostsSync sync.RWMutex
	ctx       context.Context
	hosts     map[endpoints.Name]*EmuRoute
	strategy  NetStrategy
	running   bool
	bufSize   int
}

type errEmuNetwork struct {
	errType string
	details interface{}
}

func (e errEmuNetwork) Error() string {
	return fmt.Sprintf("emu-net error - %s: %v", e.errType, e.details)
}

func ErrUnknownEmuHost(host endpoints.Name) error {
	return errEmuNetwork{errType: "Unknown host", details: host}
}

func NewEmuNetwork(nwStrategy NetStrategy, ctx context.Context) *EmuNetwork {
	return &EmuNetwork{strategy: nwStrategy, ctx: ctx}
}

func (emuNet *EmuNetwork) AddHost(ctx context.Context, host endpoints.Name) (toHost <-chan Packet, fromHost chan<- Packet) {
	emuNet.hostsSync.Lock()
	defer emuNet.hostsSync.Unlock()

	_, isPresent := emuNet.hosts[host]
	if isPresent {
		panic(fmt.Sprintf("Duplicate host: %v", host))
	}

	var routeStrategy LinkStrategy
	if emuNet.strategy != nil {
		routeStrategy = emuNet.strategy.GetLinkStrategy(host)
	}
	if routeStrategy == nil {
		routeStrategy = stubLinkStrategyValue
	}

	chanBufSize := emuNet.bufSize
	if chanBufSize <= 0 {
		chanBufSize = 10
	}

	fromHostC := make(chan Packet, chanBufSize)
	toHostC := make(chan Packet, chanBufSize)

	if emuNet.hosts == nil {
		emuNet.hosts = make(map[endpoints.Name]*EmuRoute)
	}

	route := EmuRoute{host: host, strategy: routeStrategy, toHost: toHostC, fromHost: fromHostC, network: emuNet}
	emuNet.hosts[host] = &route

	if emuNet.running {
		go route.run(ctx)
	}

	return toHostC, fromHostC
}

func (emuNet *EmuNetwork) DropHost(host endpoints.Name) bool {
	route := emuNet.getHostRoute(host)
	if route == nil {
		return false
	}

	route.closeRoute()
	return true
}

func (emuNet *EmuNetwork) SendToHost(host endpoints.Name, payload interface{}, fromHost endpoints.Name) bool {
	route := emuNet.getHostRoute(host)
	if route == nil {
		return false
	}

	targetPacket := Packet{Payload: payload, Host: fromHost}
	route.pushPacket(targetPacket)
	return true
}

func (emuNet *EmuNetwork) SendToAll(payload interface{}, fromHost endpoints.Name) {
	for _, route := range emuNet.getRoutes() {
		targetPacket := Packet{Payload: payload, Host: fromHost}
		route.pushPacket(targetPacket)
	}
}

func (emuNet *EmuNetwork) SendRandom(payload interface{}, fromHost endpoints.Name) {
	targetPacket := Packet{Payload: payload, Host: fromHost}
	routes := emuNet.getRoutes()
	routes[rand.Intn(len(routes))].pushPacket(targetPacket)
}

func (emuNet *EmuNetwork) CreateSendToAllChannel() chan<- Packet {
	inbound := make(chan Packet)
	go func() {
		for {
			inboundPacket, ok := <-inbound
			if !ok {
				return
			}
			emuNet.SendToAll(inboundPacket.Payload, inboundPacket.Host)
		}
	}()
	return inbound
}

func (emuNet *EmuNetwork) CreateChannelSendToAllFromOne(sender endpoints.Name) chan<- interface{} {
	inbound := make(chan interface{})
	go func() {
		for {
			payload, ok := <-inbound
			if !ok {
				return
			}
			emuNet.SendToAll(payload, sender)
		}
	}()
	return inbound
}

func (emuNet *EmuNetwork) CreateChannelSendToRandom(sender endpoints.Name, attempts int) chan<- interface{} {
	inbound := make(chan interface{})
	go func() {
		for {
			payload, ok := <-inbound
			if !ok {
				return
			}
			for i := 0; i < attempts; i++ {
				emuNet.SendRandom(payload, sender)
			}
		}
	}()
	return inbound
}

func (emuNet *EmuNetwork) GetHosts() []*endpoints.Name {
	emuNet.hostsSync.RLock()
	defer emuNet.hostsSync.RUnlock()

	keys := make([]*endpoints.Name, 0, len(emuNet.hosts))
	for k := range emuNet.hosts {
		keys = append(keys, &k)
	}

	return keys
}

func (emuNet *EmuNetwork) getRoutes() []*EmuRoute {
	emuNet.hostsSync.RLock()
	defer emuNet.hostsSync.RUnlock()

	routes := make([]*EmuRoute, 0, len(emuNet.hosts))
	for _, v := range emuNet.hosts {
		routes = append(routes, v)
	}

	return routes
}

func (emuNet *EmuNetwork) Start(ctx context.Context) {
	emuNet.hostsSync.Lock()
	defer emuNet.hostsSync.Unlock()

	if emuNet.running {
		return
	}
	emuNet.running = true

	for _, route := range emuNet.hosts {
		go route.run(ctx)
	}
}

func (emuNet *EmuNetwork) getHostRoute(host endpoints.Name) *EmuRoute {
	emuNet.hostsSync.RLock()
	defer emuNet.hostsSync.RUnlock()

	return emuNet.hosts[host]
}

func (emuNet *EmuNetwork) internalRemoveHost(route *EmuRoute) {
	emuNet.hostsSync.Lock()
	defer emuNet.hostsSync.Unlock()
	delete(emuNet.hosts, route.host)
}

func (emuRt *EmuRoute) run(ctx context.Context) {
	defer emuRt.closeRoute()

	for {
		select {
		case <-ctx.Done():
			return
		case originPacket, ok := <-emuRt.fromHost:
			if !ok {
				return
			}
			// strategy can modify target and payload of a packet before delivery
			emuRt.strategy.BeforeSend(&originPacket, emuRt._sendPacket)
		}
	}
}

func (emuRt *EmuRoute) pushPacket(packet Packet) {
	emuRt.strategy.BeforeReceive(&packet, emuRt._recvPacket)
}

func (emuRt *EmuRoute) _sendPacket(originPacket *Packet) {

	targetPacket := Packet{Payload: originPacket.Payload, Host: emuRt.host}

	var outRoute *EmuRoute
	if originPacket.Host.IsLocalHost() || originPacket.Host == emuRt.host {
		outRoute = emuRt
	} else {
		outRoute = emuRt.network.getHostRoute(originPacket.Host)
		if outRoute == nil {
			targetPacket.Payload = ErrUnknownEmuHost(originPacket.Host)
			// inbound strategy MUST NOT be applied to error replies
			emuRt.toHost <- targetPacket
			return
		}
	}

	outRoute.pushPacket(targetPacket)
}

func (emuRt *EmuRoute) _recvPacket(originPacket *Packet) {
	defer func() {
		if recover() == nil {
			return
		}
		outRoute := emuRt.network.getHostRoute(originPacket.Host)
		if outRoute == nil {
			// the sender receiver is not available anymore
			return
		}
		targetPacket := Packet{Payload: ErrUnknownEmuHost(emuRt.host), Host: outRoute.host}
		outRoute.toHost <- targetPacket
	}()
	emuRt.toHost <- *originPacket
}

func (emuRt *EmuRoute) closeRoute() {
	defer func() {
		recover()
		if emuRt.network != nil {
			emuRt.network.internalRemoveHost(emuRt)
		}
	}()
	close(emuRt.toHost)
}

var stubLinkStrategyValue LinkStrategy = &stubLinkStrategy{}

type stubLinkStrategy struct{}

func (stubLinkStrategy) BeforeSend(packet *Packet, out PacketFunc) {
	out(packet)
}

func (stubLinkStrategy) BeforeReceive(packet *Packet, out PacketFunc) {
	out(packet)
}
