/*
 *
 *  *    Copyright 2018 Insolar
 *  *
 *  *    Licensed under the Apache License, Version 2.0 (the "License");
 *  *    you may not use this file except in compliance with the License.
 *  *    You may obtain a copy of the License at
 *  *
 *  *        http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  *    Unless required by applicable law or agreed to in writing, software
 *  *    distributed under the License is distributed on an "AS IS" BASIS,
 *  *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  *    See the License for the specific language governing permissions and
 *  *    limitations under the License.
 *
 */

package logicrunner

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/ugorji/go/codec"
	"golang.org/x/crypto/sha3"
)

type trf = core.RecordRef

type CaseRecordType int

// Types of records
const (
	CaseRecordTypeUnexistent CaseRecordType = iota
	CaseRecordTypeMethodCall
	CaseRecordTypeConstructorCall
	CaseRecordTypeRouteCall
	CaseRecordTypeSaveAsChild
	CaseRecordTypeGetObjChildren
	CaseRecordTypeSaveAsDelegate
	CaseRecordTypeGetDelegate
)

// CaseRecord is one record of validateable object calling history
type CaseRecord struct {
	Type   CaseRecordType
	ReqSig []byte
	Resp   rpctypes.UpRespIface
}

// CaseBinder is a whole result of executor efforts on every object it seen on this pulse
type CaseBind struct {
	P core.Pulse           // pulse info for this bind
	R map[trf][]CaseRecord // ordered cases for each object
}

func HashInterface(in interface{}) []byte {
	s := []byte{}
	ch := new(codec.CborHandle)
	err := codec.NewEncoderBytes(&s, ch).Encode(in)
	if err != nil {
		panic("Can't marshal: " + err.Error())
	}
	sh := sha3.New224()
	return sh.Sum(s)
}
