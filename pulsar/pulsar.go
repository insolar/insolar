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
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/gob"
	"encoding/json"
	"errors"
	"net"

	"github.com/cenkalti/rpc2"
	"github.com/insolar/insolar/configuration"
	"golang.org/x/crypto/sha3"
)

type Pulsar struct {
	Sock      net.Listener
	RpcServer *rpc2.Server

	Neighbours map[string]*Neighbour
	PrivateKey *rsa.PrivateKey
}

// Creation new pulsar-node
func NewPulsar(configuration configuration.Pulsar, listener func(string, string) (net.Listener, error)) (*Pulsar, error) {
	// Listen for incoming connections.
	l, err := listener(configuration.ConnectionType.String(), configuration.ListenAddress)
	if err != nil {
		return nil, err
	}

	// Parse private key from config
	privateKey, err := ParseRsaPrivateKeyFromPemStr(configuration.PrivateKey)
	if err != nil {
		return nil, err
	}
	pulsar := &Pulsar{Sock: l, Neighbours: map[string]*Neighbour{}}
	pulsar.PrivateKey = privateKey

	// Adding other pulsars
	for _, neighbour := range configuration.ListOfNeighbours {
		publicKey, err := ParseRsaPublicKeyFromPemStr(neighbour.PublicKey)
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

		err = checkSignature(request)
		if err != nil {
			return err
		}

		if neighbour.Client == nil {
			neighbour.Client = client
		}

		generator := StandardEntropyGenerator{}
		convertedKey, err := ExportRsaPublicKeyAsPemStr(&pulsar.PrivateKey.PublicKey)
		if err != nil {
			return nil
		}
		message := Payload{PublicKey: convertedKey, Body: HandshakePayload{Entropy: generator.GenerateEntropy()}}
		message.Signature, err = singData(pulsar.PrivateKey, message.Body)
		response = &message

		return nil
	}
}

func (pulsar *Pulsar) EstablishConnection(pubKey *rsa.PublicKey) error {
	converted, err := ExportRsaPublicKeyAsPemStr(pubKey)
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
	convertedKey, err := ExportRsaPublicKeyAsPemStr(&pulsar.PrivateKey.PublicKey)
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

	err = checkSignature(&rep)
	if err != nil {
		return err
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

func checkSignature(request *Payload) error {
	h := sha3.New256()
	data, _ := json.Marshal(request.Body)
	h.Write(data)
	digest := h.Sum(nil)
	publicKey, err := ParseRsaPublicKeyFromPemStr(request.PublicKey)
	if err != nil {
		return nil
	}
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA3_256, digest, request.Signature)
}

func singData(privateKey *rsa.PrivateKey, data interface{}) ([]byte, error) {
	h := sha3.New256()
	marshaledData, _ := json.Marshal(data)
	h.Write(marshaledData)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA3_256, d)
}
