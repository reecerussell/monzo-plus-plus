package persistence

import (
	"context"
	"database/sql"
	"os"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/datamodel"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/model"

	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/repository"
)

type PermissionsRepository struct {
	db *sql.DB
}

func NewPermissionRepository() repository.PermissionsRepository {
	return new(PermissionsRepository)
}

// LoadCollections loads all permission/role records from the database and
// organisises them into a map[int][]string, in the format of map[permissionID]roleIDs.
func (pr *PermissionsRepository) LoadCollections() map[int][]string {
	if err := pr.openConnection(); err != nil {
		panic(err)
	}

	query := `SELECT 
					permission_id, role_id
				FROM
					role_permissions
				ORDER BY role_id;`

	ctx := context.Background()
	stmt, err := pr.db.PrepareContext(ctx, query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		panic(err)
	}

	collections := make(map[int][]string)

	for rows.Next() {
		var (
			permissionID int
			roleID       string
		)

		err = rows.Scan(&permissionID, &roleID)
		if err != nil {
			panic(err)
		}

		roles, ok := collections[permissionID]
		if !ok {
			collections[permissionID] = []string{roleID}
		} else {
			roles = append(roles, roleID)
			collections[permissionID] = roles
		}
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	return collections
}

func (pr *PermissionsRepository) Get(id int) (*model.Permission, errors.Error) {
	query := "SELECT id, name FROM permissions WHERE id = ?;"

	ctx := context.Background()
	stmt, err := pr.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.InternalError(err)
	}
	defer stmt.Close()

	var dm datamodel.Permission

	err = stmt.QueryRowContext(ctx).Scan(&dm.ID, &dm.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFound("permission not found")
		}

		return nil, errors.InternalError(err)
	}

	return model.PermissionFromDataModel(&dm), nil
}

func (pr *PermissionsRepository) openConnection() error {
	if pr.db == nil {
		db, err := sql.Open("mysql", os.Getenv("CONN_STRING"))
		if err != nil {
			return err
		}

		pr.db = db
	}

	ctx := context.Background()
	err := pr.db.PingContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
