package heavy

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/ledger/heavy/exporter"
)

type Limiter interface {
	Allow() bool
}

type noLimit struct {
}

func (l *noLimit) Allow() bool {
	return true
}

func NewNoLimit(_ int) *noLimit {
	return &noLimit{}
}

type ServerLimiter struct {
	inbound  *Limiters
	outbound *Limiters
}

func NewServerLimiter(config configuration.RateLimit) *ServerLimiter {
	return &ServerLimiter{
		inbound:  NewLimiters(config.In),
		outbound: NewLimiters(config.Out),
	}
}

type Limiters struct {
	config            configuration.Limits
	globalLimiter     Limiter
	perClientLimiters map[string]map[string]Limiter
	mutex             *sync.RWMutex
}

func NewLimiters(config configuration.Limits) *Limiters {
	// here we will use a suitable implementation of limiter with the RPS value from the config.Global
	limiter := NewNoLimit(config.Global)
	return &Limiters{
		config:            config,
		globalLimiter:     limiter,
		perClientLimiters: make(map[string]map[string]Limiter),
		mutex:             &sync.RWMutex{},
	}
}

func (l *ServerLimiter) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		method := info.FullMethod
		limiters := l.inbound
		if limiters.GlobalLimit() || limiters.PerClientLimit(ctx, method) {
			return nil, status.Errorf(codes.ResourceExhausted, "method: %s, %s", method, exporter.RateLimitExceededMsg)
		}
		return handler(ctx, req)
	}
}

func (l *ServerLimiter) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		method := info.FullMethod
		limiters := l.inbound
		if limiters.GlobalLimit() || limiters.PerClientLimit(stream.Context(), method) {
			return status.Errorf(codes.ResourceExhausted, "method: %s, %s", method, exporter.RateLimitExceededMsg)
		}
		return handler(srv, l.LimitStream(stream, method))
	}
}

func (l *ServerLimiter) LimitStream(stream grpc.ServerStream, method string) grpc.ServerStream {
	return &limitedServerStream{
		ServerStream: stream,
		outbound:     l.outbound,
		method:       method,
	}
}

func (l *Limiters) GlobalLimit() bool {
	if l.globalLimiter == nil {
		return false
	}
	return !l.globalLimiter.Allow()
}

func (l *Limiters) PerClientLimit(ctx context.Context, method string) bool {
	md, _ := metadata.FromIncomingContext(ctx)
	client := "unknown"
	if _, isContain := md[exporter.ObsID]; isContain {
		client = md.Get(exporter.ObsID)[0]
	}
	l.mutex.RLock()
	limiter := l.perClientLimiters[method][client]
	l.mutex.RUnlock()
	if limiter == nil {
		// here we will use a suitable implementation of limiter with value l.config.PerClient.Limit(method)
		rps := l.config.PerClient.Limit(method)
		limiter = NewNoLimit(rps)
		l.mutex.Lock()
		l.perClientLimiters[method] = map[string]Limiter{client: limiter}
		l.mutex.Unlock()
	}
	return !limiter.Allow()
}

type limitedServerStream struct {
	grpc.ServerStream
	outbound *Limiters
	method   string
}

func (s *limitedServerStream) SendMsg(m interface{}) error {
	limiters := s.outbound
	if limiters.GlobalLimit() || limiters.PerClientLimit(s.Context(), s.method) {
		return status.Errorf(codes.ResourceExhausted, "method: %s, %s", s.method, exporter.RateLimitExceededMsg)
	}
	return s.ServerStream.SendMsg(m)
}
