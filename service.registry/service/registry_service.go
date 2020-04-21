package service

import (
	"context"
	"sync"

	"github.com/reecerussell/monzo-plus-plus/service.registry/proto"
)

func init() {
	DefaultRegistry = &RegistryService{
		mu:    sync.RWMutex{},
		hosts: make(map[string]string),
	}
}

// DefaultRegistry is the default registry used for the gRPC server.
var DefaultRegistry *RegistryService

// RegistryService is used to register plugins.
type RegistryService struct {
	mu    sync.RWMutex
	hosts map[string]string
}

// Register adds a plugin host to the regsitry.
func (rs *RegistryService) Register(ctx context.Context, in *proto.RegisterType) (*proto.EmptyResponse, error) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	rs.hosts[in.GetName()] = in.GetHost()

	return &proto.EmptyResponse{}, nil
}

// Unregister removes a plugin from the registry.
func (rs *RegistryService) Unregister(ctx context.Context, in *proto.UnregisterType) (*proto.EmptyResponse, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	if h, ok := rs.hosts[in.GetName()]; ok {
		delete(rs.hosts, h)
	}

	return &proto.EmptyResponse{}, nil
}

// GetHost returns the host address for a given plugin.
func (rs *RegistryService) GetHost(ctx context.Context, in *proto.GetHostType) (*proto.HostResponseType, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	h, ok := rs.hosts[in.GetName()]
	if ok {
		return &proto.HostResponseType{
			Host: h,
		}, nil
	}

	return &proto.HostResponseType{
		Host: "",
	}, nil
}
