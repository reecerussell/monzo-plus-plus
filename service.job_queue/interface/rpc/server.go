package rpc

import (
	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/service.job_queue/interface/rpc/proto"
	"github.com/reecerussell/monzo-plus-plus/service.job_queue/processing"
	"google.golang.org/grpc"
)

// Build returns a new RPCServer
func Build(jp processing.JobProcessor) *bootstrap.RPCServer {
	server := grpc.NewServer()
	service := NewJobsService(jp)

	proto.RegisterJobsServiceServer(server, service)

	return bootstrap.BuildRPCServer(server)
}
