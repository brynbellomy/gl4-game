package entity

import (
	"errors"
	"reflect"

	"github.com/listenonrepeat/backend/common/structomancer"
)

type (
	Factory struct {
		componentTypes map[string]cmptType
	}

	cmptType struct {
		z      *structomancer.Structomancer
		initer IComponentIniter
	}

	IComponentIniter interface {
		InitComponent(cmpt IComponent)
	}
)

func NewFactory() *Factory {
	return &Factory{
		componentTypes: map[string]cmptType{},
	}
}

func (f *Factory) RegisterComponentType(typeName string, cmpt IComponent, cmptIniter IComponentIniter) {
	f.componentTypes[typeName] = cmptType{
		z:      structomancer.New(cmpt, "config"),
		initer: cmptIniter,
	}
}

var entityIDType = reflect.TypeOf(ID(0))

func (f *Factory) EntityFromConfig(cfg map[string]interface{}) (ID, []IComponent, error) {
	id, exists := cfg["id"]
	if !exists {
		return 0, nil, errors.New("key 'id' is missing")
	} else if !reflect.TypeOf(id).ConvertibleTo(entityIDType) {
		return 0, nil, errors.New("key 'id' is wrong type")
	}

	eid := reflect.ValueOf(id).Convert(entityIDType).Interface().(ID)

	cmptcfgs, exists := cfg["components"].([]interface{})
	if !exists {
		return 0, nil, errors.New("key 'components' is missing")
	}

	cmpts := make([]IComponent, len(cmptcfgs))

	for i, cmptcfg := range cmptcfgs {
		cmpt, err := f.ComponentFromConfig(cmptcfg.(map[string]interface{}))
		if err != nil {
			return 0, nil, err
		}
		cmpts[i] = cmpt
	}

	return eid, cmpts, nil
}

func (f *Factory) ComponentFromConfig(cmptcfg map[string]interface{}) (IComponent, error) {
	t, exists := cmptcfg["type"].(string)
	if !exists {
		return nil, errors.New("missing 'type' key")
	}

	cfg, exists := cmptcfg["config"].(map[string]interface{})
	if !exists {
		return nil, errors.New("missing 'config' key")
	}

	ctype, exists := f.componentTypes[t]
	if !exists {
		return nil, errors.New("component type '" + t + "' is not registered")
	}

	maybeCmpt, err := ctype.z.MapToStruct(cfg)
	if err != nil {
		return nil, err
	}

	cmpt := maybeCmpt.(IComponent)

	if ctype.initer != nil {
		ctype.initer.InitComponent(cmpt)
	}

	return cmpt, nil
}
