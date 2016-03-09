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

func NewEntityFactory() *EntityFactory {
	return &EntityFactory{
		cmptFactory: NewComponentFactory(),
	}
}

func (f *EntityFactory) RegisterComponentType(typeName string, cmpt IComponent) {
	f.cmptFactory.RegisterComponentType(typeName, cmpt)
}

func (f *EntityFactory) EntityFromConfig(cfg map[string]interface{}) (ID, []IComponent, error) {
	c, err := entityConfigType.MapToStruct(cfg)
	if err != nil {
		return 0, nil, errors.New("error deserializing entity from config: " + err.Error())
	}

	config := c.(*entityConfig)

	cmpts := make([]IComponent, len(config.Components))
	for i, cmptcfg := range config.Components {
		cmpt, err := f.cmptFactory.ComponentFromConfig(cmptcfg)
		if err != nil {
			return 0, nil, errors.New("error deserializing component from config: " + err.Error())
		}
		cmpts[i] = cmpt
	}

	return config.ID, cmpts, nil
}
