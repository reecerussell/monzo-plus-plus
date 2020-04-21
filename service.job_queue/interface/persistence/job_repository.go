package persistence

import (
	"github.com/reecerussell/monzo-plus-plus/libraries/database"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.job_queue/domain/datamodel"
	"github.com/reecerussell/monzo-plus-plus/service.job_queue/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.job_queue/domain/repository"
)

// jobRepository is a MySQL instance of JobRepository.
type jobRepository struct {
	db *database.DB
}

// NewJobRepository returns a new instance of jobRepository.
func NewJobRepository() repository.JobRepository {
	return &jobRepository{
		db: database.New(),
	}
}

// Add inserts a new Job record into the database.
func (jr *jobRepository) Add(j *model.Job) errors.Error {
	dm := j.DataModel()
	query := "CALL create_job(?,?,?,?);"
	args := []interface{}{dm.UserID, dm.PluginID, dm.Received, dm.Data}

	_, err := jr.db.Execute(query, args...)
	if err != nil {
		return err
	}

	return nil
}

// Update modifies an existing Job recor din the database.
func (jr *jobRepository) Update(j *model.Job) errors.Error {
	dm := j.DataModel()
	query := "CALL update_job(?,?,?,?);"
	args := []interface{}{dm.ID, dm.RetryCount, dm.Completed, dm.LastExecuted}

	_, err := jr.db.Execute(query, args...)
	if err != nil {
		return err
	}

	return nil
}

// GetN returns N number of Jobs which are the next jobs to be processed.
func (jr *jobRepository) GetN(n int) ([]*model.Job, errors.Error) {
	dms, err := jr.db.Read("CALL get_latest_jobs(?);", func(s database.ScannerFunc) (interface{}, errors.Error) {
		var dm datamodel.Job

		err := s(
			&dm.ID,
			&dm.UserID,
			&dm.AccountID,
			&dm.PluginID,
			&dm.PluginName,
			&dm.RetryCount,
			&dm.Received,
			&dm.Completed,
			&dm.LastExecuted,
			&dm.Data,
		)
		if err != nil {
			return nil, errors.InternalError(err)
		}

		return &dm, nil
	}, n)
	if err != nil {
		return nil, err
	}

	models := make([]*model.Job, len(dms))

	for i, dm := range dms {
		models[i] = model.JobFromDataModel(dm.(*datamodel.Job))
	}

	return models, nil
}
