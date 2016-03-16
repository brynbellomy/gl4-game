package entity

import (
	"errors"
	"fmt"

	"github.com/brynbellomy/gl4-game/systems/assetsys"
)

type (
	Manager struct {
		systems       []ISystem
		entities      []entityRecord
		componentSets map[ComponentKind]IComponentSet

		templateCache *TemplateCache
		entityFactory *EntityFactory
		cmptRegistry  *ComponentRegistry

		cullable []ID

		idCounter ID
		usedIDs   map[ID]bool
	}

	entityRecord struct {
		id             ID
		mask           ComponentMask
		componentKinds []ComponentKind
	}
)

func NewManager(fs assetsys.IFilesystem, systems []ISystem) *Manager {
	cmptRegistry := NewComponentRegistry()
	entityFactory := NewEntityFactory(cmptRegistry)
	templateCache := NewTemplateCache(fs, entityFactory)

	componentSets := map[ComponentKind]IComponentSet{}
	for _, sys := range systems {
		ctypes := sys.ComponentTypes()
		for name, cmpt := range ctypes {
			kind := cmptRegistry.RegisterComponentType(name, cmpt.Cmpt)
			componentSets[kind] = NewComponentSet(cmpt.CmptSlice)
		}
	}

	m := &Manager{
		systems:       systems,
		entities:      []entityRecord{},
		cullable:      []ID{},
		usedIDs:       map[ID]bool{},
		entityFactory: entityFactory,
		templateCache: templateCache,
		cmptRegistry:  cmptRegistry,
		componentSets: componentSets,
	}

	for _, sys := range systems {
		sys.WillJoinManager(m)
	}

	return m
}

func (m *Manager) newEntityID() ID {
	for id := m.idCounter; ; id++ {
		if m.usedIDs[id] == false {
			// m.usedIDs[id] = true
			// m.idCounter = id + 1
			return id
		}
	}
}

func (m *Manager) setIDUsed(eid ID) {
	if m.usedIDs[eid] == true {
		panic("entity.Manager.setIDUsed: id already in use")
	}
	m.usedIDs[eid] = true
	if m.idCounter <= eid {
		m.idCounter = eid + 1
	}
}

func (m *Manager) GetComponentSet(name string) (IComponentSet, error) {
	cmptType, exists := m.cmptRegistry.GetComponentType(name)
	if !exists {
		return nil, errors.New("entity.Manager.GetComponentSet: unregistered component type '" + name + "'")
	}

	set, exists := m.componentSets[cmptType.Kind()]
	if !exists {
		return nil, errors.New("entity.Manager.GetComponentSet: unregistered component type '" + name + "'")
	}

	return set, nil
}

func (m *Manager) EntitiesMatching(cmptMask ComponentMask) []ID {
	matching := make([]ID, 0)
	for _, ent := range m.entities {
		if ent.mask.HasAll(cmptMask) {
			matching = append(matching, ent.id)
		}
	}
	return matching
}

func (m *Manager) MakeCmptQuery(cmptTypes []string) (ComponentMask, error) {
	var cmptQuery ComponentMask
	for _, typeName := range cmptTypes {
		t, exists := m.cmptRegistry.GetComponentType(typeName)
		if !exists {
			return 0, errors.New("entity.Manager.MakeCmptQuery: component type '" + typeName + "' is not registered")
		}

		cmptQuery = cmptQuery.Add(t.Kind())
	}
	return cmptQuery, nil
}

func (m *Manager) EntityFromTemplate(name string) (Entity, error) {
	ent, err := m.templateCache.Load(name)
	if err != nil {
		return Entity{}, err
	}

	ent.ID = m.newEntityID()
	m.setIDUsed(ent.ID)

	return ent, nil
}

func (m *Manager) EntityFromConfig(config map[string]interface{}) (Entity, error) {
	entity, err := m.entityFactory.EntityFromConfig(config)
	if err != nil {
		return Entity{}, err
	}

	if m.usedIDs[entity.ID] == true {
		return Entity{}, fmt.Errorf("entity.Manager: entity ID '%s' is already in use", entity.ID)
	}
	m.setIDUsed(entity.ID)

	return entity, nil
}

func (m *Manager) AddEntity(entity Entity) {
	m.entities = append(m.entities, entityRecord{id: entity.ID, mask: entity.ComponentMask, componentKinds: entity.ComponentKinds})

	for i, cmpt := range entity.Components {
		kind := entity.ComponentKinds[i]
		m.componentSets[kind].Add(entity.ID, cmpt)
	}
}

func (m *Manager) RemoveEntity(eid ID) {
	m.cullable = append(m.cullable, eid)
}

func (m *Manager) CullEntities() {
	for _, eid := range m.cullable {
		removedIdx := -1
		for i := range m.entities {
			if m.entities[i].id == eid {
				removedIdx = i
				break
			}
		}

		if removedIdx >= 0 {
			ent := m.entities[removedIdx]
			for _, kind := range ent.componentKinds {
				m.componentSets[kind].Remove(eid)
			}

			m.entities = append(m.entities[:removedIdx], m.entities[removedIdx+1:]...)
		}
	}
	m.cullable = []ID{}
}
