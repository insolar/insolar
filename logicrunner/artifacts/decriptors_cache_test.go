// Copyright 2020 Insolar Network Ltd.
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

package artifacts

import (
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func TestNewDescriptorsCache(t *testing.T) {
	dc := NewDescriptorsCache()
	require.NotNil(t, dc)
}

func Test_descriptorsCache(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	ctx := inslogger.TestContext(t)
	protoRef := gen.Reference()
	codeRef := gen.Reference()

	type fields struct {
		Client     Client
		codeCache  cache
		protoCache cache
	}
	tests := []struct {
		name   string
		fields fields
		obj    ObjectDescriptor
		proto  ObjectDescriptor
		code   CodeDescriptor
		err    bool
	}{
		{
			name: "success",
			obj:  NewObjectDescriptorMock(mc).PrototypeMock.Return(&protoRef, nil),
			fields: fields{
				protoCache: NewCacheMock(mc).getMock.Return(
					NewPrototypeDescriptorMock(mc).
						CodeMock.Return(&codeRef).
						HeadRefMock.Return(&protoRef),
					nil,
				),
				codeCache: NewCacheMock(mc).getMock.Return(
					NewCodeDescriptorMock(mc).RefMock.Return(&codeRef), nil,
				),
			},
		},
		{
			name: "objDesc.Prototype fails -> error",
			err:  true,
			obj:  NewObjectDescriptorMock(mc).PrototypeMock.Return(nil, errors.New("has no prototype")),
		},
		{
			name: "no such prototype -> error",
			err:  true,
			obj:  NewObjectDescriptorMock(mc).PrototypeMock.Return(&protoRef, nil),
			fields: fields{
				protoCache: NewCacheMock(mc).getMock.Return(
					nil, errors.New("no proto"),
				),
			},
		},
		{
			name: "is not prototype -> error",
			err:  true,
			obj:  NewObjectDescriptorMock(mc).PrototypeMock.Return(nil, nil),
			fields: fields{
				protoCache: NewCacheMock(mc),
			},
		},
		{
			name: "bad code reference -> error",
			err:  true,
			obj:  NewObjectDescriptorMock(mc).PrototypeMock.Return(&protoRef, nil),
			fields: fields{
				protoCache: NewCacheMock(mc).getMock.Return(
					NewPrototypeDescriptorMock(mc).
						CodeMock.Return(&codeRef),
					nil,
				),
				codeCache: NewCacheMock(mc).getMock.Return(
					nil, errors.New("no code"),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &descriptorsCache{
				Client:     tt.fields.Client,
				codeCache:  tt.fields.codeCache,
				protoCache: tt.fields.protoCache,
			}
			pd, cd, err := c.ByObjectDescriptor(ctx, tt.obj)
			if tt.err {
				require.Error(t, err)
				require.Nil(t, pd)
				require.Nil(t, cd)
			} else {
				require.NoError(t, err)
				require.Equal(t, &protoRef, pd.HeadRef())
				require.Equal(t, &codeRef, cd.Ref())
			}
		})
	}
}
