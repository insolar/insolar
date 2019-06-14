package bootstrap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/hostnetwork/host"
)

func TestAuthorizer_Authorize(t *testing.T) {

	var cert certificate.Certificate

	d := cert.GetDiscoveryNodes()[0]
	h, err := host.NewHostN(d.GetHost(), *d.GetNodeRef())
	assert.NoError(t, err)

	options := common.ConfigureOptions(configuration.NewConfiguration())

	a := NewAuthorizer(options)
	resp, err := a.Authorize(context.Background(), h, &cert)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
