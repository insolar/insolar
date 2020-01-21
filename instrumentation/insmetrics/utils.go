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

package insmetrics

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"strings"
)

// SumMetricsValueByNamePrefix summarizes values of metrics with prefix.
// Reader expects to provide text stream in OpenMetrics format,
func SumMetricsValueByNamePrefix(r io.Reader, prefix string) float64 {
	var acc float64
	for _, s := range FindMetricsByNamePrefix(r, prefix) {
		vStr := ExtractValue(s)
		v, err := strconv.ParseFloat(vStr, 64)
		if err != nil {
			log.Printf("fail to parse value %v (line: %v)\n", vStr, s)
		}
		acc += v
	}
	return acc
}

// FindMetricsByNamePrefix finds all metrics with prefix.
// Reader expects to provide text stream in OpenMetrics format,
func FindMetricsByNamePrefix(r io.Reader, prefix string) []string {
	var result []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.HasPrefix(s, prefix) {
			result = append(result, s)
		}
	}
	return result
}

// ExtractValue extracts value of metric from line in OpenMetrics format.
func ExtractValue(s string) string {
	return s[strings.LastIndex(s, " ")+1:]
}
