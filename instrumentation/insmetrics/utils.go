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
