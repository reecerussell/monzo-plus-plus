package usecase

import (
	"context"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/domain/dto"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/domain/repository"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/domain/service"
)

// PluginUsecase is a high-level usecase used to manage the plugin domain.
type PluginUsecase interface {
	All(ctx context.Context, term string) ([]*dto.Plugin, errors.Error)
	Get(ctx context.Context, id string) (*dto.Plugin, errors.Error)
	Create(d *dto.CreatePlugin) (*dto.Plugin, errors.Error)
	Update(ctx context.Context, d *dto.UpdatePlugin) errors.Error
	Delete(id string) errors.Error
}

type pluginUsecase struct {
	serv *service.PluginService
	repo repository.PluginRepository
}

// NewPluginUsecase returns a new instance of PluginUsecase, populated
// with the given dependencies.
func NewPluginUsecase(serv *service.PluginService, repo repository.PluginRepository) PluginUsecase {
	return &pluginUsecase{
		serv: serv,
		repo: repo,
	}
}

// All returns an array of all plugins matching the term below. The term
// can be empty, which will result in all plugins.
func (pu *pluginUsecase) All(ctx context.Context, term string) ([]*dto.Plugin, errors.Error) {
	plugins, err := pu.repo.GetList(ctx, term)
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.Plugin, len(plugins))

	for i, p := range plugins {
		dtos[i] = p.DTO()
	}

	return dtos, nil
}

// Get returns a single plugin record.
func (pu *pluginUsecase) Get(ctx context.Context, id string) (*dto.Plugin, errors.Error) {
	plugin, err := pu.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return plugin.DTO(), nil
}

// Create instansiates a new plugin record then inserts it into the database.
func (pu *pluginUsecase) Create(d *dto.CreatePlugin) (*dto.Plugin, errors.Error) {
	p, err := model.NewPlugin(d)
	if err != nil {
		return nil, err
	}

	err = pu.serv.EnsureUniqueName(p)
	if err != nil {
		return nil, err
	}

	err = pu.repo.Create(p)
	if err != nil {
		return nil, err
	}

	return p.DTO(), nil
}

// Update updates a plugin record using the data given by the dto.
func (pu *pluginUsecase) Update(ctx context.Context, d *dto.UpdatePlugin) errors.Error {
	p, err := pu.repo.Get(ctx, d.ID)
	if err != nil {
		return err
	}

	err = p.Update(d)
	if err != nil {
		return err
	}

	if !p.HasBeenUpdated() {
		return nil
	}

	err = pu.serv.EnsureUniqueName(p)
	if err != nil {
		return err
	}

	err = pu.repo.Update(p)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a plugin record from the database.
func (pu *pluginUsecase) Delete(id string) errors.Error {
	return pu.repo.Delete(id)
}
