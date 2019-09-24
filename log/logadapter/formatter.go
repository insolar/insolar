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
	"github.com/insolar/insolar/insolar"
	"reflect"
)

type FormatFunc func(...interface{}) string
type FormatfFunc func(string, ...interface{}) string

func GetDefaultLogMsgFormatter() MsgFormatConfig {
	return MsgFormatConfig{
		Sformat:  fmt.Sprint,
		Sformatf: fmt.Sprintf,
	}
}

type MsgFormatConfig struct {
	Sformat  FormatFunc
	Sformatf FormatfFunc
}

func (v MsgFormatConfig) TryLogObject(a ...interface{}) (insolar.LogObjectMarshaller, string) {
	if len(a) == 1 {
		switch v := a[0].(type) {
		case nil: // the most obvious case(s)
			break
		case string: // the most obvious case(s)
			return nil, v
		case insolar.LogObjectMarshaller:
			return v, ""
		default:
			vt := reflect.ValueOf(v)
			if vt.Kind() == reflect.Struct && len(vt.Type().Name()) == 0 {
				return GetInlineLogObjectMarshaller(vt), ""
			}
		}
	}
	return nil, v.Sformat(a...)
}

func GetInlineLogObjectMarshaller(v reflect.Value) insolar.LogObjectMarshaller {
	return defaultLogObjectMarshaller{v}
}

type defaultLogObjectMarshaller struct {
	v reflect.Value
}

func (v defaultLogObjectMarshaller) MarshalLogObject(insolar.LogObjectWriter) string {
	return fmt.Sprintf("%+v", v.v)
}
