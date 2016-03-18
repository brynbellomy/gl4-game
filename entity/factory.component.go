package entity

import (
	"errors"

	"github.com/brynbellomy/go-structomancer"
)

type (
	ComponentFactory struct {
		componentRegistry *ComponentRegistry
	}

	componentConfig struct {
		Type   string                 `config:"type"`
		Config map[string]interface{} `config:"config"`
	}
)

var componentConfigType = structomancer.New(&componentConfig{}, "config")

func NewComponentFactory(componentRegistry *ComponentRegistry) *ComponentFactory {
	return &ComponentFactory{
		componentRegistry: componentRegistry,
	}
}

func (f *ComponentFactory) ComponentFromConfig(cmptcfg map[string]interface{}) (IComponent, ComponentKind, error) {
	c, err := componentConfigType.MapToStruct(cmptcfg)
	if err != nil {
		return nil, 0, err
	}

	cfg := c.(*componentConfig)
	if cfg.Type == "" {
		return nil, 0, errors.New("missing 'type' key")
	} else if cfg.Config == nil {
		return nil, 0, errors.New("missing 'config' key")
	}

	ctype, exists := f.componentRegistry.GetComponentType(cfg.Type)
	if !exists {
		return nil, 0, errors.New("component type '" + cfg.Type + "' is not registered")
	}

	if cfg.Config == nil {
		return nil, 0, errors.New("error deserializing component (type = " + cfg.Type + "): config is nil")
	}

	cmpt, err := ctype.Coder.Decode(cfg.Config)
	if err != nil {
		return nil, 0, errors.New("error deserializing component (type = " + cfg.Type + ")" + err.Error())
	}

	return cmpt.(IComponent), ctype.Kind, nil
}
