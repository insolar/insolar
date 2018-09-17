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
	assert.NoError(t, err)

	counter, err := m.AddCounter("test_total", "metrics", "just test counter")
	assert.NoError(t, err)

	gauge, err := m.AddGauge("temperature_celsius", "metrics", "just test gauge")
	assert.NoError(t, err)

	counter.Add(55)
	gauge.Set(100)

	response, err := http.Get("http://" + cfg.ListenAddress + "/metrics")
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	contentText := string(content)
	assert.NoError(t, err)

	assert.True(t, strings.Contains(contentText, "insolar_metrics_test_total 55"))
	assert.True(t, strings.Contains(contentText, "insolar_metrics_temperature_celsius 100"))

	assert.NoError(t, m.Stop())
}
