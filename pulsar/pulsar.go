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

package pulsar

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"math/big"
	"net"

	"github.com/cenkalti/rpc2"
	"github.com/insolar/insolar/configuration"
	"golang.org/x/crypto/sha3"
)

type Pulsar struct {
	Sock      net.Listener
	RpcServer *rpc2.Server

	Neighbours map[string]*Neighbour
	PrivateKey *ecdsa.PrivateKey
}

// Creation new pulsar-node
func NewPulsar(configuration configuration.Pulsar, listener func(string, string) (net.Listener, error)) (*Pulsar, error) {
	// Listen for incoming connections.
	l, err := listener(configuration.ConnectionType.String(), configuration.ListenAddress)
	if err != nil {
		return nil, err
	}

	// Parse private key from config
	privateKey, err := importPrivateKey(configuration.PrivateKey)
	if err != nil {
		return nil, err
	}
	pulsar := &Pulsar{Sock: l, Neighbours: map[string]*Neighbour{}}
	pulsar.PrivateKey = privateKey

	// Adding other pulsars
	for _, neighbour := range configuration.ListOfNeighbours {
		if len(neighbour.PublicKey) == 0 {
			continue
		}
		publicKey, err := importPublicKey(neighbour.PublicKey)
		if err != nil {
			continue
		}
		pulsar.Neighbours[neighbour.PublicKey] = &Neighbour{
			ConnectionType:    neighbour.ConnectionType,
			ConnectionAddress: neighbour.Address,
			PublicKey:         publicKey}
	}

	gob.Register(Payload{})
	gob.Register(HandshakePayload{})

	return pulsar, nil
}

func (pulsar *Pulsar) Start() {
	// Adding rpc-server listener
	srv := rpc2.NewServer()
	ConfigureHandlers(pulsar, srv)
	pulsar.RpcServer = srv
	srv.Accept(pulsar.Sock)
}

func ConfigureHandlers(pulsar *Pulsar, server *rpc2.Server) {
	server.Handle(Handshake.String(), pulsar.HandshakeHandler())
}

func ConfigureHandlersForClient(pulsar *Pulsar, server *rpc2.Client) {
	server.Handle(Handshake.String(), pulsar.HandshakeHandler())
}

func (pulsar *Pulsar) HandshakeHandler() func(client *rpc2.Client, request *Payload, response *Payload) error {

	return func(client *rpc2.Client, request *Payload, response *Payload) error {
		neighbour, err := pulsar.fetchNeighbour(request.PublicKey)
		if err != nil {
			return err
		}

		result, err := checkSignature(request)
		if err != nil {
			return err
		}
		if !result {
			return errors.New("Signature check failed")
		}

		if neighbour.Client == nil {
			neighbour.Client = client
		}

		generator := StandardEntropyGenerator{}
		convertedKey, err := exportPublicKey(&pulsar.PrivateKey.PublicKey)
		if err != nil {
			return err
		}
		message := Payload{PublicKey: convertedKey, Body: HandshakePayload{Entropy: generator.GenerateEntropy()}}
		message.Signature, err = singData(pulsar.PrivateKey, message.Body)
		if err != nil {
			return err
		}
		*response = message

		return nil
	}
}

func (pulsar *Pulsar) EstablishConnection(pubKey *ecdsa.PublicKey) error {
	converted, err := exportPublicKey(pubKey)
	if err != nil {
		return err
	}
	neighbour, err := pulsar.fetchNeighbour(converted)
	if err != nil {
		return err
	}
	if neighbour.Client != nil {
		return nil
	}

	conn, err := net.Dial(neighbour.ConnectionType.String(), neighbour.ConnectionAddress)
	if err != nil {
		return err
	}

	clt := rpc2.NewClient(conn)
	ConfigureHandlersForClient(pulsar, clt)
	go clt.Run()

	generator := StandardEntropyGenerator{}
	convertedKey, err := exportPublicKey(&pulsar.PrivateKey.PublicKey)
	if err != nil {
		return nil
	}
	var rep Payload
	message := Payload{PublicKey: convertedKey, Body: HandshakePayload{Entropy: generator.GenerateEntropy()}}
	message.Signature, err = singData(pulsar.PrivateKey, message.Body)
	if err != nil {
		return err
	}
	err = clt.Call(Handshake.String(), message, &rep)
	if err != nil {
		return err
	}

	result, err := checkSignature(&rep)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("Signature check failed")
	}
	neighbour.Client = clt

	return nil
}

func (pulsar *Pulsar) fetchNeighbour(pubKey string) (*Neighbour, error) {
	neighbour, ok := pulsar.Neighbours[pubKey]
	if !ok {
		return nil, errors.New("Forbidden connection")
	}

	return neighbour, nil
}

func checkSignature(request *Payload) (bool, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(request.Body)
	if err != nil {
		return false, err
	}

	r := big.Int{}
	s := big.Int{}
	sigLen := len(request.Signature)
	r.SetBytes(request.Signature[:(sigLen / 2)])
	s.SetBytes(request.Signature[(sigLen / 2):])

	h := sha3.New256()
	h.Write(b.Bytes())
	hash := h.Sum(nil)
	publicKey, err := importPublicKey(request.PublicKey)
	if err != nil {
		return false, nil
	}

	return ecdsa.Verify(publicKey, hash, &r, &s), nil
}

func singData(privateKey *ecdsa.PrivateKey, data interface{}) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(data)
	if err != nil {
		return nil, err
	}

	h := sha3.New256()
	h.Write(b.Bytes())
	hash := h.Sum(nil)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash)
	if err != nil {
		return nil, err
	}

	return append(r.Bytes(), s.Bytes()...), nil
}
