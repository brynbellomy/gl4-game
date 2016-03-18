package entity

import (
	"errors"

	"github.com/brynbellomy/go-structomancer"
)

type (
	EntityFactory struct {
		cmptFactory *ComponentFactory
	}

	entityConfig struct {
		ID         ID                       `config:"id"`
		Components []map[string]interface{} `config:"components"`
	}
)

var entityConfigType = structomancer.New(&entityConfig{}, "config")

func NewEntityFactory(cmptRegistry *ComponentRegistry) *EntityFactory {
	return &EntityFactory{
		cmptFactory: NewComponentFactory(cmptRegistry),
	}
}

func (f *EntityFactory) EntityFromConfig(cfg map[string]interface{}) (Entity, error) {
	c, err := entityConfigType.MapToStruct(cfg)
	if err != nil {
		return Entity{}, errors.New("error deserializing entity from config: " + err.Error())
	}

	config := c.(*entityConfig)

	cmpts := make([]IComponent, len(config.Components))
	kinds := make([]ComponentKind, len(config.Components))
	mask := ComponentMask(0)
	for i, cmptcfg := range config.Components {
		cmpt, kind, err := f.cmptFactory.ComponentFromConfig(cmptcfg)
		if err != nil {
			return Entity{}, errors.New("error deserializing component from config: " + err.Error())
		}
		cmpts[i] = cmpt
		kinds[i] = kind
		mask = mask.Add(kind)
	}

	return Entity{ID: config.ID, ComponentMask: mask, ComponentKinds: kinds, Components: cmpts}, nil
}
