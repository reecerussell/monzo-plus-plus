package model

import (
	"github.com/google/uuid"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/domain/datamodel"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/domain/dto"
)

type Plugin struct {
	id          string
	name        string
	displayName string
	description string
}

func NewPlugin(d *dto.CreatePlugin) (*Plugin, errors.Error) {
	p := new(Plugin)
	p.id = uuid.New().String()

	err := p.UpdateName(d.Name)
	if err != nil {
		return nil, err
	}

	err = p.UpdateDisplayName(d.DisplayName)
	if err != nil {
		return nil, err
	}

	err = p.UpdateDescription(d.Description)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// GetID returns the plugin's id.
func (p *Plugin) GetID() string {
	return p.id
}

// GetName returns the plugin's name.
func (p *Plugin) GetName() string {
	return p.name
}

func (p *Plugin) Update(d *dto.UpdatePlugin) errors.Error {
	err := p.UpdateName(d.Name)
	if err != nil {
		return err
	}

	err = p.UpdateDisplayName(d.DisplayName)
	if err != nil {
		return err
	}

	err = p.UpdateDescription(d.Description)
	if err != nil {
		return err
	}

	return nil
}

func (p *Plugin) UpdateName(name string) errors.Error {
	if len(name) < 1 {
		return errors.BadRequest("name cannot be empty")
	}

	if len(name) > 45 {
		return errors.BadRequest("name cannot be greater than 45 characters long")
	}

	p.name = name

	return nil
}

func (p *Plugin) UpdateDisplayName(name string) errors.Error {
	if len(name) < 1 {
		return errors.BadRequest("display name cannot be empty")
	}

	if len(name) > 45 {
		return errors.BadRequest("display name cannot be greater than 45 characters long")
	}

	p.displayName = name

	return nil
}

func (p *Plugin) UpdateDescription(text string) errors.Error {
	if len(text) < 1 {
		return errors.BadRequest("description cannot be empty")
	}

	p.description = text

	return nil
}

func (p *Plugin) Datamodel() *datamodel.Plugin {
	return &datamodel.Plugin{
		ID:          p.id,
		Name:        p.name,
		DisplayName: p.displayName,
		Description: p.description,
	}
}

func (p *Plugin) DTO() *dto.Plugin {
	return &dto.Plugin{
		ID:          p.id,
		Name:        p.name,
		DisplayName: p.displayName,
		Description: p.description,
	}
}

func PluginFromDataModel(d *datamodel.Plugin) *Plugin {
	return &Plugin{
		id:          d.ID,
		name:        d.Name,
		displayName: d.DisplayName,
		description: d.Description,
	}
}
