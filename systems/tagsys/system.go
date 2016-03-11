package tagsys

import (
	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	System struct {
		entities  []entityAspect
		entityMap map[entity.ID]*entityAspect
		idsByTag  map[string]entity.ID
	}

	entityAspect struct {
		id      entity.ID
		tagCmpt *Component
	}
)

func New() *System {
	return &System{
		entities:  make([]entityAspect, 0),
		idsByTag:  map[string]entity.ID{},
		entityMap: make(map[entity.ID]*entityAspect),
	}
}

func (s *System) EntityWithTag(tag string) (entity.ID, bool) {
	eid, exists := s.idsByTag[tag]
	return eid, exists
}

func (s *System) Update(t common.Time) {
	// no-op
}

func (s *System) ComponentTypes() map[string]entity.IComponent {
	return map[string]entity.IComponent{
		"tag": &Component{},
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	// no-op
}

func (s *System) EntityComponentsChanged(eid entity.ID, components []entity.IComponent) {
	var tagCmpt *Component

	for _, cmpt := range components {
		if ac, is := cmpt.(*Component); is {
			tagCmpt = ac
		}

		if tagCmpt != nil {
			break
		}
	}

	if tagCmpt != nil {
		if _, exists := s.entityMap[eid]; !exists {
			s.entities = append(s.entities, entityAspect{
				id:      eid,
				tagCmpt: tagCmpt,
			})

			s.idsByTag[tagCmpt.GetTag()] = eid
			s.entityMap[eid] = &s.entities[len(s.entities)-1]
		}

	} else {
		if _, exists := s.entityMap[eid]; exists {
			idx := -1
			for i := range s.entities {
				if s.entities[i].id == eid {
					idx = i
					break
				}
			}

			if idx >= 0 {
				s.entities = append(s.entities[:idx], s.entities[idx+1:]...)
			}

			delete(s.entityMap, eid)
		}
	}
}
