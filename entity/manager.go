package entity

import (
	"github.com/brynbellomy/gl4-game/systems/assetsys"
)

type (
	Manager struct {
		entities      []Entity
		systems       []ISystem
		templateCache *TemplateCache
		entityFactory *EntityFactory

		cullable []ID

		idCounter ID
		usedIDs   map[ID]bool
	}

	Entity struct {
		ID         ID
		Components []IComponent
	}
)

func NewManager(fs assetsys.IFilesystem, systems []ISystem) *Manager {
	entityFactory := NewEntityFactory()
	templateCache := NewTemplateCache(fs, entityFactory)

	m := &Manager{
		systems:       systems,
		entities:      []Entity{},
		cullable:      []ID{},
		entityFactory: entityFactory,
		templateCache: templateCache,
		usedIDs:       map[ID]bool{},
	}

	for _, sys := range systems {
		sys.WillJoinManager(m)
	}

	return m
}

func (m *Manager) newEntityID() ID {
	for id := m.idCounter; ; id++ {
		if m.usedIDs[id] == false {
			m.usedIDs[id] = true
			m.idCounter = id + 1
			return id
		}
	}
}

func (m *Manager) setIDUsed(eid ID) {
	m.usedIDs[eid] = true
	if m.idCounter <= eid {
		m.idCounter = eid + 1
	}
}

func (m *Manager) RegisterComponentType(typeName string, cmpt IComponent) {
	m.entityFactory.RegisterComponentType(typeName, cmpt)
}

func (m *Manager) EntityFromTemplate(name string) (ID, []IComponent, error) {
	eid := m.newEntityID()

	ent, err := m.templateCache.Load(name)
	if err != nil {
		return 0, nil, err
	}

	return eid, ent, nil
}

func (m *Manager) EntityFromConfig(config map[string]interface{}) (ID, []IComponent, error) {
	return m.entityFactory.EntityFromConfig(config)
}

func (m *Manager) AddComponents(eid ID, components []IComponent) {
	m.entities = append(m.entities, Entity{ID: eid, Components: components})
	m.setIDUsed(eid)

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
