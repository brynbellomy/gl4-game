package positionsys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	System struct {
		entities  []entityAspect
		entityMap map[entity.ID]*entityAspect
	}

	entityAspect struct {
		id           entity.ID
		positionCmpt *Component
	}
)

func New() *System {
	return &System{
		entities:  make([]entityAspect, 0),
		entityMap: make(map[entity.ID]*entityAspect),
	}
}

func (s *System) GetPos(eid entity.ID) mgl32.Vec2 {
	if e, exists := s.entityMap[eid]; exists {
		return e.positionCmpt.GetPos()
	} else {
		panic("entity does not exist")
	}
}

func (s *System) SetPos(eid entity.ID, pos mgl32.Vec2) {
	if e, exists := s.entityMap[eid]; exists {
		e.positionCmpt.SetPos(pos)
	} else {
		panic("entity does not exist")
	}
}

func (s *System) Update(t common.Time) {
	// no-op
}

func (s *System) ComponentTypes() map[string]entity.IComponent {
	return map[string]entity.IComponent{
		"position": &Component{},
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	// 	em.RegisterComponentType("position", &Component{})
}

func (s *System) EntityComponentsChanged(eid entity.ID, components []entity.IComponent) {
	var positionCmpt *Component

	for _, cmpt := range components {
		if cmpt, is := cmpt.(*Component); is {
			positionCmpt = cmpt
			break
		}
	}

	if positionCmpt != nil {
		if _, exists := s.entityMap[eid]; !exists {
			aspect := entityAspect{id: eid, positionCmpt: positionCmpt}
			s.entities = append(s.entities, aspect)
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

// func (s *System) ComponentsWillLeave(eid entity.ID, components []entity.IComponent) {
// 	remove := false
// 	for _, cmpt := range components {
// 		switch cmpt.(type) {
// 		case *Component:
// 			remove = true
// 			break
// 		}
// 	}

// 	if remove {
// 		removedIdx := -1
// 		for i := range s.entities {
// 			if s.entities[i].id == eid {
// 				removedIdx = i
// 				break
// 			}
// 		}
// 		if removedIdx >= 0 {
// 			s.entities = append(s.entities[:removedIdx], s.entities[removedIdx+1:]...)
// 		}
// 		delete(s.entityMap, eid)
// 	}
// }
