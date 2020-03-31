package handler

import (
	"context"
	"database/sql"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"

	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/event"
)

// DisablePluginForUser is an event handler for the event.DisablePluginForUser event.
type DisablePluginForUser struct{}

// Invoke is used to handle a event.DisablePluginForUser event.
func (*DisablePluginForUser) Invoke(ctx context.Context, tx *sql.Tx, e interface{}) errors.Error {
	evt := e.(*event.DisablePluginForUser)
	query := "CALL disable_plugin_for_user(?,?);"
	args := []interface{}{evt.PluginID, evt.UserID}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return errors.InternalError(err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return errors.InternalError(err)
	}

	if rows, _ := res.RowsAffected(); rows < 1 {
		return errors.BadRequest("plugin either doesn't exist or is already disabled")
	}

	return nil
}
