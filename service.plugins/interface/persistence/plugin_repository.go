package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/reecerussell/monzo-plus-plus/libraries/util"

	"github.com/reecerussell/monzo-plus-plus/libraries/database"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/domain/datamodel"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/domain/repository"
)

type pluginRepository struct {
	db *database.DB
}

// NewPluginRepository returns a new instance of repository.PluginRepository.
func NewPluginRepository() repository.PluginRepository {
	return &pluginRepository{
		db: database.New(),
	}
}

func (pr *pluginRepository) GetList(ctx context.Context, term string) ([]*model.Plugin, errors.Error) {
	query := `CALL get_plugins(?,?)`
	term = fmt.Sprintf("%%%s%%", term)
	args := []interface{}{term, ctx.Value(util.ContextKey("user_id"))}

	items, err := pr.db.Read(query, func(s database.ScannerFunc) (interface{}, errors.Error) {
		var dm datamodel.Plugin

		if err := s(
			&dm.ID,
			&dm.Name,
			&dm.DisplayName,
			&dm.Description,
			&dm.ConsumedBy,
			&dm.ConsumedByUser); err != nil {
			log.Printf("ERROR: %v", err)
			return nil, errors.InternalError(database.ErrScanFailed)
		}

		return &dm, nil
	}, args...)
	if err != nil {
		return nil, err
	}

	plugins := make([]*model.Plugin, len(items))

	for i, dm := range items {
		plugins[i] = model.PluginFromDataModel(dm.(*datamodel.Plugin))
	}

	return plugins, nil
}

func (pr *pluginRepository) Get(ctx context.Context, id string) (*model.Plugin, errors.Error) {
	query := "CALL get_plugin(?,?);"
	args := []interface{}{id, ctx.Value(util.ContextKey("user_id"))}

	item, err := pr.db.ReadOne(query, func(s database.ScannerFunc) (interface{}, errors.Error) {
		var dm datamodel.Plugin

		if err := s(&dm.ID, &dm.Name, &dm.DisplayName, &dm.Description, &dm.ConsumedBy, &dm.ConsumedByUser); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.NotFound("plugin not found")
			}

			log.Printf("ERROR: %v", err)
			return nil, errors.InternalError(database.ErrScanFailed)
		}

		return &dm, nil
	}, args...)
	if err != nil {
		return nil, err
	}

	return model.PluginFromDataModel(item.(*datamodel.Plugin)), nil
}

func (pr *pluginRepository) Create(p *model.Plugin) errors.Error {
	query := "INSERT INTO `plugins` (`id`,`name`,`display_name`,`description`) VALUES (?,?,?,?);"
	dm := p.Datamodel()
	args := []interface{}{
		dm.ID, dm.Name, dm.DisplayName, dm.Description,
	}

	_, err := pr.db.Execute(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (pr *pluginRepository) Update(p *model.Plugin) errors.Error {
	query := "UPDATE `plugins` SET `name` = ?, display_name = ?, `description` = ? WHERE id = ?;"
	dm := p.Datamodel()
	args := []interface{}{
		dm.Name, dm.DisplayName, dm.Description, dm.ID,
	}

	ra, err := pr.db.Execute(query, args...)
	if err != nil {
		return err
	}

	if ra < 1 {
		return errors.NotFound("plugin not found")
	}

	return nil
}

func (pr *pluginRepository) Delete(id string) errors.Error {
	query := "DELETE FROM `plugins` WHERE id = ?;"

	ra, err := pr.db.Execute(query, id)
	if err != nil {
		return err
	}

	if ra < 1 {
		return errors.NotFound("plugin not found")
	}

	return nil
}
