//
// Copyright 2019 Insolar Technologies GmbH
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
//

package logadapter

import (
	"fmt"
	"reflect"
)

type LogStringer interface {
	LogString() string
}

type FormatFunc func(...interface{}) string
type FormatfFunc func(string, ...interface{}) string

type MsgFormatConfig struct {
	Sformat  FormatFunc
	Sformatf FormatfFunc
}

func GetDefaultLogMsgFormatter() MsgFormatConfig {
	return MsgFormatConfig{
		Sformat:  SpecialSprint,
		Sformatf: fmt.Sprintf,
	}
}

func SpecialSprint(a ...interface{}) string {
	if len(a) == 1 {
		switch v := a[0].(type) {
		case nil: // the most obvious case(s)
			break
		case string: // the most obvious case(s)
			return v
		case LogStringer:
			return v.LogString()
		default:
			vt := reflect.TypeOf(v)
			if vt.Kind() == reflect.Struct && len(vt.Name()) == 0 {
				return fmt.Sprintf("%+v", v)
			}
		}
	}

	return fmt.Sprint(a...)
}
