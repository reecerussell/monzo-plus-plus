package handler

import (
	"context"
	"database/sql"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/event"
)

// DeleteToken is an event handler for the event.DeleteToken event.
type DeleteToken struct{}

// Invoke is used to handle the event and delete the token data.
func (*DeleteToken) Invoke(ctx context.Context, tx *sql.Tx, e interface{}) errors.Error {
	evt := e.(*event.DeleteToken)

	query := "CALL delete_user_token(?);"
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return errors.InternalError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, evt.UserID)
	if err != nil {
		return errors.InternalError(err)
	}

	return nil
}
