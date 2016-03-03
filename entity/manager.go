package entity

type (
	Manager struct {
		entities []Entity
		systems  []ISystem

		idCounter ID
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

func (m *Manager) NewEntityID() ID {
	cur := m.idCounter
	m.idCounter++
	return cur
}

func (m *Manager) AddComponents(eid ID, components []IComponent) {
	m.entities = append(m.entities, Entity{ID: eid, Components: components})

	for _, sys := range m.systems {
		sys.ComponentsWillJoin(eid, components)
	}
}
