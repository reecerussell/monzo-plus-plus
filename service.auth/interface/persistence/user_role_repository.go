package persistence

import (
	"context"
	"database/sql"
	"os"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/repository"
)

type userRoleRepository struct {
	db *sql.DB
}

func NewUserRoleRepository() repository.UserRoleRepository {
	return new(userRoleRepository)
}

func (urr *userRoleRepository) Insert(userID, roleID string) errors.Error {
	return urr.execute("CALL add_user_to_role(?,?);", userID, roleID)
}

func (urr *userRoleRepository) Delete(userID, roleID string) errors.Error {
	return urr.execute("CALL remove_user_from_role(?,?);", userID, roleID)
}

func (urr *userRoleRepository) openConnection() errors.Error {
	if urr.db == nil {
		db, err := sql.Open("mysql", os.Getenv("CONN_STRING"))
		if err != nil {
			return errors.InternalError(err)
		}

		urr.db = db
	}

	ctx := context.Background()
	err := urr.db.PingContext(ctx)
	if err != nil {
		return errors.InternalError(err)
	}

	return nil
}

func (urr *userRoleRepository) execute(query string, args ...interface{}) errors.Error {
	if openErr := urr.openConnection(); openErr != nil {
		return openErr
	}

	ctx := context.Background()
	stmt, err := urr.db.PrepareContext(ctx, query)
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
