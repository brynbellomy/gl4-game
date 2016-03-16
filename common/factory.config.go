package common

import (
	"errors"

	"github.com/brynbellomy/go-structomancer"
)

type (
	ConfigFactory struct {
		configTypes map[string]configType
	}

	configType struct {
		z *structomancer.Structomancer
	}

	configWrapper struct {
		Type   string                 `config:"type"`
		Config map[string]interface{} `config:"config"`
	}
)

var configWrapperType = structomancer.New(&configWrapper{}, "config")

func NewConfigFactory() *ConfigFactory {
	return &ConfigFactory{
		configTypes: map[string]configType{},
	}
}

func (f *ConfigFactory) Register(name string, specimen interface{}) {
	f.configTypes[name] = configType{z: structomancer.New(specimen, "config")}
}

func (f *ConfigFactory) Build(data map[string]interface{}) (interface{}, error) {
	c, err := configWrapperType.MapToStruct(data)
	if err != nil {
		return nil, err
	}

	cfg := c.(*configWrapper)
	if cfg.Type == "" {
		return nil, errors.New("missing 'type' key")
	} else if cfg.Config == nil {
		return nil, errors.New("missing 'config' key")
	}

	ctype, exists := f.configTypes[cfg.Type]
	if !exists {
		return nil, errors.New("config type '" + cfg.Type + "' is not registered")
	}

	config, err := ctype.z.MapToStruct(cfg.Config)
	if err != nil {
		return nil, errors.New("error deserializing config (type = " + cfg.Type + ")" + err.Error())
	}

	return config, nil
}
