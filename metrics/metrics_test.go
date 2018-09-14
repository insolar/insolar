package metrics

import (
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/stretchr/testify/assert"
)

func TestMetrics_NewMetrics(t *testing.T) {

	m, err := NewMetrics(configuration.NewMetrics())
	assert.NoError(t, err)
	err = m.Start(nil)
	assert.NoError(t, err)
	m.Stop()

}
