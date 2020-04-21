package rpc

import (
	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/service.auth/interface/rpc/proto"

	"google.golang.org/grpc"
)

func Build(ctn *di.Container) *bootstrap.RPCServer {
	s := grpc.NewServer()
	ps := NewPermissionService(ctn)

	proto.RegisterPermissionServiceServer(s, ps)

	return bootstrap.BuildRPCServer(s)
}
