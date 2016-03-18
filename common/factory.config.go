package common

import (
	"errors"

	"github.com/brynbellomy/go-structomancer"
)

type (
	ConfigFactory struct {
		configTypes map[string]*Coder
	}

	configWrapper struct {
		Type   string                 `config:"type"`
		Config map[string]interface{} `config:"config"`
	}
)

var configWrapperType = structomancer.New(&configWrapper{}, "config")

func NewConfigFactory() *ConfigFactory {
	return &ConfigFactory{
		configTypes: map[string]*Coder{},
	}
}

func (f *ConfigFactory) Register(name string, coder *Coder) {
	f.configTypes[name] = coder
}

func (f *ConfigFactory) Decode(data map[string]interface{}) (interface{}, error) {
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

	coder, exists := f.configTypes[cfg.Type]
	if !exists {
		return nil, errors.New("config type '" + cfg.Type + "' is not registered")
	}

	thing, err := coder.Decode(cfg.Config)
	if err != nil {
		return nil, errors.New("error deserializing config (type = " + cfg.Type + "): " + err.Error())
	}

	return thing, nil
}

func (f *ConfigFactory) Encode(data map[string]interface{}) (interface{}, error) {
	panic("not implemented")
}
