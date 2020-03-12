package bootstrap

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap/proto"
	"google.golang.org/grpc"
)

// RegistryRPCHost is an environment variable for the registry host address.
var RegistryRPCHost = os.Getenv("REGISTRY_HOST")

// Register registers a plugin with its name and host address.
func Register(name, host string) error {
	conn, err := grpc.Dial(RegistryRPCHost, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("dial: %v", err)
	}
	defer conn.Close()

	client := proto.NewRegistryServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	data := &proto.RegisterType{
		Name: name,
		Host: host,
	}

	_, err = client.Register(ctx, data)
	if err != nil {
		return fmt.Errorf("reg: %v", err)
	}

	return nil
}

// Unregister unregisters a plugin.
func Unregister(name string) error {
	conn, err := grpc.Dial(RegistryRPCHost, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("dial: %v", err)
	}
	defer conn.Close()

	client := proto.NewRegistryServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	data := &proto.UnregisterType{
		Name: name,
	}

	_, err = client.Unregister(ctx, data)
	if err != nil {
		return fmt.Errorf("reg: %v", err)
	}

	return nil
}

// GetHost returns a plugins host address.
func GetHost(name string) (string, error) {
	conn, err := grpc.Dial(RegistryRPCHost, grpc.WithInsecure())
	if err != nil {
		return "", fmt.Errorf("dial: %v", err)
	}
	defer conn.Close()

	client := proto.NewRegistryServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	data := &proto.GetHostType{
		Name: name,
	}

	res, err := client.GetHost(ctx, data)
	if err != nil {
		return "", fmt.Errorf("reg: %v", err)
	}

	return res.GetHost(), nil
}
