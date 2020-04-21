package datamodel

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

// User acts as a persistence model/data access object, used to map the User
// domain model to the data model.
type User struct {
	ID           string
	Username     string
	PasswordHash string
	StateToken   string
	Enabled      mysql.NullTime
	AccountID    sql.NullString
}
