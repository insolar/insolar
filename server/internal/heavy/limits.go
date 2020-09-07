package heavy

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	sw "github.com/insolar/ratelimiter/slidingwindow"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/ledger/heavy/exporter"
)

type limiter interface {
	Allow() bool
}

type serverLimiters struct {
	inbound  *limiters
	outbound *limiters
}

func newServerLimiters(config configuration.RateLimit) *serverLimiters {
	return &serverLimiters{
		inbound:  newLimiters(config.In),
		outbound: newLimiters(config.Out),
	}
}

type limiters struct {
	config            configuration.Limits
	globalLimiter     limiter
	perClientLimiters map[string]map[string]limiter
	mutex             *sync.RWMutex
}

func newLimiters(config configuration.Limits) *limiters {
	gl, _ := sw.NewLimiter(time.Second, int64(config.Global), func() (sw.Window, sw.StopFunc) {
		return sw.NewLocalWindow()
	})
	return &limiters{
		config:            config,
		globalLimiter:     gl,
		perClientLimiters: make(map[string]map[string]limiter),
		mutex:             &sync.RWMutex{},
	}
}

func (l *serverLimiters) unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		method := info.FullMethod
		limiters := l.inbound
		if limiters.isGlobalLimitExceeded() || limiters.isClientLimitExceeded(ctx, method) {
			return nil, status.Errorf(codes.ResourceExhausted, "method: %s, %s", method, exporter.RateLimitExceededMsg)
		}
		return handler(ctx, req)
	}
}

func (l *serverLimiters) streamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		method := info.FullMethod
		limiters := l.inbound
		if limiters.isGlobalLimitExceeded() || limiters.isClientLimitExceeded(stream.Context(), method) {
			return status.Errorf(codes.ResourceExhausted, "method: %s, %s", method, exporter.RateLimitExceededMsg)
		}
		return handler(srv, l.limitStream(stream, method))
	}
}

func (l *serverLimiters) limitStream(stream grpc.ServerStream, method string) grpc.ServerStream {
	return &limitedServerStream{
		ServerStream: stream,
		outbound:     l.outbound,
		method:       method,
	}
}

func (l *limiters) isGlobalLimitExceeded() bool {
	if l.globalLimiter == nil {
		return false
	}
	return !l.globalLimiter.Allow()
}

func (l *limiters) isClientLimitExceeded(ctx context.Context, method string) bool {
	md, _ := metadata.FromIncomingContext(ctx)
	client := "unknown"
	if _, isContain := md[exporter.ObsID]; isContain {
		client = md.Get(exporter.ObsID)[0]
	}

	var cl limiter
	func() {
		l.mutex.RLock()
		defer l.mutex.RUnlock()
		cl = l.perClientLimiters[method][client]
	}()

	if cl == nil {
		rps := l.config.PerClient.Limit(method)
		cl, _ = sw.NewLimiter(time.Second, int64(rps), func() (sw.Window, sw.StopFunc) {
			return sw.NewLocalWindow()
		})
		func() {
			l.mutex.Lock()
			defer l.mutex.Unlock()
			lm, isContain := l.perClientLimiters[method]
			if !isContain {
				lm = make(map[string]limiter)
				l.perClientLimiters[method] = lm
			}
			lm[client] = cl
		}()
	}

	return !cl.Allow()
}

type limitedServerStream struct {
	grpc.ServerStream
	outbound *limiters
	method   string
}

func (s *limitedServerStream) SendMsg(m interface{}) error {
	limiters := s.outbound
	if limiters.isGlobalLimitExceeded() || limiters.isClientLimitExceeded(s.Context(), s.method) {
		return status.Errorf(codes.ResourceExhausted, "method: %s, %s", s.method, exporter.RateLimitExceededMsg)
	}
	return s.ServerStream.SendMsg(m)
}
