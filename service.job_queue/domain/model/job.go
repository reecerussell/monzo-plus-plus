package model

import (
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.job_queue/domain/datamodel"
)

// ProcessFunc is a standard function used to handle the execution of a job.
type ProcessFunc func(userID, accountID, pluginName, data string) errors.Error

// Job is a domain object for job records which is used to create,
// read and execute jobs in the queue.
type Job struct {
	id           int
	userID       string
	accountID    string
	pluginID     string
	pluginName   string
	retryCount   int
	received     time.Time
	completed    *time.Time
	lastExecuted *time.Time
	data         string
}

// NewJob is used to create a new Job domain object.
func NewJob(userID, pluginID, data string) *Job {
	return &Job{
		userID:       userID,
		pluginID:     pluginID,
		retryCount:   0,
		received:     time.Now().UTC(),
		completed:    nil,
		lastExecuted: nil,
		data:         data,
	}
}

// Add for debugging purposes.
// TODO: remove method.
func (j *Job) GetID() int {
	return j.id
}

// MarkAsCompleted is used to set the completed field to the current time
// in UTC. Once a job has been marked as completed, it will not be added
// to the queue again.
//
// A non-nil error is returned if the job has already been completed.
func (j *Job) MarkAsCompleted() errors.Error {
	if j.completed != nil {
		return errors.BadRequest("this job has already been completed")
	}

	t := time.Now().UTC()
	j.completed = &t

	return nil
}

// Execute processes the job using the given processor and then updates
// the required fields accordingly. The job must not already be completed
// and not have been retried more than 3 times.
//
// A non-nil error is returned if either the job has already been
// completed, the job has been retried three times or if it failed
// to be processed.
func (j *Job) Execute(p ProcessFunc) errors.Error {
	if j.completed != nil {
		return errors.BadRequest("this job has already been completed")
	}

	if j.retryCount >= 3 {
		return errors.BadRequest("this job has already been execute the maximum number of times")
	}

	err := p(j.userID, j.accountID, j.pluginName, j.data)
	if err != nil {
		j.retryCount++
		j.setLastExecutedDate()

		return err
	}

	j.setLastExecutedDate()
	j.MarkAsCompleted()

	return nil
}

func (j *Job) setLastExecutedDate() {
	t := time.Now().UTC()
	j.lastExecuted = &t
}

// DataModel returns an instance of *datamodel.Job containing the job's data.
func (j *Job) DataModel() *datamodel.Job {
	dm := &datamodel.Job{
		ID:         j.id,
		UserID:     j.userID,
		PluginID:   j.pluginID,
		RetryCount: j.retryCount,
		Received:   j.received,
		Data:       j.data,
	}

	if j.completed == nil {
		dm.Completed = mysql.NullTime{
			Valid: false,
		}
	} else {
		dm.Completed = mysql.NullTime{
			Valid: true,
			Time:  *j.completed,
		}
	}

	if j.lastExecuted == nil {
		dm.LastExecuted = mysql.NullTime{
			Valid: false,
		}
	} else {
		dm.LastExecuted = mysql.NullTime{
			Valid: true,
			Time:  *j.lastExecuted,
		}
	}

	return dm
}

// JobFromDataModel returns a new instance of Job from the given data model.
func JobFromDataModel(dm *datamodel.Job) *Job {
	j := &Job{
		id:         dm.ID,
		userID:     dm.UserID,
		accountID:  dm.AccountID,
		pluginID:   dm.PluginID,
		pluginName: dm.PluginName,
		retryCount: dm.RetryCount,
		received:   dm.Received,
		data:       dm.Data,
	}

	if dm.Completed.Valid {
		j.completed = &dm.Completed.Time
	}

	if dm.LastExecuted.Valid {
		j.lastExecuted = &dm.LastExecuted.Time
	}

	return j
}
