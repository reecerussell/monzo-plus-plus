package handler

import (
	"context"
	"database/sql"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/event"
)

// AddUserToRole is an event handler for event.AddUserToRole.
type AddUserToRole struct{}

// Invoke is used to handle the event, and call a stored proceedure to
// assign a user to a role.
func (*AddUserToRole) Invoke(ctx context.Context, tx *sql.Tx, e interface{}) errors.Error {
	evt := e.(*event.AddUserToRole)

	query := "CALL add_user_to_role(?,?)"
	args := []interface{}{evt.UserID, evt.RoleID}

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
