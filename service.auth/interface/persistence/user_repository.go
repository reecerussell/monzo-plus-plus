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
	query := "CALL get_user_by_id(?);"
	readers := []database.ReaderFunc{readUser, readUserRole, readToken}
	results, err := ur.db.ReadMultiple(query, readers, id)
	if err != nil {
		return nil, err
	}

	user := results[0][0].(*datamodel.User)
	roles := make([]*datamodel.Role, len(results[1]))

	for i, dm := range results[1] {
		roles[i] = dm.(*datamodel.Role)
	}

	token := results[2][0].(*datamodel.UserToken)

	return model.UserFromDataModel(user, roles, token), nil
}

func readUser(s database.ScannerFunc) (interface{}, errors.Error) {
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

func readUserRole(s database.ScannerFunc) (interface{}, errors.Error) {
	var dm datamodel.Role

	err := s(&dm.ID, &dm.Name)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.InternalError(fmt.Errorf("read role: %v", err))
	}

	return &dm, nil
}

func readToken(s database.ScannerFunc) (interface{}, errors.Error) {
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

// GetByUsername attempts to get a user from the database with the given username.
// With the given username, an attempt to get the user's id is made, on success
// Get() is then called, which handles the getting and reading the user.
func (ur *userRepository) GetByUsername(username string) (*model.User, errors.Error) {
	query := "CALL get_user_id_by_username(?);"
	id, err := ur.db.ReadOne(query, func(s database.ScannerFunc) (interface{}, errors.Error) {
		var userID string

		err := s(&userID)
		if err != nil {
			return nil, errors.InternalError(err)
		}

		return userID, nil
	}, username)
	if err != nil {
		return nil, err
	}

	return ur.Get(id.(string))
}

func (ur *userRepository) GetByStateToken(stateToken string) (*model.User, errors.Error) {
	query := "CALL get_user_id_by_state_token(?);"
	id, err := ur.db.ReadOne(query, func(s database.ScannerFunc) (interface{}, errors.Error) {
		var userID string

		err := s(&userID)
		if err != nil {
			return nil, errors.InternalError(err)
		}

		return userID, nil
	}, stateToken)
	if err != nil {
		return nil, err
	}

	return ur.Get(id.(string))
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

func (ur *userRepository) EnsureExists(id string) errors.Error {
	query := "SELECT COUNT(*) FROM users WHERE id = ?;"

	c, err := ur.db.Count(query, id)
	if err != nil {
		return err
	}

	if c < 1 {
		return errors.BadRequest("user not found")
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

	_, err := ur.db.Execute(query, args...)
	if err != nil {
		return err
	}

	return nil
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
	tx, err := ur.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Close()

	_, err = tx.Execute(query, args...)
	if err != nil {
		return err
	}

	err = updateToken(ctx, tx, u)
	if err != nil {
		tx.InternalTx.Rollback()
		return err
	}

	err = u.DispatchEvents(ctx, tx.InternalTx)
	if err != nil {
		tx.InternalTx.Rollback()
		return err
	}

	return nil
}

func updateToken(ctx context.Context, tx *database.Tx, u *model.User) errors.Error {
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

	_, err := tx.Execute(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) Delete(id string) errors.Error {
	query := "DELETE FROM users WHERE id = ?;"

	c, err := ur.db.Execute(query, id)
	if err != nil {
		return err
	}

	if c < 1 {
		return errors.NotFound("user not found")
	}

	return nil
}

func (ur *userRepository) openConnection() errors.Error {
	if ur.sql == nil {
		db, err := sql.Open("mysql", os.Getenv("CONN_STRING"))
		if err != nil {
			return errors.InternalError(err)
		}

		ur.sql = db
	}

	err := ur.sql.Ping()
	if err != nil {
		return errors.InternalError(err)
	}

	return nil
}
