// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

/*
Package metrics is a gateway for Prometheus monitoring system, it based on Prometheus golang client.
Package contains metrics collectors descriptions of entire project.
Component starts http server on http://0.0.0.0:8080/metrics by default(can be changed in configuration)

Example:

	// starts metrics server
	cfg := configuration.NewMetrics()
	m, _ := NewMetrics(cfg)
    m.Start(nil)

    // manipulate with network metrics
	NetworkMessageSentTotal.Inc()
	NetworkPacketSentTotal.WithLabelValues("ping").Add(55)

*/
package metrics
