package entity

import (
	"github.com/brynbellomy/go-structomancer"
)

type (
	ComponentRegistry struct {
		componentTypes map[string]CmptType
		kindCounter    uint64
	}

	CmptType struct {
		name string
		kind ComponentKind
		z    *structomancer.Structomancer
	}
)

func NewComponentRegistry() *ComponentRegistry {
	return &ComponentRegistry{
		componentTypes: map[string]CmptType{},
		kindCounter:    0,
	}
}

func (r *ComponentRegistry) nextComponentKind() ComponentKind {
	next := r.kindCounter //1 << r.kindCounter
	r.kindCounter++
	if r.kindCounter == 64 {
		panic("entity.ComponentFactory: maximum number of component kinds exceeded (64)")
	}
	return ComponentKind(next)
}

func (r *ComponentRegistry) RegisterComponentType(typeName string, cmpt IComponent) {
	r.componentTypes[typeName] = CmptType{
		name: typeName,
		kind: r.nextComponentKind(),
		z:    structomancer.New(cmpt, "config"),
	}
}

func (r *ComponentRegistry) GetComponentType(typeName string) (CmptType, bool) {
	k, exists := r.componentTypes[typeName]
	return k, exists
}

func (c CmptType) Kind() ComponentKind {
	return c.kind
}

func (c CmptType) DeserializeConfig(cfg map[string]interface{}) (interface{}, error) {
	return c.z.MapToStruct(cfg)
}
