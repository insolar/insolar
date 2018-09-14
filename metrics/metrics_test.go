package metrics

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/stretchr/testify/assert"
)

func TestMetrics_NewMetrics(t *testing.T) {
	cfg := configuration.NewMetrics()
	m, err := NewMetrics(cfg)
	assert.NoError(t, err)
	err = m.Start(nil)

	response, err := http.Get("http://" + cfg.ListenAddress + "/metrics")
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	contentText := string(content)
	assert.NoError(t, err)

	assert.True(t, strings.Contains(contentText, "insolar_network_messages_sent_total 0"))
	assert.NoError(t, m.Stop())
}
