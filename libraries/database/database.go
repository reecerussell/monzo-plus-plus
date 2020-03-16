package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
)

// ConnectionString is a connection read from the CONN_STRING
// environment variable, therefore requires it to be set.
var ConnectionString = os.Getenv("CONN_STRING")

// A set of friendlier, generic error messages.
var (
	ErrConnectionFailed    = fmt.Errorf("failed to open a connection to the database")
	ErrPingFailed          = fmt.Errorf("failed to communicate with the database")
	ErrTransactionFailed   = fmt.Errorf("failed to begin a transaction")
	ErrPrepareFailed       = fmt.Errorf("an error occured while preparing the database")
	ErrExecutionFailed     = fmt.Errorf("failed to execute a command on the database")
	ErrFailedToReadResults = fmt.Errorf("failed to read the results of the database")
	ErrScanFailed          = fmt.Errorf("failed to record data")
)

func init() {
	if ConnectionString == "" {
		panic("environment variable 'CONN_STRING' must be set")
	}
}

// ScannerFunc is used to provide a reader with a scan method.
type ScannerFunc func(dsts ...interface{}) error

// ReaderFunc is used as a generic read function to read
// query results for different records and types.
type ReaderFunc func(s ScannerFunc) (interface{}, errors.Error)

// DB is a wrapper type used to provide more abstract set
// of features and methods. These methods make interacting with
// sql.DB easier.
type DB struct {
	sql *sql.DB
}

// New returns a new instance of DB.
func New() *DB {
	return new(DB)
}

// EnsureConnected ensures there is an open connection to the database,
// ready for use. If no connection has been opened, a new one is made.
// Otherwise, the database is pinged to ensure a connection is alive.
func (db *DB) EnsureConnected() errors.Error {
	if db.sql == nil {
		sqlDB, err := sql.Open("mysql", ConnectionString)
		if err != nil {
			log.Printf("ERROR: %v", err)
			return errors.InternalError(ErrConnectionFailed)
		}

		db.sql = sqlDB
	}

	if err := db.sql.PingContext(context.Background()); err != nil {
		log.Printf("ERROR: %v", err)
		return errors.InternalError(ErrPingFailed)
	}

	return nil
}

// Execute takes a query and a set of arguments, then executes the query
// on the database. Returned is an integer value which shows the number of
// rows affected in the execution. Additionally, an error interface which will
// only have a non-nil value if an error occured during execution.
func (db *DB) Execute(query string, args ...interface{}) (int, errors.Error) {
	if err := db.EnsureConnected(); err != nil {
		return 0, err
	}

	ctx := context.Background()
	tx, err := db.sql.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	if err != nil {
		log.Printf("ERROR: %v", err)
		return 0, errors.InternalError(ErrTransactionFailed)
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
		log.Printf("ERROR: %v", err)
		return 0, errors.InternalError(ErrPrepareFailed)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return 0, errors.InternalError(ErrExecutionFailed)
	}

	rows, _ := res.RowsAffected()

	return int(rows), nil
}

// ReadOne is used to read a single record from a result set, given the query
// and arguments provided. The interface returned comes directly from the
// ReaderFunc. An error is only returned if the reader returns an error or
// if there was an issue with speaking to the database.
func (db *DB) ReadOne(query string, reader ReaderFunc, args ...interface{}) (interface{}, errors.Error) {
	if err := db.EnsureConnected(); err != nil {
		return nil, err
	}

	ctx := context.Background()
	stmt, err := db.sql.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return nil, errors.InternalError(ErrPrepareFailed)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, args...)
	item, readerErr := reader(row.Scan)
	if readerErr != nil {
		log.Printf("ERROR: %v", readerErr.Text())
		return nil, readerErr
	}

	return item, nil
}

// Read reads a set of records from a result set for the given query and arguments.
// Returned is an array of the records read, alongside an error. The error will
// only have a non-nil value if the ReaderFunc returned an error or if there was
// a problem communicating with the database.
func (db *DB) Read(query string, reader ReaderFunc, args ...interface{}) ([]interface{}, errors.Error) {
	if err := db.EnsureConnected(); err != nil {
		return nil, err
	}

	ctx := context.Background()
	stmt, err := db.sql.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return nil, errors.InternalError(ErrPrepareFailed)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return nil, errors.InternalError(ErrFailedToReadResults)
	}
	defer rows.Close()

	items := []interface{}{}

	for rows.Next() {
		item, readerErr := reader(rows.Scan)
		if readerErr != nil {
			log.Panicf("ERROR: %s", readerErr.Text())
			return nil, readerErr
		}

		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			return make([]interface{}, 0), nil
		}

		log.Printf("ERROR: %v", err)
		return nil, errors.InternalError(err)
	}

	return items, nil
}

// Count is used to read a record with a single integer value. An example
// of which could be a "SELECT COUNT(*)". The result is then returned.
func (db *DB) Count(query string, args ...interface{}) (int, errors.Error) {
	if err := db.EnsureConnected(); err != nil {
		return 0, err
	}

	ctx := context.Background()
	stmt, err := db.sql.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return 0, errors.InternalError(ErrPrepareFailed)
	}
	defer stmt.Close()

	var count int64

	err = stmt.QueryRowContext(ctx, args...).Scan(&count)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return 0, errors.InternalError(ErrFailedToReadResults)
	}

	return int(count), nil
}

// Exists uses the Count function to determine if a record exists. If the
// number result of Count is greater than 0, true is returned, otherwise false.
func (db *DB) Exists(query string, args ...interface{}) (bool, errors.Error) {
	c, err := db.Count(query, args...)
	if err != nil {
		return false, err
	}

	return c > 0, nil
}
