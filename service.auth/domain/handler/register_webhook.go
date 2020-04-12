package handler

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/reecerussell/monzo-plus-plus/libraries/monzo"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/event"
)

// RegisterWebhook is an event handler for the event.RegisterWebhook event.
type RegisterWebhook struct{}

// Invoke is used to make a HTTP request to Monzo's API to register a
// webhook for the user's account.
func (*RegisterWebhook) Invoke(ctx context.Context, tx *sql.Tx, e interface{}) errors.Error {
	evt := e.(*event.RegisterWebhook)

	err := monzo.RegisterHook(evt.AccountID, evt.AccessToken)
	if err != nil {
		return errors.InternalError(fmt.Errorf("register webhook: %v", err))
	}

	return nil
}
