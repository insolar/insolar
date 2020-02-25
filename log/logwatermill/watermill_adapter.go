// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package logwatermill

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/insolar/insolar/insolar"
)

type WatermillLogAdapter struct {
	log insolar.Logger
}

func NewWatermillLogAdapter(log insolar.Logger) *WatermillLogAdapter {
	return &WatermillLogAdapter{
		log: log.WithField("service", "watermill"),
	}
}

func (w *WatermillLogAdapter) addFields(fields watermill.LogFields) insolar.Logger {
	return w.log.WithFields(fields)
}

func (w *WatermillLogAdapter) event(fields watermill.LogFields, level insolar.LogLevel, args ...interface{}) {
	w.addFields(fields).Embeddable().EmbeddedEvent(level, args...)
}

func (w *WatermillLogAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	l := w.addFields(fields)
	return &WatermillLogAdapter{log: l}
}

func (w *WatermillLogAdapter) Error(msg string, err error, fields watermill.LogFields) {
	w.event(fields, insolar.ErrorLevel, msg, " | Error: "+err.Error())
}

func (w *WatermillLogAdapter) Info(msg string, fields watermill.LogFields) {
	w.event(fields, insolar.InfoLevel, msg)
}

func (w *WatermillLogAdapter) Debug(msg string, fields watermill.LogFields) {
	w.event(fields, insolar.DebugLevel, msg)
}

func (w *WatermillLogAdapter) Trace(msg string, fields watermill.LogFields) {
	// don't use w.Debug(), value of the "file=..." field would be incorrect
	// in the output
	w.event(fields, insolar.DebugLevel, msg)
}
