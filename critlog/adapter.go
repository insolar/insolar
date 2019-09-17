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

package critlog

import (
	"errors"
	"github.com/rs/zerolog"
	"io"
)

type Flusher interface {
	Flush() error
}

type Syncer interface {
	Sync() error
}

var _ zerolog.LevelWriter = &writerAdapter{}

func AsLevelWriter(w io.Writer) zerolog.LevelWriter {
	if lw, ok := w.(zerolog.LevelWriter); ok {
		return lw
	}
	return &writerAdapter{w}
}

type writerAdapter struct {
	w io.Writer
}

func (w *writerAdapter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	return w.w.Write(p)
}

func (w *writerAdapter) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

func (w *writerAdapter) Flush() error {
	if f, ok := w.w.(Flusher); ok {
		return f.Flush()
	}
	return errors.New("unsupported: Flush")
}

func (w *writerAdapter) Close() error {
	if f, ok := w.w.(io.Closer); ok {
		return f.Close()
	}
	return errors.New("unsupported: Close")
}

func (w *writerAdapter) Sync() error {
	if f, ok := w.w.(Syncer); ok {
		return f.Sync()
	}
	return errors.New("unsupported: Sync")
}
