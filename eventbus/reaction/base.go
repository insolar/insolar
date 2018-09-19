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

// Package reaction represents responses to messages of the eventbus
package reaction

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

// Type is a enum type of reaction
type Type byte

const (
	// WrongReactionType - incorrect type (0)
	WrongReactionType = Type(iota)
	// CommonReactionType - two binary fields: data and results
	CommonReactionType
	// ObjectBodyReactionType - reaction with body, class reference, code reference ...
	ObjectBodyReactionType
)

func getEmptyReaction(t Type) (core.Reaction, error) {
	switch t {
	case WrongReactionType:
		return nil, errors.New("no empty reaction for 'wrong' reaction")
	case CommonReactionType:
		return &CommonReaction{}, nil
	case ObjectBodyReactionType:
		return &ObjectBodyReaction{}, nil
	default:
		return nil, errors.Errorf("unimplemented reaction type: '%d'", t)
	}
}

func serialize(reaction core.Reaction, t Type) (io.Reader, error) {
	buff := &bytes.Buffer{}
	_, err := buff.Write([]byte{byte(t)})
	if err != nil {
		return nil, err
	}

	enc := gob.NewEncoder(buff)
	err = enc.Encode(reaction)
	return buff, err
}

// Deserialize returns a reaction.
func Deserialize(buff io.Reader) (core.Reaction, error) {
	b := make([]byte, 1)
	_, err := buff.Read(b)
	if err != nil {
		return nil, errors.New("too short input to deserialize an event reaction")
	}

	reaction, err := getEmptyReaction(Type(b[0]))
	if err != nil {
		return nil, err
	}
	enc := gob.NewDecoder(buff)
	err = enc.Decode(reaction)
	return reaction, err
}

func init() {
	gob.Register(&CommonReaction{})
	gob.Register(&ObjectBodyReaction{})
}
