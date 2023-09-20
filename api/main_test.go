package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type MainAPISuite struct {
	suite.Suite
}

func (suite *MainAPISuite) TestNewApiRunnerNilConfig() {
	_, err := NewRunner(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, Options{})
	suite.Contains(err.Error(), "config is nil")
}

func (suite *MainAPISuite) TestNewApiRunnerNoRequiredParams() {
	cfg := configuration.APIRunner{}
	_, err := NewRunner(&cfg, nil, nil, nil, nil, nil, nil, nil, nil, nil, Options{})
	suite.Contains(err.Error(), "Address must not be empty")

	cfg.Address = "address:100"
	_, err = NewRunner(&cfg, nil, nil, nil, nil, nil, nil, nil, nil, nil, Options{})
	suite.Contains(err.Error(), "RPC must exist")

	cfg.RPC = "test"
	_, err = NewRunner(&cfg, nil, nil, nil, nil, nil, nil, nil, nil, nil, Options{})
	suite.Contains(err.Error(), "Missing openAPI spec file path")

	cfg.SwaggerPath = "testdata/api-exported.yaml"
	runner, err := NewRunner(&cfg, nil, nil, nil, nil, nil, nil, nil, nil, nil, Options{})
	suite.NoError(err)
	suite.NoError(runner.Stop(context.Background()))
}

func TestMainTestSuite(t *testing.T) {
	ctx, _ := inslogger.WithTraceField(context.Background(), "APItests")
	http.DefaultServeMux = new(http.ServeMux)
	cfg := configuration.NewAPIRunner(false)
	cfg.SwaggerPath = "testdata/api-exported.yaml"
	api, err := NewRunner(&cfg, nil, nil, nil, nil, nil, nil, nil, nil, nil, Options{})
	require.NoError(t, err, "new runner constructor")

	cm := certificate.NewCertificateManager(&certificate.Certificate{})
	api.CertificateManager = cm
	api.Start(ctx)

	suite.Run(t, new(MainAPISuite))

	api.Stop(ctx)
}
