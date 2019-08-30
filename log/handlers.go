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

package log

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"

	"github.com/insolar/insolar/insolar"
)

// LogHandler is an HTTP handler that changes the log level.
type LogHandler struct {
	mux           *chi.Mux
	logController insolar.LogController
}

func NewLogHandler(base string, logController insolar.LogController) *LogHandler {
	r := chi.NewRouter()
	lh := &LogHandler{
		mux:           r,
		logController: logController,
	}
	r.Use(recoverer)
	r.Use(middleware.Throttle(1))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route(base, func(r chi.Router) {
		r.Get("/", lh.getRules)
		r.Post("/", lh.setRule)
		r.Delete("/", lh.deleteRule)
	})

	return lh
}

type jsonErr struct {
	Error string
}

type RuleItem struct {
	Level  string `json:"level"`
	Prefix string `json:"prefix"`
}

func (lh *LogHandler) deleteRule(w http.ResponseWriter, r *http.Request) {
	prefixStr := r.URL.Query().Get("prefix")
	ok := lh.logController.Del(prefixStr)
	if !ok {
		render.Status(r, http.StatusNoContent)
		return
	}

	render.JSON(w, r, struct {
		Prefix string `json:"prefix"`
	}{Prefix: prefixStr})
}

func (lh *LogHandler) setRule(w http.ResponseWriter, r *http.Request) {
	levelStr := r.URL.Query().Get("level")
	level, err := insolar.ParseLevel(levelStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r,
			jsonErr{fmt.Sprintf("Invalid level '%v': %v\n", levelStr, err)})
		return
	}

	if level == insolar.NoLevel {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r,
			jsonErr{"level value should not be empty"})
		return
	}

	prefixStr := r.URL.Query().Get("prefix")
	lh.logController.Set(prefixStr, level)
	render.JSON(w, r, struct {
		Level  string `json:"level"`
		Prefix string `json:"prefix"`
	}{
		Level:  level.String(),
		Prefix: prefixStr,
	})
}

func (lh *LogHandler) getRules(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, struct {
		GlobalLevel string     `json:"global_level"`
		Items       []RuleItem `json:"items"`
	}{
		GlobalLevel: zerolog.GlobalLevel().String(),
		Items: func() (items []RuleItem) {
			for _, item := range lh.logController.List() {
				items = append(items, RuleItem{
					Level:  item.Level.String(),
					Prefix: item.Prefix,
				})
			}
			return
		}(),
	})
}

// ServeHTTP implements http.Handler.
func (lh *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if lh.logController == nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("log controller not initialized"))
		return
	}
	lh.mux.ServeHTTP(w, r)
}

// recoverer is a chi middleware. handles panics in handlers.
func recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Log Controller API Panic: %+v\n", rvr)
				debug.PrintStack()
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
