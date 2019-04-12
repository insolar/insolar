//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package messagebus

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/platformpolicy/keys"
)

type parcelFactory struct {
	Cryptography insolar.CryptographyService `inject:""`
}

// NewParcelFactory returns new instance of parcelFactory
func NewParcelFactory() message.ParcelFactory {
	return &parcelFactory{}
}

func (pf *parcelFactory) Create(ctx context.Context, msg insolar.Message, sender insolar.Reference, token insolar.DelegationToken, currentPulse insolar.Pulse) (insolar.Parcel, error) {
	if msg == nil {
		return nil, errors.New("failed to signature a nil message")
	}

	serialized := message.ToBytes(msg)
	signature, err := pf.Cryptography.Sign(serialized)
	if err != nil {
		return nil, err
	}

	serviceData := message.ServiceData{
		LogTraceID:    inslogger.TraceID(ctx),
		LogLevel:      inslogger.GetLoggerLevel(ctx),
		TraceSpanData: instracer.MustSerialize(ctx),
	}

	return &message.Parcel{
		Msg:         msg,
		Signature:   signature.Bytes(),
		Sender:      sender,
		Token:       token,
		PulseNumber: currentPulse.PulseNumber,
		ServiceData: serviceData,
	}, nil
}

func (pf *parcelFactory) Validate(publicKey keys.PublicKey, parcel insolar.Parcel) error {
	ok := pf.Cryptography.Verify(publicKey, insolar.SignatureFromBytes(parcel.GetSign()), message.ToBytes(parcel.Message()))
	if !ok {
		return errors.New("parcel isn't valid")
	}
	return nil
}
