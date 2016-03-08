package entity

type (
	Manager struct {
		entities []Entity
		systems  []ISystem
		factory  *Factory

		cullable []ID

		idCounter ID
	}

	Entity struct {
		ID         ID
		Components []IComponent
	}
)

func NewManager(systems []ISystem) *Manager {
	m := &Manager{
		systems:  systems,
		entities: []Entity{},
		cullable: []ID{},
		factory:  NewFactory(),
	}

	for _, sys := range systems {
		sys.WillJoinManager(m)
	}

	return m
}

func (m *Manager) NewEntityID() ID {
	cur := m.idCounter
	m.idCounter++
	return cur
}

func (m *Manager) RegisterEntityTemplate(name string, cmpts []IComponent) {
	m.factory.RegisterEntityTemplate(name, cmpts)
}

func (m *Manager) RegisterComponentType(typeName string, cmpt IComponent, initer IComponentIniter) {
	m.factory.RegisterComponentType(typeName, cmpt, initer)
}

func (m *Manager) EntityFromConfig(config map[string]interface{}) (ID, []IComponent, error) {
	return m.factory.EntityFromConfig(config)
}

// func (m *Manager) AddEntityFromConfig(config map[string]interface{}) error {
// 	eid, cmpts, err := m.factory.EntityFromConfig(config)
// 	if err != nil {
// 		return err
// 	}

// 	m.AddComponents(eid, cmpts)

// 	return nil
// }

func (m *Manager) AddComponents(eid ID, components []IComponent) {
	m.entities = append(m.entities, Entity{ID: eid, Components: components})

	for _, sys := range m.systems {
		sys.ComponentsWillJoin(eid, components)
	}
}

func (m *Manager) RemoveEntity(eid ID) {
	m.cullable = append(m.cullable, eid)
}

func (m *Manager) CullEntities() {
	for _, eid := range m.cullable {
		removedIdx := -1
		for i := range m.entities {
			if m.entities[i].ID == eid {
				removedIdx = i
				break
			}
		}

		var cmpts []IComponent
		if removedIdx >= 0 {
			cmpts = m.entities[removedIdx].Components
			m.entities = append(m.entities[:removedIdx], m.entities[removedIdx+1:]...)

			for _, sys := range m.systems {
				sys.ComponentsWillLeave(eid, cmpts)
			}
		}
	}
	m.cullable = []ID{}
}
