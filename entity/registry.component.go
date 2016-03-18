package entity

import (
	"fmt"

	"github.com/brynbellomy/gl4-game/common"
)

type (
	ComponentRegistry struct {
		componentTypes map[string]CmptRegistryEntry
		kindCounter    uint64
	}

	CmptRegistryEntry struct {
		Name  string
		Kind  ComponentKind
		Coder *common.Coder
	}
)

func NewComponentRegistry() *ComponentRegistry {
	return &ComponentRegistry{
		componentTypes: map[string]CmptRegistryEntry{},
		kindCounter:    0,
	}
}

func (r *ComponentRegistry) nextComponentKind() ComponentKind {
	next := 1 << r.kindCounter
	r.kindCounter++
	if r.kindCounter == 64 {
		panic("entity.ComponentFactory: maximum number of component kinds exceeded (64)")
	}
	return ComponentKind(next)
}

func (r *ComponentRegistry) RegisterComponentType(typeName string, coder *common.Coder) ComponentKind {
	kind := r.nextComponentKind()

	r.componentTypes[typeName] = CmptRegistryEntry{
		Name:  typeName,
		Kind:  kind,
		Coder: coder,
	}

	fmt.Printf("Registering component type '%s' (kind = %d)\n", typeName, kind)

	return kind
}

func (r *ComponentRegistry) GetComponentType(typeName string) (CmptRegistryEntry, bool) {
	k, exists := r.componentTypes[typeName]
	return k, exists
}
