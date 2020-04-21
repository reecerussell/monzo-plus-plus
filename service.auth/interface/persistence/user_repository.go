package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/reecerussell/monzo-plus-plus/libraries/database"
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
	db  *database.DB
	sql *sql.DB
}

// NewUserRepository returns a new instance of the UserRepository interface.
func NewUserRepository() repository.UserRepository {
	return &userRepository{
		db: database.New(),
	}
}

func (ur *userRepository) Get(id string) (*model.User, errors.Error) {
	openErr := ur.openConnection()
	if openErr != nil {
		return nil, openErr
	}

	query := "CALL get_user_by_id(?);"

	ctx := context.Background()
	stmt, err := ur.sql.PrepareContext(ctx, query)
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
	var token *datamodel.UserToken

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

// GetByUsername attempts to get a user from the database with the given username.
// With the given username, an attempt to get the user's id is made, on success
// Get() is then called, which handles the getting and reading the user.
func (ur *userRepository) GetByUsername(username string) (*model.User, errors.Error) {
	if openErr := ur.openConnection(); openErr != nil {
		return nil, openErr
	}

	query := "CALL get_user_id_by_username(?);"

	ctx := context.Background()
	stmt, err := ur.sql.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.InternalError(fmt.Errorf("prepare: %v", err))
	}
	defer stmt.Close()

	var userID string

	err = stmt.QueryRowContext(ctx, username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFound("user not found")
		}

		return nil, errors.InternalError(fmt.Errorf("read: %v", err))
	}

	return ur.Get(userID)
}

func (ur *userRepository) GetByStateToken(stateToken string) (*model.User, errors.Error) {
	if openErr := ur.openConnection(); openErr != nil {
		return nil, openErr
	}

	query := "CALL get_user_id_by_state_token(?);"

	ctx := context.Background()
	stmt, err := ur.sql.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.InternalError(fmt.Errorf("prepare: %v", err))
	}
	defer stmt.Close()

	var userID string

	err = stmt.QueryRowContext(ctx, stateToken).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFound("user not found")
		}

		return nil, errors.InternalError(fmt.Errorf("read: %v", err))
	}

	return ur.Get(userID)
}

func (ur *userRepository) GetList(term string) ([]*model.User, errors.Error) {
	dms, err := ur.db.Read("CALL get_users(?);", func(s database.ScannerFunc) (interface{}, errors.Error) {
		var dm datamodel.User

		err := s(
			&dm.ID,
			&dm.Username,
			&dm.PasswordHash,
			&dm.StateToken,
			&dm.Enabled,
			&dm.AccountID,
		)
		if err != nil {
			return nil, errors.InternalError(err)
		}

		return &dm, nil
	}, term)
	if err != nil {
		return nil, err
	}

	u := make([]*model.User, len(dms))

	for i, dm := range dms {
		u[i] = model.UserFromDataModel(dm.(*datamodel.User), nil, nil)
	}

	return u, nil
}

func (ur *userRepository) GetPending(term string) ([]*model.User, errors.Error) {
	dms, err := ur.db.Read("CALL get_pending_users(?);", func(s database.ScannerFunc) (interface{}, errors.Error) {
		var dm datamodel.User

		err := s(
			&dm.ID,
			&dm.Username,
			&dm.PasswordHash,
			&dm.StateToken,
			&dm.Enabled,
			&dm.AccountID,
		)
		if err != nil {
			return nil, errors.InternalError(err)
		}

		return &dm, nil
	}, term)
	if err != nil {
		return nil, err
	}

	u := make([]*model.User, len(dms))

	for i, dm := range dms {
		u[i] = model.UserFromDataModel(dm.(*datamodel.User), nil, nil)
	}

	return u, nil
}

func readUser(s scannerFunc) (*datamodel.User, errors.Error) {
	var dm datamodel.User

	err := s(
		&dm.ID,
		&dm.Username,
		&dm.PasswordHash,
		&dm.StateToken,
		&dm.Enabled,
		&dm.AccountID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFound("user not found")
		}

		return nil, errors.InternalError(fmt.Errorf("read user: %v", err))
	}

	return &dm, nil
}

func readUserRole(s scannerFunc) (*datamodel.Role, errors.Error) {
	var dm datamodel.Role

	err := s(&dm.ID, &dm.Name)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.InternalError(fmt.Errorf("read role: %v", err))
	}

	return &dm, nil
}

func readToken(s scannerFunc) (*datamodel.UserToken, errors.Error) {
	var dm datamodel.UserToken

	err := s(
		&dm.AccessToken,
		&dm.RefreshToken,
		&dm.Expires,
		&dm.TokenType,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.InternalError(fmt.Errorf("read token: %v", err))
	}

	return &dm, nil
}

func (ur *userRepository) EnsureExists(id string) errors.Error {
	if openErr := ur.openConnection(); openErr != nil {
		return openErr
	}

	query := "SELECT id FROM users WHERE id = ?;"

	ctx := context.Background()
	stmt, err := ur.sql.PrepareContext(ctx, query)
	if err != nil {
		return errors.InternalError(err)
	}
	defer stmt.Close()

	var userID string
	err = stmt.QueryRowContext(ctx, id).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NotFound("user not found")
		}

		return errors.InternalError(err)
	}

	return nil
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
	query := "CALL update_user(?,?,?,?,?,?);"
	args := []interface{}{
		dm.ID,
		dm.Username,
		dm.AccountID,
		dm.StateToken,
		dm.PasswordHash,
		dm.Enabled,
	}

	openErr := ur.openConnection()
	if openErr != nil {
		return openErr
	}

	ctx := context.Background()
	tx, err := ur.sql.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
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

	err = updateToken(ctx, tx, u)
	if err != nil {
		return errors.InternalError(err)
	}

	return u.DispatchEvents(ctx, tx)
}

func updateToken(ctx context.Context, tx *sql.Tx, u *model.User) error {
	ut := u.GetToken()
	if ut == nil {
		return nil
	}

	dm := ut.DataModel()

	query := "CALL update_user_token(?,?,?,?,?);"
	args := []interface{}{
		u.GetID(),
		dm.AccessToken,
		dm.RefreshToken,
		dm.Expires,
		dm.TokenType,
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	return nil
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
	tx, err := ur.sql.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
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

		ur.sql = db
	}

	ctx := context.Background()
	err := ur.sql.PingContext(ctx)
	if err != nil {
		return errors.InternalError(err)
	}

	return nil
}
