///
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
///

package logwatermill

import (
	"github.com/ThreeDotsLabs/watermill"

	"github.com/insolar/insolar/insolar"
)

func NewWatermillLogAdapter(log insolar.Logger) *WatermillLogAdapter {
	return &WatermillLogAdapter{
		log: log.WithField("service", "watermill"),
	}
}

type WatermillLogAdapter struct {
	log insolar.Logger
}

func (w *WatermillLogAdapter) event(fields watermill.LogFields, level insolar.LogLevel, msg string) {
	// don't use w.Debug() etc, value of the "file=..." field would be incorrect
	if fn := w.log.Embeddable().NewEventStruct(level); fn != nil {
		fn(insolar.LogObjectFields{Msg: msg, Fields: fields})
	}
}

func (w *WatermillLogAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	l := w.log.WithFields(fields)
	return &WatermillLogAdapter{log: l}
}

func (w *WatermillLogAdapter) Error(msg string, err error, fields watermill.LogFields) {
	w.event(fields, insolar.ErrorLevel, msg+" | Error: "+err.Error())
}

func (w *WatermillLogAdapter) Info(msg string, fields watermill.LogFields) {
	w.event(fields, insolar.InfoLevel, msg)
}

func (w *WatermillLogAdapter) Debug(msg string, fields watermill.LogFields) {
	w.event(fields, insolar.DebugLevel, msg)
}

func (w *WatermillLogAdapter) Trace(msg string, fields watermill.LogFields) {
	w.event(fields, insolar.DebugLevel, msg)
}
