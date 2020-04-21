package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/model"
)

// UserService is used to provide extra functionality, such as validation,
// for the User domain.
type UserService struct {
	db *sql.DB
}

// NewUserService returns a new instance of UserService.
func NewUserService() *UserService {
	return new(UserService)
}

// ValidateUsername validates a user's username, to ensure it is unique
// and hasn't been taken. An error is returned if it has been taken, or
// if there was an error communicating with the database.
func (us *UserService) ValidateUsername(u *model.User) errors.Error {
	if openErr := us.openConnection(); openErr != nil {
		return openErr
	}

	query := "SELECT COUNT(*) FROM users WHERE username LIKE ? AND id != ?;"

	ctx := context.Background()
	stmt, err := us.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.InternalError(err)
	}
	defer stmt.Close()

	var count int

	dm := u.DataModel()
	err = stmt.QueryRowContext(ctx, dm.Username, dm.ID).Scan(&count)
	if err != nil {
		return errors.InternalError(err)
	}

	if count > 0 {
		return errors.BadRequest(fmt.Sprintf("the username '%s' is already taken", dm.Username))
	}

	return nil
}

func (us *UserService) openConnection() errors.Error {
	if us.db == nil {
		db, err := sql.Open("mysql", os.Getenv("CONN_STRING"))
		if err != nil {
			return errors.InternalError(err)
		}

		us.db = db
	}

	ctx := context.Background()
	err := us.db.PingContext(ctx)
	if err != nil {
		return errors.InternalError(err)
	}

	return nil
}
