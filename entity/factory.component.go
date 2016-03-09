package entity

import (
	"errors"

	"github.com/brynbellomy/go-structomancer"
)

type (
	ComponentFactory struct {
		componentTypes map[string]cmptType
	}

	componentConfig struct {
		Type   string                 `config:"type"`
		Config map[string]interface{} `config:"config"`
	}

	cmptType struct {
		z *structomancer.Structomancer
	}
)

var componentConfigType = structomancer.New(&componentConfig{}, "config")

func NewComponentFactory() *ComponentFactory {
	return &ComponentFactory{
		componentTypes: map[string]cmptType{},
	}
}

func (f *ComponentFactory) RegisterComponentType(typeName string, cmpt IComponent) {
	f.componentTypes[typeName] = cmptType{
		z: structomancer.New(cmpt, "config"),
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

	ctype, exists := f.componentTypes[cfg.Type]
	if !exists {
		return nil, errors.New("component type '" + cfg.Type + "' is not registered")
	}

	cmpt, err := ctype.z.MapToStruct(cfg.Config)
	if err != nil {
		return nil, errors.New("error deserializing component (type = " + cfg.Type + ")" + err.Error())
	}

	return cmpt.(IComponent), nil
}
