package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/model"
)

// RoleService is used to provide extra functionality, such as validation,
// for the Role domain.
type RoleService struct {
	db *sql.DB
}

// NewRoleService returns a new instance of RoleService.
func NewRoleService() *RoleService {
	return new(RoleService)
}

// ValidateName validates a roles's name, to ensure it is unique
// and hasn't been taken. An error is returned if it has been taken, or
// if there was an error communicating with the database.
func (rs *RoleService) ValidateName(r *model.Role) errors.Error {
	if openErr := rs.openConnection(); openErr != nil {
		return openErr
	}

	query := "SELECT COUNT(*) FROM roles WHERE `name` LIKE ? AND id != ?;"

	ctx := context.Background()
	stmt, err := rs.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.InternalError(err)
	}
	defer stmt.Close()

	var count int

	dm := r.DataModel()
	err = stmt.QueryRowContext(ctx, dm.Name, dm.ID).Scan(&count)
	if err != nil {
		return errors.InternalError(err)
	}

	if count > 0 {
		return errors.BadRequest(fmt.Sprintf("the name '%s' is already taken", dm.Name))
	}

	return nil
}

func (rs *RoleService) openConnection() errors.Error {
	if rs.db == nil {
		db, err := sql.Open("mysql", os.Getenv("CONN_STRING"))
		if err != nil {
			return errors.InternalError(err)
		}

		rs.db = db
	}

	ctx := context.Background()
	err := rs.db.PingContext(ctx)
	if err != nil {
		return errors.InternalError(err)
	}

	return nil
}
