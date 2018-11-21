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

package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/insolar/insolar/core/utils"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

func (ar *Runner) infoHandler() func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, req *http.Request) {

		ctx, inslog := inslogger.WithTraceField(context.Background(), utils.RandTraceID())

		rootDomain := ar.GenesisDataProvider.GetRootDomain(ctx)
		if rootDomain == nil {
			inslog.Error("[ INFO ] Can't get rootDomain ref")
		}
		rootMember := ar.GenesisDataProvider.GetRootMember(ctx)
		if rootMember == nil {
			inslog.Error("[ INFO ] Can't get rootMember ref")
		}
		nodeDomain := ar.GenesisDataProvider.GetNodeDomain(ctx)
		if nodeDomain == nil {
			inslog.Error("[ INFO ] Can't get nodeDomain ref")
		}
		data, err := json.MarshalIndent(map[string]interface{}{
			"root_domain": rootDomain.String(),
			"root_member": rootMember.String(),
			"node_domain": nodeDomain.String(),
		}, "", "   ")
		if err != nil {
			inslog.Error(errors.Wrap(err, "[ INFO ] Can't marshal response"))
		}

		response.Header().Add("Content-Type", "application/json")
		_, err = response.Write(data)
		if err != nil {
			inslog.Error(errors.Wrap(err, "[ INFO ] Can't write response"))
		}
	}
}
