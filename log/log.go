//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package log

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log/logadapter"
	"github.com/insolar/insolar/log/zlogadapter"
)

// NewLog creates logger instance with particular configuration
func NewLog(cfg configuration.Log) (insolar.Logger, error) {
	return NewLogExt(cfg, 0)
}

// NewLog creates logger instance with particular configuration
func NewLogExt(cfg configuration.Log, skipFrameBaselineAdjustment int8) (insolar.Logger, error) {

	defaults := insolar.DefaultLoggerSettings()
	pCfg, err := insolar.ParseLogConfigWithDefaults(cfg, defaults)

	if err == nil {
		var logger insolar.Logger

		pCfg.SkipFrameBaselineAdjustment = skipFrameBaselineAdjustment

		msgFmt := logadapter.GetDefaultLogMsgFormatter()

		switch strings.ToLower(cfg.Adapter) {
		case "zerolog":
			logger, err = zlogadapter.NewZerologAdapter(pCfg, msgFmt)
		default:
			err = errors.New("unknown adapter")
		}

		if err == nil {
			if logger != nil {
				return logger, nil
			}
			return nil, errors.New("logger was not initialized")
		}
	}
	return nil, errors.Wrap(err, "invalid logger config")
}
