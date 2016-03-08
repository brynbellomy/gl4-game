package entity

import (
	"errors"

	"github.com/brynbellomy/go-structomancer"
)

type (
	Factory struct {
		entityTemplates map[string][]IComponent
		componentTypes  map[string]cmptType
	}

	cmptType struct {
		z      *structomancer.Structomancer
		initer IComponentIniter
	}

	IComponentIniter interface {
		InitComponent(cmpt IComponent) error
	}
)

func NewFactory() *Factory {
	return &Factory{
		entityTemplates: map[string][]IComponent{},
		componentTypes:  map[string]cmptType{},
	}
}

func (f *Factory) RegisterEntityTemplate(name string, cmpts []IComponent) {
	f.entityTemplates[name] = cmpts
}

func (f *Factory) RegisterComponentType(typeName string, cmpt IComponent, cmptIniter IComponentIniter) {
	f.componentTypes[typeName] = cmptType{
		z:      structomancer.New(cmpt, "config"),
		initer: cmptIniter,
	}
}

// func (f *Factory) EntityFromTemplate(name string) ([]IComponent, error) {
// 	// f.entityTemplates[]
// }

type entityConfig struct {
	ID         ID                       `config:"id"`
	Components []map[string]interface{} `config:"components"`
}

var entityZ = structomancer.New(&entityConfig{}, "config")

func (f *Factory) EntityFromConfig(cfg map[string]interface{}) (ID, []IComponent, error) {
	c, err := entityZ.MapToStruct(cfg)
	if err != nil {
		return 0, nil, errors.New("error deserializing entity from config: " + err.Error())
	}

	config := c.(*entityConfig)

	cmpts := make([]IComponent, len(config.Components))
	for i, cmptcfg := range config.Components {
		cmpt, err := f.ComponentFromConfig(cmptcfg)
		if err != nil {
			return 0, nil, errors.New("error deserializing component from config: " + err.Error())
		}
		cmpts[i] = cmpt
	}

	return config.ID, cmpts, nil
}

type componentConfig struct {
	Type   string                 `config:"type"`
	Config map[string]interface{} `config:"config"`
}

var componentZ = structomancer.New(&componentConfig{}, "config")

func (f *Factory) ComponentFromConfig(cmptcfg map[string]interface{}) (IComponent, error) {
	c, err := componentZ.MapToStruct(cmptcfg)
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

	maybeCmpt, err := ctype.z.MapToStruct(cfg.Config)
	if err != nil {
		return nil, errors.New("error deserializing component (type = " + cfg.Type + ")" + err.Error())
	}

	cmpt := maybeCmpt.(IComponent)

	if ctype.initer != nil {
		err := ctype.initer.InitComponent(cmpt)
		if err != nil {
			return nil, err
		}
	}

	return cmpt, nil
}
