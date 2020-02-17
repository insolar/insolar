// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

type KeeperRsp struct {
	Available bool `json:"available"`
}

type PromRsp struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Installation string `json:"installation"`
				Instance     string `json:"instance"`
				Job          string `json:"job"`
				Role         string `json:"role"`
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`

	ErrorType string `json:"errorType"`
	Error     string `json:"error"`
}
