package common

import (
	"github.com/brynbellomy/go-structomancer"
)

type (
	Coder struct {
		configStruct *structomancer.Structomancer
		encode       func(interface{}) (interface{}, error)
		decode       func(interface{}) (interface{}, error)
	}

	CoderConfig struct {
		ConfigType interface{}
		Tag        string
		Encode     func(interface{}) (interface{}, error)
		Decode     func(interface{}) (interface{}, error)
	}
)

func NewCoder(cfg CoderConfig) *Coder {
	return &Coder{
		configStruct: structomancer.New(cfg.ConfigType, cfg.Tag),
		encode:       cfg.Encode,
		decode:       cfg.Decode,
	}
}

func (c *Coder) Decode(cfgMap map[string]interface{}) (interface{}, error) {
	cfg, err := c.configStruct.MapToStruct(cfgMap)
	if err != nil {
		return nil, err
	}

	return c.decode(cfg)
}

func (c *Coder) Encode(x interface{}) (map[string]interface{}, error) {
	cfg, err := c.encode(x)
	if err != nil {
		return nil, err
	}

	return c.configStruct.StructToMap(cfg)
}
