package rpc

import (
	"context"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/interface/rpc/proto"
	"github.com/reecerussell/monzo-plus-plus/service.auth/permission"
	"github.com/reecerussell/monzo-plus-plus/service.auth/registry"
	"github.com/reecerussell/monzo-plus-plus/service.auth/usecase"
)

type PermissionService struct {
	userAuthUsecase usecase.UserAuthUsecase
}

func NewPermissionService(ctn *di.Container) *PermissionService {
	uu := ctn.Resolve(registry.ServiceUserAuthUsecase).(usecase.UserAuthUsecase)

	return &PermissionService{uu}
}

func (ps *PermissionService) HasPermission(ctx context.Context, in *proto.PermissionData) (*proto.Error, error) {
	ctx, err := ps.userAuthUsecase.WithUser(ctx, in.GetAccessToken())
	if err != nil {
		return &proto.Error{
			StatusCode: int32(err.ErrorCode()),
			Message:    err.Text(),
		}, nil
	}

	if !permission.Has(ctx, int(in.GetPermission())) {
		err := errors.Forbidden()
		return &proto.Error{
			StatusCode: int32(err.ErrorCode()),
			Message:    err.Text(),
		}, nil
	}

	return &proto.Error{
		StatusCode: 200,
		Message:    "",
	}, nil
}
