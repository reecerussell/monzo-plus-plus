package handler

import (
	"context"
	"database/sql"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"

	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/event"
)

// RemovePermissionFromRole provides a handler for the RemovePermissionFromRole event.
type RemovePermissionFromRole struct {
}

// Invoke is used to handle a RemovePermissionFromRole event and remove a permission from a role.
func (*RemovePermissionFromRole) Invoke(ctx context.Context, tx *sql.Tx, e interface{}) errors.Error {
	evt := e.(*event.RemovePermissionFromRole)

	query := "CALL remove_permission_from_role(?,?)"
	args := []interface{}{evt.PermissionID, evt.RoleID}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return errors.InternalError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return errors.InternalError(err)
	}

	return nil
}
