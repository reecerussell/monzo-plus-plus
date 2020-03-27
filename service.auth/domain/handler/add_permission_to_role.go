package handler

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"

	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/event"
)

// AddPermissionToRole provides a handler for the AddPermissionToRole event.
type AddPermissionToRole struct {
}

// Invoke is used to handle a AddPermissionToRole event and assign a permission to a role.
func (*AddPermissionToRole) Invoke(ctx context.Context, tx *sql.Tx, e interface{}) errors.Error {
	evt := e.(*event.AddPermissionToRole)

	fmt.Println(e)
	fmt.Println(evt)

	query := "CALL add_permission_to_role(?,?)"
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
