package datamodel

import (
	"time"

	"github.com/go-sql-driver/mysql"
)

// Job is a data model object which is used to read data
// from the database, then used to transfer to a domain object.
type Job struct {
	ID           int
	UserID       string
	AccountID    string
	PluginID     string
	PluginName   string
	RetryCount   int
	Received     time.Time
	Completed    mysql.NullTime
	LastExecuted mysql.NullTime
	Data         string
}
