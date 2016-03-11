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

func (f *ComponentFactory) ComponentFromConfig(cmptcfg map[string]interface{}) (IComponent, error) {
	c, err := componentConfigType.MapToStruct(cmptcfg)
	if err != nil {
		return nil, err
	}

	cfg := c.(*componentConfig)
	if cfg.Type == "" {
		return nil, errors.New("missing 'type' key")
	} else if cfg.Config == nil {
		return nil, errors.New("missing 'config' key")
	}

	ctype, exists := f.componentRegistry.GetComponentType(cfg.Type)
	if !exists {
		return nil, errors.New("component type '" + cfg.Type + "' is not registered")
	}

	cmpt, err := ctype.DeserializeConfig(cfg.Config)
	if err != nil {
		return nil, errors.New("error deserializing component (type = " + cfg.Type + ")" + err.Error())
	}

	cmpt.(IComponent).SetKind(ctype.kind)

	return cmpt.(IComponent), nil
}
