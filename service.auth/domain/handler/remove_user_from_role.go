package handler

import (
	"context"
	"database/sql"

	"github.com/reecerussell/monzo-plus-plus/libraries/domain"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/event"
)

// RemoveUserFromRole is an event handler for event.RemoveUserFromRole.
type RemoveUserFromRole struct{}

// Invoke is used to handle the event, and call a stored proceedure to
// unassign a user from a role.
func (*RemoveUserFromRole) Invoke(ctx context.Context, tx *sql.Tx, e domain.Event) errors.Error {
	evt := e.(*event.RemoveUserFromRole)

	query := "CALL remove_user_from_role(?,?)"
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
