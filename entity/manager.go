package entity

import (
	"github.com/brynbellomy/gl4-game/common"
)

type (
	Manager struct {
		entities []Entity
		systems  []ISystem
	}

	Entity struct {
		ID         ID
		Components []IComponent
	}
)

func NewManager(systems []ISystem) Manager {
	return Manager{
		systems:  systems,
		entities: []Entity{},
	}
}

func (m *Manager) AddComponents(eid ID, components []IComponent) {
	m.entities = append(m.entities, Entity{ID: eid, Components: components})

	for _, sys := range m.systems {
		sys.ComponentsWillJoin(eid, components)
	}
}

func (m *Manager) Update(t common.Time) {
	for _, sys := range m.systems {
		sys.Update(t)
	}
}
