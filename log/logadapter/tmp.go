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

package logadapter

import (
	"fmt"
	"reflect"

	"github.com/insolar/insolar/log/logcommon"
)

var prepareObjTypes = []struct { // MUST be ordered, not map
	t  reflect.Type
	fn func(interface{}) (interface{}, bool)
}{
	{reflect.TypeOf((*logcommon.LogStringer)(nil)).Elem(), func(value interface{}) (interface{}, bool) {
		if vv, ok := value.(logcommon.LogStringer); ok {
			return vv.LogString(), true
		}
		return value, false
	}},
	{reflect.TypeOf((*fmt.Stringer)(nil)).Elem(), func(value interface{}) (interface{}, bool) {
		if vv, ok := value.(fmt.Stringer); ok {
			return vv.String(), true
		}
		return value, false
	}},
}

func defaultStrValuePrepare(iv interface{}) (string, bool) {
	switch vv := iv.(type) {
	case string:
		return vv, true
	case *string:
		if vv == nil {
			return "", false
		}
		return *vv, true
	case func() string:
		return vv(), true
	}
	if vv, ok := tryDefaultValuePrepare(iv); ok {
		return vv.(string), true
	}
	return "", false
}

func tryDefaultValuePrepare(iv interface{}) (interface{}, bool) {
	for _, f := range prepareObjTypes {
		if vv, ok := f.fn(iv); ok {
			return vv, true
		}
	}
	return iv, false
}
