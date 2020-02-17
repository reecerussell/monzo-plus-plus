package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/reecerussell/monzo-plus-plus/service.mpp/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/domain/repository"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

// Environment variables.
var (
	ConnectionString = os.Getenv("CONN_STRING")
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository() repository.UserRepository {
	return &userRepository{}
}

func (ur *userRepository) FindByStateToken(token string) (*model.User, error) {
	err := ur.openConnection()
	if err != nil {
		return nil, err
	}

	query := "CALL get_user_by_state_token(?);"

	ctx := context.Background()
	stmt, err := ur.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare: %v", err)
	}
	defer stmt.Close()

	var (
		id, stateToken string
	)

	err = stmt.QueryRowContext(ctx, token).Scan(&id, &stateToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrUserNotFoundWithStateToken
		}

		return nil, fmt.Errorf("scan: %v", err)
	}

	return model.NewUserFromStateToken(id, stateToken), nil
}

func (ur *userRepository) Get(userID string) (*model.User, error) {
	err := ur.openConnection()
	if err != nil {
		return nil, err
	}

	query := "CALL get_user_by_id(?);"

	ctx := context.Background()
	stmt, err := ur.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare: %v", err)
	}
	defer stmt.Close()

	u, dsts := model.NewUserFromScan()
	err = stmt.QueryRowContext(ctx, userID).Scan(dsts...)
	if err != nil {
		return nil, fmt.Errorf("scan: %v", err)
	}

	return u, nil
}

func (ur *userRepository) Create(u *model.User) error {
	err := ur.openConnection()
	if err != nil {
		return err
	}

	query := "CALL create_new_user(?,?);"

	ctx := context.Background()
	tx, err := ur.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	if err != nil {
		return fmt.Errorf("tx: %v", err)
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
		return fmt.Errorf("prepare: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.GetID(), u.GetStateToken())
	if err != nil {
		return fmt.Errorf("exec: %v", err)
	}

	return nil
}

func (ur *userRepository) Update(u *model.User) error {
	err := ur.openConnection()
	if err != nil {
		return err
	}

	query := "CALL update_user(?,?,?,?,?,?);"
	args := []interface{}{
		u.GetID(),
		u.GetMonzoID(),
		u.GetToken().GetAccessToken(),
		u.GetToken().GetRefreshToken(),
		u.GetToken().GetExpiryDate(),
		u.GetToken().GetTokenType(),
	}

	ctx := context.Background()
	tx, err := ur.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	if err != nil {
		return fmt.Errorf("tx: %v", err)
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
		return fmt.Errorf("prepare: %v", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return fmt.Errorf("exec: %v", err)
	}

	if c, _ := res.RowsAffected(); c < 1 {
		return repository.ErrNoUserWasUpdated
	}

	return nil
}

// openConnection opens a connection to the database if one doesn't exist.
func (ur *userRepository) openConnection() error {
	if ur.db == nil {
		db, err := sql.Open("mysql", ConnectionString)
		if err != nil {
			return fmt.Errorf("open connection: %v", err)
		}

		ur.db = db
	}

	ctx := context.Background()
	err := ur.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("open connection: ping: %v", err)
	}

	return nil
}
