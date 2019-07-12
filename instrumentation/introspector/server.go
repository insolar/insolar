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

// introspector provides grpc/rest introspection API endpoint on shared tcp port.
package introspector

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/introspector/introproto"
	"github.com/pkg/errors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server implements introspection API server.
type Server struct {
	addr   string
	pubSrv introproto.PublisherServer

	cancel context.CancelFunc
	fin    chan error
}

// NewServer creates configured introspection API server.
func NewServer(addr string, ps introproto.PublisherServer) *Server {
	return &Server{
		addr:   addr,
		pubSrv: ps,
	}
}

// Start starts introspection http/grpc endpoint.
func (s *Server) Start(ctx context.Context) error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return errors.Wrapf(err, "failed to start introspection server on %s", s.addr)
	}

	inslogger.FromContext(ctx).Infof("started introspection server on %s\n", l.Addr())
	return s.run(ctx, l)
}

// Stop stops introspection http/grpc endpoint.
func (s *Server) Stop(ctx context.Context) error {
	if s.cancel != nil {
		s.cancel()
		return <-s.fin
	}

	inslogger.FromContext(ctx).Warn("stop called for not started introspection server")
	return nil
}

func (s *Server) run(ctx context.Context, l net.Listener) error {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	s.fin = make(chan error)

	grpcServer := grpc.NewServer()
	introproto.RegisterPublisherServer(grpcServer, s.pubSrv)
	reflection.Register(grpcServer)

	mux := http.NewServeMux()

	customMarshaller := &runtime.JSONPb{
		EmitDefaults: true,
	}
	muxOpts := runtime.WithMarshalerOption(runtime.MIMEWildcard, customMarshaller)
	gwMux := runtime.NewServeMux(muxOpts)

	mux.Handle("/", gwMux)
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		_, _ = io.Copy(w, strings.NewReader(publisherSwagger))
	})

	dOpts := []grpc.DialOption{grpc.WithInsecure()}
	err := introproto.RegisterPublisherHandlerFromEndpoint(ctx, gwMux, l.Addr().String(), dOpts)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Handler: grpcHandlerFunc(grpcServer, mux),
	}

	inslog := inslogger.FromContext(ctx)
	go func() {
		<-ctx.Done()
		inslog.Debug("shutdown introspection server...")
		s.fin <- srv.Shutdown(ctx)
	}()
	go func() {
		err = srv.Serve(l)
		inslog.Debugf("introspection server stopped: %s", err)
	}()
	return nil
}

// grpcHandlerFunc provides routing to proper gRPC/gateway server.
// It works without TLS thanks to h2c lib.
func grpcHandlerFunc(grpcHandler http.Handler, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contentType := r.Header.Get("Content-Type")
			if r.ProtoMajor == 2 && strings.Contains(contentType, "application/grpc") {
				grpcHandler.ServeHTTP(w, r)
			} else {
				otherHandler.ServeHTTP(w, r)
			}
		}),
		&http2.Server{},
	)
}
