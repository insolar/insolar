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

package servicenetwork

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/messagerouter"
	"github.com/insolar/insolar/network/nodenetwork"
)

// Service provides a route between MessageRouter and Nodenetwork.
type Service struct {
	references []core.RecordRef

	nodes        map[string]*nodenetwork.Node // key - node ID, value - node ptr.
	referenceMap map[string]string            // key - reference ID, value - node ID.
}

// NewService returns a new service.
func NewService() *Service {
	service := &Service{
		nodes:        make(map[string]*nodenetwork.Node),
		referenceMap: make(map[string]string),
		references:   make([]core.RecordRef, 0),
	}

	service.readFile("../configuration/testdata/ids.txt")
	return service
}

// AddNode adds a node to service.
func (service *Service) AddNode(node *nodenetwork.Node) {
	if node != nil {
		service.nodes[node.GetNodeID()] = node
	}
}

// SendMessage sends a message from MessageRouter.
func (service Service) SendMessage(reference record.Reference, msg messagerouter.Message) {
	args := make([][]byte, 0)
	args[0] = msg.Arguments
	domainID := string(reference.Domain.Hash[:bytes.IndexByte(reference.Domain.Hash, 0)])
	service.nodes[service.referenceMap[domainID]].SendPacket(msg.Method, args)
}

func (service *Service) readFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		service.referenceMap[line[0]] = line[1]
	}

	return nil
}
