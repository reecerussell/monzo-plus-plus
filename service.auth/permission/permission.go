package permission

import (
	"context"
	"sync"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"

	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/repository"

	"github.com/reecerussell/monzo-plus-plus/libraries/util"
)

// Permissions
const (
	PermissionCreateUser = 1
	PermissionGetUser    = 2
	PermissionGetList    = 3
	PermissionGetPending = 4
	PermissionUpdateUser = 5
	PermissionEnableUser = 6
	PermissionDeleteUser = 7
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
func Build(ctn *di.Container) {
	r := ctn.Resolve("perm_repo").(repository.PermissionsRepository)
	perms = r.LoadCollections()
}

// Has is used to determine if the current user has permssion to access a
// protected resource or operation.
//
// The standard HTTP code for a false value is http.StatusForbidden (403).
func Has(ctx context.Context, perm int) bool {
	mu.Lock()
	defer mu.Unlock()

	user := ctx.Value(userKey).(*model.User)
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
