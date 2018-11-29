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

package genesisdataprovider

import (
	"context"
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func mockContractRequesterWithError(t *testing.T) *testutils.ContractRequesterMock {
	contractRequesterMock := testutils.NewContractRequesterMock(t)
	contractRequesterMock.SendRequestFunc = func(p context.Context, p1 *core.RecordRef, p2 string, p3 []interface{}) (r core.Reply, r1 error) {
		return nil, errors.New("test reasons")
	}
	return contractRequesterMock
}

func mockContractRequester(t *testing.T, res core.Reply) *testutils.ContractRequesterMock {
	contractRequesterMock := testutils.NewContractRequesterMock(t)
	contractRequesterMock.SendRequestFunc = func(p context.Context, p1 *core.RecordRef, p2 string, p3 []interface{}) (r core.Reply, r1 error) {
		return res, nil
	}
	return contractRequesterMock
}

func mockCertificate(t *testing.T, rootDomainRef *core.RecordRef) *testutils.CertificateMock {
	certificateMock := testutils.NewCertificateMock(t)
	certificateMock.GetRootDomainReferenceFunc = func() (r *core.RecordRef) {
		return rootDomainRef
	}
	return certificateMock
}

func TestNew(t *testing.T) {
	contractRequester := mockContractRequester(t, nil)
	certificate := mockCertificate(t, nil)

	result, err := New()

	cm := &component.Manager{}
	cm.Inject(contractRequester, certificate, result)

	require.NoError(t, err)
	require.Equal(t, result.Certificate, certificate)
	require.Equal(t, result.ContractRequester, contractRequester)
}
