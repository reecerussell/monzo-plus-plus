package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/domain/repository"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

// Environment variables.
var (
	ConnectionString = os.Getenv("CONN_STRING")
)

// preferencesRepository is used to manage persistent data to a MySQL source.
type preferencesRepository struct {
	db *sql.DB
}

// NewPreferencesRepository returns a new instance of PreferencesRepository.
func NewPreferencesRepository() repository.PreferencesRepository {
	return &preferencesRepository{}
}

func (pr *preferencesRepository) Get(userID string) (*model.Preferences, error) {
	err := pr.openConnection()
	if err != nil {
		return nil, fmt.Errorf("open: %v", err)
	}

	query := `CALL get_user_budget_preferences(?);`

	ctx := context.Background()
	tx, err := pr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	if err != nil {
		return nil, fmt.Errorf("tx: %v", err)
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
		return nil, fmt.Errorf("prepare: %v", err)
	}
	defer stmt.Close()

	p, dsts := model.PreferencesFromScan()
	err = stmt.QueryRowContext(ctx, userID).Scan(dsts...)
	if err != nil {
		return nil, fmt.Errorf("scan: %v", err)
	}

	return p, nil
}

func (pr *preferencesRepository) Update(p *model.Preferences) error {
	err := pr.openConnection()
	if err != nil {
		return fmt.Errorf("open: %v", err)
	}

	query := "CALL update_user_budget_preferences(?,?);"

	ctx := context.Background()
	tx, err := pr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
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

	_, err = stmt.ExecContext(ctx, p.GetUserID(), p.GetMonthlyBudget())
	if err != nil {
		return fmt.Errorf("exec: %v", err)
	}

	return nil
}

// openConnection opens a connection to the database if one doesn't exist.
func (pr *preferencesRepository) openConnection() error {
	if pr.db == nil {
		db, err := sql.Open("mysql", ConnectionString)
		if err != nil {
			return fmt.Errorf("open connection: %v", err)
		}

		pr.db = db
	}

	ctx := context.Background()
	err := pr.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("open connection: ping: %v", err)
	}

	return nil
}
