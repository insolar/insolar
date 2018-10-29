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

package instracer

import (
	"context"

	"github.com/ugorji/go/codec"
)

// MustDeserialize encode baggage entries frotom bytes, panics on error.
func MustSerialize(ctx context.Context) []byte {
	b, err := Serialize(ctx)
	if err != nil {
		panic(err)
	}
	return b
}

// Serialize encode baggage entries to bytes.
func Serialize(ctx context.Context) ([]byte, error) {
	entries := GetBaggage(ctx)
	ch := new(codec.CborHandle)
	var b []byte
	err := codec.NewEncoderBytes(&b, ch).Encode(entries)
	return b, err
}

// MustDeserialize decode baggage entries from bytes, panics on error.
func MustDeserialize(b []byte) []Entry {
	bag, err := Deserialize(b)
	if err != nil {
		panic(err)
	}
	return bag
}

// Deserialize decode baggage entries from bytes.
func Deserialize(b []byte) ([]Entry, error) {
	var bag []Entry
	ch := new(codec.CborHandle)
	err := codec.NewDecoderBytes(b, ch).Decode(&bag)
	return bag, err
}
