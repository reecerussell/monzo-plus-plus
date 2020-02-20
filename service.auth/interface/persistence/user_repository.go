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

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

type scannerFunc func(dsts ...interface{}) error

// userRepository is an implementation of the repository.UserRepository interface for MySQL.
type userRepository struct {
	db *sql.DB
}

// NewUserRepository returns a new instance of the UserRepository interface.
func NewUserRepository() repository.UserRepository {
	return new(userRepository)
}

func (ur *userRepository) Get(id string) (*model.User, errors.Error) {
	openErr := ur.openConnection()
	if openErr != nil {
		return nil, openErr
	}

	query := "CALL get_user_by_id(?);"

	ctx := context.Background()
	stmt, err := ur.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.InternalError(fmt.Errorf("prepare: %v", err))
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, errors.InternalError(fmt.Errorf("rows: %v", err))
	}
	defer rows.Close()

	var user *datamodel.User
	var roles []*datamodel.Role
	var token *datamodel.Token

	if rows.Next() {
		u, readErr := readUser(rows.Scan)
		if readErr != nil {
			return nil, readErr
		}

		user = u
	} else {
		return nil, errors.NotFound("user not found")
	}

	if err = rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFound("user not found")
		}

		return nil, errors.InternalError(err)
	}

	if rows.NextResultSet() {
		for rows.Next() {
			r, readErr := readRole(rows.Scan)
			if err != nil {
				return nil, readErr
			}

			roles = append(roles, r)
		}
	}

	if rows.NextResultSet() {
		if rows.Next() {
			t, readErr := readToken(rows.Scan)
			if err != nil {
				return nil, readErr
			}

			token = t
		}
	}

	return model.UserFromDataModel(user, roles, token), nil
}

func (ur *userRepository) GetList(term string) ([]*model.User, errors.Error) {
	return ur.readUsers("CALL get_users(?);", term)
}

func (ur *userRepository) GetPending(term string) ([]*model.User, errors.Error) {
	return ur.readUsers("CALL get_pending_users(?);", term)
}

func (ur *userRepository) readUsers(query string, args ...interface{}) ([]*model.User, errors.Error) {
	openErr := ur.openConnection()
	if openErr != nil {
		return nil, openErr
	}

	ctx := context.Background()
	stmt, err := ur.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.InternalError(fmt.Errorf("prepare: %v", err))
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.InternalError(fmt.Errorf("rows: %v", err))
	}
	defer rows.Close()

	var users []*model.User

	for rows.Next() {
		user, readErr := readUser(rows.Scan)
		if readErr != nil {
			return nil, readErr
		}

		users = append(users, model.UserFromDataModel(user, nil, nil))
	}
	if err := rows.Err(); err != nil {
		return nil, errors.InternalError(err)
	}

	return users, nil
}

func readUser(s scannerFunc) (*datamodel.User, errors.Error) {
	var dm datamodel.User

	err := s(
		&dm.ID,
		&dm.Username,
		&dm.PasswordHash,
		&dm.StateToken,
		&dm.Enabled,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFound("user not found")
		}

		return nil, errors.InternalError(fmt.Errorf("read user: %v", err))
	}

	return &dm, nil
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

func readToken(s scannerFunc) (*datamodel.Token, errors.Error) {
	var dm datamodel.Token

	err := s(
		&dm.UserID,
		&dm.AccessToken,
		&dm.RefreshToken,
		&dm.Expires,
		&dm.TokenType,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFound("token not found")
		}

		return nil, errors.InternalError(fmt.Errorf("read token: %v", err))
	}

	return &dm, nil
}

func (ur *userRepository) Insert(u *model.User) errors.Error {
	dm := u.DataModel()
	query := "INSERT INTO users (`id`, `username`, `password_hash`, `state_token`, `enabled`) VALUES (?,?,?,?,?);"
	args := []interface{}{
		dm.ID,
		dm.Username,
		dm.PasswordHash,
		dm.StateToken,
		dm.Enabled,
	}

	return ur.execute(query, args...)
}

func (ur *userRepository) Update(u *model.User) errors.Error {
	dm := u.DataModel()
	query := "UPDATE users SET username = ?, state_token = ?, password_hash = ?, enabled = ? WHERE id = ?;"
	args := []interface{}{
		dm.Username,
		dm.StateToken,
		dm.PasswordHash,
		dm.Enabled,
		dm.ID,
	}

	return ur.execute(query, args...)
}

func (ur *userRepository) Delete(id string) errors.Error {
	query := "DELETE FROM users WHERE id = ?;"

	return ur.execute(query, id)
}

func (ur *userRepository) execute(query string, args ...interface{}) errors.Error {
	openErr := ur.openConnection()
	if openErr != nil {
		return openErr
	}

	ctx := context.Background()
	tx, err := ur.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
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

	return nil
}

func (ur *userRepository) openConnection() errors.Error {
	if ur.db == nil {
		db, err := sql.Open("mysql", os.Getenv("CONN_STRING"))
		if err != nil {
			return errors.InternalError(err)
		}

		ur.db = db
	}

	ctx := context.Background()
	err := ur.db.PingContext(ctx)
	if err != nil {
		return errors.InternalError(err)
	}

	return nil
}
