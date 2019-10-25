package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	jsonrpc "github.com/insolar/rpc/v2/json2"
	"github.com/pkg/errors"

	"github.com/insolar/rpc/v2"

	"github.com/insolar/insolar/application/api/requester"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/storage"
)

type ApiService struct {
	PulseAccessor storage.PulseAccessor `inject:""`
	NetworkStatus insolar.NetworkStatus `inject:""`

	handler   http.Handler
	server    *http.Server
	rpcServer *rpc.Server
	cfg       configuration.APIRunner
}

func NewApiService(cfg configuration.APIRunner) *ApiService {
	return &ApiService{cfg: cfg}
}

func (s *ApiService) Start(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	logger.Warn("=========================  Starting ApiService ")

	addrStr := fmt.Sprint(s.cfg.Address)
	s.server = &http.Server{Addr: addrStr}
	s.rpcServer = rpc.NewServer()

	s.rpcServer.RegisterCodec(jsonrpc.NewCodec(), "application/json")

	err := s.rpcServer.RegisterService(&NodeService{s}, "node")
	if err != nil {
		return errors.Wrap(err, "[ registerServices ] Can't RegisterService: node")
	}

	router := http.NewServeMux()
	s.server.Handler = router

	router.Handle(s.cfg.RPC, s.rpcServer)
	s.handler = router

	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return errors.Wrap(err, "Can't start listening")
	}
	go func() {
		err := s.server.Serve(listener)
		if err != nil {
			logger.Error("Http server: ListenAndServe() error: ", err)

		}

	}()
	logger.Warnf("=========================  ApiService started on %s", s.server.Addr)
	return nil
}

func (s *ApiService) Stop(ctx context.Context) error {
	const timeOut = 5

	inslogger.FromContext(ctx).Infof("Shutting down server gracefully ...(waiting for %d seconds)", timeOut)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeOut)*time.Second)
	defer cancel()
	err := s.server.Shutdown(ctxWithTimeout)
	if err != nil {
		return errors.Wrap(err, "Can't gracefully stop API server")
	}

	return nil
}

type NodeService struct {
	runner *ApiService
}

func (s *NodeService) GetStatus(r *http.Request, args *interface{}, requestBody *rpc.RequestBody, reply *requester.StatusResponse) error {
	traceID := utils.RandTraceID()
	ctx, inslog := inslogger.WithTraceField(context.Background(), traceID)

	inslog.Infof("[ NodeService.GetStatus ] Incoming request: %s", r.RequestURI)
	// if !s.runner.cfg.IsAdmin {
	// 	return errors.New("method not allowed")
	// }
	statusReply := s.runner.NetworkStatus.GetNetworkStatus()

	reply.NetworkState = statusReply.NetworkState.String()
	reply.ActiveListSize = statusReply.ActiveListSize
	reply.WorkingListSize = statusReply.WorkingListSize

	nodes := make([]requester.Node, reply.ActiveListSize)
	for i, node := range statusReply.Nodes {
		nodes[i] = requester.Node{
			Reference: node.ID().String(),
			Role:      node.Role().String(),
			IsWorking: node.GetPower() > 0,
			ID:        uint32(node.ShortID()),
		}
	}
	reply.Nodes = nodes

	reply.Origin = requester.Node{
		Reference: statusReply.Origin.ID().String(),
		Role:      statusReply.Origin.Role().String(),
		IsWorking: statusReply.Origin.GetPower() > 0,
		ID:        uint32(statusReply.Origin.ShortID()),
	}

	reply.NetworkPulseNumber = uint32(statusReply.Pulse.PulseNumber)

	p, err := s.runner.PulseAccessor.GetLatestPulse(ctx)
	if err != nil {
		p = *insolar.GenesisPulse
	}
	reply.PulseNumber = uint32(p.PulseNumber)

	reply.Entropy = statusReply.Pulse.Entropy[:]
	reply.Version = statusReply.Version
	reply.StartTime = statusReply.StartTime
	reply.Timestamp = statusReply.Timestamp

	return nil
}
