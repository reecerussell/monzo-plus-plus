package repository

import (
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.job_queue/domain/model"
)

// JobRepository is a high-level interface used to read jobs from a data source.
type JobRepository interface {
	Add(j *model.Job) errors.Error
	Update(j *model.Job) errors.Error
	GetN(n int) ([]*model.Job, errors.Error)
}
