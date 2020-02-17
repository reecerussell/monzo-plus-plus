package rpc

import (
	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/interface/rpc/proto"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/interface/rpc/service"
	"google.golang.org/grpc"
)

func Build(ctn *di.Container) *bootstrap.RPCServer {
	s := grpc.NewServer()
	bs := service.NewBudgetService(ctn)

	proto.RegisterBudgetServiceServer(s, bs)

	return bootstrap.BuildRPCServer(s)
}
