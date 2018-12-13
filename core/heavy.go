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

package core

import (
	"context"
)

// HeavySync provides methods for sync on heavy node.
//go:generate minimock -i github.com/insolar/insolar/core.HeavySync -o ../testutils -s _mock.go
type HeavySync interface {
	Start(ctx context.Context, jet RecordID, pn PulseNumber) error
	Store(ctx context.Context, jet RecordID, pn PulseNumber, kvs []KV) error
	Stop(ctx context.Context, jet RecordID, pn PulseNumber) error
	Reset(ctx context.Context, jet RecordID, pn PulseNumber) error
}
