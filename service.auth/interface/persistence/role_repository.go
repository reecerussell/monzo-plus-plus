package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/datamodel"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/repository"
)

type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository() repository.RoleRepository {
	return new(roleRepository)
}

func (rr *roleRepository) Get(id string) (*model.Role, errors.Error) {
	if openErr := rr.openConnection(); openErr != nil {
		return nil, openErr
	}

	query := "CALL get_role(?);"

	ctx := context.Background()
	stmt, err := rr.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.InternalError(err)
	}
	defer stmt.Close()

	var role datamodel.Role
	var perms []*datamodel.Permission

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, errors.InternalError(err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&role.ID, &role.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.NotFound("role not found")
			}

			return nil, errors.InternalError(err)
		}

		if rows.NextResultSet() {
			for rows.Next() {
				var dm datamodel.Permission

				err = rows.Scan(&dm.ID, &dm.Name)
				if err != nil {
					return nil, errors.InternalError(err)
				}

				perms = append(perms, &dm)
			}
		}
	}
	if err = rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFound("role not found")
		}

		return nil, errors.InternalError(err)
	}

	return model.RoleFromDataModel(&role, perms...), nil
}

func (rr *roleRepository) GetList(term string) ([]*model.Role, errors.Error) {
	if openErr := rr.openConnection(); openErr != nil {
		return nil, openErr
	}

	query := "SELECT id, `name` FROM roles WHERE `name` LIKE ?;"

	ctx := context.Background()
	stmt, err := rr.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.InternalError(err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, fmt.Sprintf("%%%s%%", term))
	if err != nil {
		return nil, errors.InternalError(err)
	}
	defer rows.Close()

	roles := []*model.Role{}

	for rows.Next() {
		dm, rErr := readRole(rows.Scan)
		if rErr != nil {
			return nil, rErr
		}

		roles = append(roles, model.RoleFromDataModel(dm))
	}
	if err := rows.Err(); err != nil {
		return nil, errors.InternalError(err)
	}

	return roles, nil
}

func readRole(s scannerFunc) (*datamodel.Role, errors.Error) {
	var dm datamodel.Role

	err := s(&dm.ID, &dm.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFound("role not found")
		}

		return nil, errors.InternalError(fmt.Errorf("read role: %v", err))
	}

	return &dm, nil
}

func (rr *roleRepository) EnsureExists(id string) errors.Error {
	if openErr := rr.openConnection(); openErr != nil {
		return openErr
	}

	query := "SELECT id FROM roles WHERE id = ?;"

	ctx := context.Background()
	stmt, err := rr.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.InternalError(err)
	}
	defer stmt.Close()

	var roleID string
	err = stmt.QueryRowContext(ctx, id).Scan(&roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NotFound("role not found")
		}

		return errors.InternalError(err)
	}

	return nil
}

func (rr *roleRepository) Insert(r *model.Role) errors.Error {
	if openErr := rr.openConnection(); openErr != nil {
		return openErr
	}

	query := "INSERT INTO roles (`id`,`name`) VALUES (?,?);"
	dm := r.DataModel()
	args := []interface{}{
		dm.ID,
		dm.Name,
	}

	ctx := context.Background()
	stmt, err := rr.db.PrepareContext(ctx, query)
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

func (rr *roleRepository) Update(r *model.Role) errors.Error {
	if openErr := rr.openConnection(); openErr != nil {
		return openErr
	}

	query := "UPDATE roles SET `name` = ? WHERE id = ?;"
	dm := r.DataModel()
	args := []interface{}{
		dm.Name,
		dm.ID,
	}

	ctx := context.Background()
	tx, err := rr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	if err != nil {
		return errors.InternalError(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return errors.InternalError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return errors.InternalError(err)
	}

	if err := r.DispatchEvents(ctx, tx); err != nil {
		return err
	}

	return nil
}

func (rr *roleRepository) Delete(id string) errors.Error {
	if openErr := rr.openConnection(); openErr != nil {
		return openErr
	}

	query := "DELETE FROM roles WHERE id = ?;"

	ctx := context.Background()
	stmt, err := rr.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.InternalError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return errors.InternalError(err)
	}

	return nil
}

func (rr *roleRepository) openConnection() errors.Error {
	if rr.db == nil {
		db, err := sql.Open("mysql", os.Getenv("CONN_STRING"))
		if err != nil {
			return errors.InternalError(err)
		}

		rr.db = db
	}

	ctx := context.Background()
	err := rr.db.PingContext(ctx)
	if err != nil {
		return errors.InternalError(err)
	}

	return nil
}
