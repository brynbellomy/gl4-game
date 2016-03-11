package entity

import (
	"errors"
	"fmt"

	"github.com/brynbellomy/gl4-game/systems/assetsys"
)

type (
	Manager struct {
		systems  []ISystem
		entities []Entity

		templateCache *TemplateCache
		entityFactory *EntityFactory
		cmptRegistry  *ComponentRegistry

		cullable []ID

		idCounter ID
		usedIDs   map[ID]bool
	}

	Entity struct {
		ID           ID
		ComponentSet ComponentSet
		Components   []IComponent
	}
)

func NewManager(fs assetsys.IFilesystem, systems []ISystem) *Manager {
	cmptRegistry := NewComponentRegistry()
	entityFactory := NewEntityFactory(cmptRegistry)
	templateCache := NewTemplateCache(fs, entityFactory)

	m := &Manager{
		systems:       systems,
		entities:      []Entity{},
		cullable:      []ID{},
		usedIDs:       map[ID]bool{},
		entityFactory: entityFactory,
		templateCache: templateCache,
		cmptRegistry:  cmptRegistry,
	}

	for _, sys := range systems {
		ctypes := sys.ComponentTypes()
		for name, cmpt := range ctypes {
			cmptRegistry.RegisterComponentType(name, cmpt)
		}
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

func (m *Manager) MakeCmptQuery(cmptTypes []string) (ComponentSet, error) {
	var cmptQuery ComponentSet
	for _, typeName := range cmptTypes {
		t, exists := m.cmptRegistry.GetComponentType(typeName)
		if !exists {
			return 0, errors.New("entity.Manager.MakeCmptQuery: component type '" + typeName + "' is not registered")
		}

		cmptQuery = cmptQuery.Add(t.Kind())
	}
	return cmptQuery, nil
}

// func (m *Manager) RegisterComponentType(typeName string, cmpt IComponent) {
// 	m.cmptRegistry.RegisterComponentType(typeName, cmpt)
// }

func (m *Manager) EntityFromTemplate(name string) (ID, []IComponent, error) {
	eid := m.newEntityID()

	ent, err := m.templateCache.Load(name)
	if err != nil {
		return 0, nil, err
	}

	return eid, ent, nil
}

func (m *Manager) EntityFromConfig(config map[string]interface{}) (ID, []IComponent, error) {
	eid, cmpts, err := m.entityFactory.EntityFromConfig(config)
	if err != nil {
		return 0, nil, err
	}

	if m.usedIDs[eid] == true {
		return 0, nil, fmt.Errorf("entity.Manager: entity ID '%s' is already in use", eid)
	}

	m.setIDUsed(eid)

	return eid, cmpts, nil
}

func (m *Manager) SetComponents(eid ID, components []IComponent) {
	var mask ComponentSet
	for _, cmpt := range components {
		mask = mask.Add(cmpt.Kind())
	}

	ent := Entity{ID: eid, ComponentSet: mask, Components: components}

	m.entities = append(m.entities, ent)
	m.setIDUsed(eid)

	for _, sys := range m.systems {
		sys.EntityComponentsChanged(eid, components)
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

		if removedIdx >= 0 {
			m.entities = append(m.entities[:removedIdx], m.entities[removedIdx+1:]...)

			for _, sys := range m.systems {
				sys.EntityComponentsChanged(eid, []IComponent{})
			}
		}
	}
	m.cullable = []ID{}
}
