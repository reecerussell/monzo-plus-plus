package persistence

import (
	"context"
	"database/sql"
	"os"

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
