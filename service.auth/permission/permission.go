package permission

import (
	"context"
	"log"
	"sync"

	"github.com/reecerussell/monzo-plus-plus/libraries/util"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/repository"
)

// Permissions
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

var (
	mu      = &sync.RWMutex{}
	userKey = util.ContextKey("user")
	perms   map[int][]string
)

// Build is used to initialise the permissions service, by loading the permissions
// from the database.
//
// If an error occurs, it will panic.
func Build(repo repository.PermissionsRepository) {
	perms = repo.LoadCollections()

	for p, roles := range perms {
		log.Printf("Roles with permission: %d\n", p)

		for _, role := range roles {
			log.Printf("    - %s\n", role)
		}
	}
}

// Has is used to determine if the current user has permssion to access a
// protected resource or operation.
//
// The standard HTTP code for a false value is http.StatusForbidden (403).
func Has(ctx context.Context, perm int) bool {
	mu.Lock()
	defer mu.Unlock()

	val := ctx.Value(userKey)
	if val == nil {
		return false
	}

	user := val.(*model.User)
	if user == nil {
		return false
	}

	allowedRoles, ok := perms[perm]
	if !ok {
		return false
	}

	for _, roleID := range user.GetRoles() {
		for _, allowedRoleID := range allowedRoles {
			if roleID == allowedRoleID {
				return true
			}
		}
	}

	return false
}
