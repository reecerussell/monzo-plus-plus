package permission

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/reecerussell/monzo-plus-plus/libraries/permission/proto"
	"google.golang.org/grpc"
)

// Permission keys.
const (
	PermissionCreateUser      = 1
	PermissionGetUser         = 2
	PermissionGetUserList     = 3
	PermissionGetPendingUsers = 4
	PermissionUpdateUser      = 5
	PermissionEnableUser      = 6
	PermissionDeleteUser      = 7
	PermissionRoleManager     = 8
	PermissionCreateRole      = 9
	PermissionGetRole         = 10
	PermissionGetRoleList     = 11
	PermissionUpdateRole      = 12
	PermissionDeleteRole      = 13
	PermissionPluginManager   = 14
)

// AuthRPCAddress is an environment variable for the auth rpc target.
var AuthRPCAddress = os.Getenv("AUTH_RPC_ADDRESS")

// Has returns whether the user, the access token belongs to, has the
// given permission.
func Has(accessToken string, permission int) bool {
	conn, err := grpc.Dial(AuthRPCAddress, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to dial auth rpc service: %v", err)
		return false
	}
	defer conn.Close()

	client := proto.NewPermissionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	data := &proto.PermissionData{
		AccessToken: accessToken,
		Permission:  int32(permission),
	}

	pErr, err := client.HasPermission(ctx, data)
	if err != nil {
		log.Printf("Failed to call client method: HasPermission: %v", err)
		return false
	}

	return pErr == nil
}
