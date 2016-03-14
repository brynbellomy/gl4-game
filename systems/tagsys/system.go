package tagsys

import (
	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	System struct {
		entityManager  *entity.Manager
		componentQuery entity.ComponentMask
		tagCmptSet     entity.IComponentSet

		idsByTag map[string]entity.ID
	}
)

func New() *System {
	return &System{
		idsByTag: map[string]entity.ID{},
	}
}

func (s *System) EntityWithTag(tag string) (entity.ID, bool) {
	sl := s.tagCmptSet.Slice().(ComponentSlice)
	for idx, x := range sl {
		if tag == x.GetTag() {
			return s.tagCmptSet.IDForIndex(idx)
		}
	}
	return entity.InvalidID, false
}

func (s *System) Update(t common.Time) {
	// no-op
}

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"tag": {Component{}, ComponentSlice{}},
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	s.entityManager = em

	componentQuery, err := s.entityManager.MakeCmptQuery([]string{"tag"})
	if err != nil {
		panic(err)
	}
	s.componentQuery = componentQuery

	tagCmptSet, err := s.entityManager.GetComponentSet("tag")
	if err != nil {
		panic(err)
	}
	s.tagCmptSet = tagCmptSet
}

// func (s *System) EntityComponentsChanged(eid entity.ID, components []entity.IComponent) {
// 	var tagCmpt *Component

// 	for _, cmpt := range components {
// 		if ac, is := cmpt.(*Component); is {
// 			tagCmpt = ac
// 		}

// 		if tagCmpt != nil {
// 			break
// 		}
// 	}

// 	if tagCmpt != nil {
// 		if _, exists := s.entityMap[eid]; !exists {
// 			s.entities = append(s.entities, entityAspect{
// 				id:      eid,
// 				tagCmpt: tagCmpt,
// 			})

// 			s.idsByTag[tagCmpt.GetTag()] = eid
// 			s.entityMap[eid] = &s.entities[len(s.entities)-1]
// 		}

// 	} else {
// 		if _, exists := s.entityMap[eid]; exists {
// 			idx := -1
// 			for i := range s.entities {
// 				if s.entities[i].id == eid {
// 					idx = i
// 					break
// 				}
// 			}

// 			if idx >= 0 {
// 				s.entities = append(s.entities[:idx], s.entities[idx+1:]...)
// 			}

// 			delete(s.entityMap, eid)
// 		}
// 	}
// }
