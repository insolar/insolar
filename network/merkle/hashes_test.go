package merkle

import (
	"context"
	"crypto"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/insolar/insolar/network"
	network2 "github.com/insolar/insolar/testutils/network"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	component "github.com/insolar/component-manager"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar/pulsartestutils"
	"github.com/insolar/insolar/testutils"
)

func (t *calculatorHashesSuite) TestGetPulseHash() {
	pulseEntry := &PulseEntry{Pulse: t.pulse}
	ph, _, err := t.calculator.GetPulseProof(pulseEntry)
	t.Assert().NoError(err)

	expectedHash, _ := hex.DecodeString(
		"bd18c009950389026c5c6f85c838b899d188ec0d667f77948aa72a49747c3ed31835b1bdbb8bd1d1de62846b5f308ae3eac5127c7d36d7d5464985004122cc90",
	)

	t.Assert().Equal(OriginHash(expectedHash), ph)
}

func (t *calculatorHashesSuite) TestGetGlobuleHash() {
	pulseEntry := &PulseEntry{Pulse: t.pulse}
	ph, pp, err := t.calculator.GetPulseProof(pulseEntry)
	t.Assert().NoError(err)

	prevCloudHash, _ := hex.DecodeString(
		"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	)

	globuleEntry := &GlobuleEntry{
		PulseEntry: pulseEntry,
		PulseHash:  ph,
		ProofSet: map[insolar.NetworkNode]*PulseProof{
			t.originProvider.GetOrigin(): pp,
		},
		PrevCloudHash: prevCloudHash,
		GlobuleID:     0,
	}
	gh, _, err := t.calculator.GetGlobuleProof(globuleEntry)
	t.Assert().NoError(err)

	expectedHash, _ := hex.DecodeString(
		"68cd36762548acd48795678c2e308978edd1ff74de2f5daf0511c1b52cf7a7bef44e09d5dd5806e99aa4ed4253aca88390e6b376e0c5f5a49ff48a8f9547e5c5",
	)

	t.Assert().Equal(OriginHash(expectedHash), gh)
}

func (t *calculatorHashesSuite) TestGetCloudHash() {
	pulseEntry := &PulseEntry{Pulse: t.pulse}
	ph, pp, err := t.calculator.GetPulseProof(pulseEntry)
	t.Assert().NoError(err)

	prevCloudHash, _ := hex.DecodeString(
		"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	)

	globuleEntry := &GlobuleEntry{
		PulseEntry: pulseEntry,
		PulseHash:  ph,
		ProofSet: map[insolar.NetworkNode]*PulseProof{
			t.originProvider.GetOrigin(): pp,
		},
		PrevCloudHash: prevCloudHash,
		GlobuleID:     0,
	}
	_, gp, err := t.calculator.GetGlobuleProof(globuleEntry)

	ch, _, err := t.calculator.GetCloudProof(&CloudEntry{
		ProofSet:      []*GlobuleProof{gp},
		PrevCloudHash: prevCloudHash,
	})

	t.Assert().NoError(err)

	expectedHash, _ := hex.DecodeString(
		"68cd36762548acd48795678c2e308978edd1ff74de2f5daf0511c1b52cf7a7bef44e09d5dd5806e99aa4ed4253aca88390e6b376e0c5f5a49ff48a8f9547e5c5",
	)

	fmt.Println(hex.EncodeToString(ch))

	t.Assert().Equal(OriginHash(expectedHash), ch)
}

type calculatorHashesSuite struct {
	suite.Suite

	pulse          *insolar.Pulse
	originProvider network.OriginProvider
	service        insolar.CryptographyService

	calculator Calculator
}

func TestCalculatorHashes(t *testing.T) {
	calculator := &calculator{}

	key, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	require.NotNil(t, key)

	service := testutils.NewCryptographyServiceMock(t)
	service.SignMock.Set(func(p []byte) (r *insolar.Signature, r1 error) {
		signature := insolar.SignatureFromBytes([]byte("signature"))
		return &signature, nil
	})
	service.GetPublicKeyMock.Set(func() (r crypto.PublicKey, r1 error) {
		return "key", nil
	})

	stater := staterMock{
		stateFunc: func() []byte {
			return []byte("state")
		},
	}
	scheme := platformpolicy.NewPlatformCryptographyScheme()
	op := network2.NewOriginProviderMock(t)
	op.GetOriginMock.Set(func() insolar.NetworkNode {
		return createOrigin()
	})

	th := testutils.NewTerminationHandlerMock(t)

	cm := component.NewManager(nil)
	cm.Inject(th, op, &stater, calculator, service, scheme)

	require.NotNil(t, calculator.Stater)
	require.NotNil(t, calculator.OriginProvider)
	require.NotNil(t, calculator.CryptographyService)
	require.NotNil(t, calculator.PlatformCryptographyScheme)

	err := cm.Init(context.Background())
	require.NoError(t, err)

	pulse := &insolar.Pulse{
		PulseNumber:     insolar.PulseNumber(1337),
		NextPulseNumber: insolar.PulseNumber(1347),
		Entropy:         pulsartestutils.MockEntropyGenerator{}.GenerateEntropy(),
	}

	s := &calculatorHashesSuite{
		Suite:          suite.Suite{},
		calculator:     calculator,
		pulse:          pulse,
		originProvider: op,
		service:        service,
	}
	suite.Run(t, s)
}
