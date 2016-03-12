package movesys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/physicssys"
)

type (
	System struct {
		// entities  []entityAspect
		// entityMap map[entity.ID]*entityAspect
		entityManager  *entity.Manager
		componentQuery entity.ComponentMask
	}

	entityAspect struct {
		id          entity.ID
		moveCmpt    *Component
		physicsCmpt *physicssys.Component
	}
)

func New() *System {
	return &System{
	// entities:  make([]entityAspect, 0),
	// entityMap: make(map[entity.ID]*entityAspect),
	}
}

func (s *System) Update(t common.Time) {
	matchIDs := s.entityManager.EntitiesMatching(s.componentQuery)
	moveCmpts := s.entityManager.GetComponentSet("move").Visitor(matchIDs)
	physicsCmpts := s.entityManager.GetComponentSet("physics").Visitor(matchIDs)

	for i := 0; i < moveCmpts.Len(); i++ {
		moveCmpt := moveCmpts.Get().(*Component)
		physicsCmpt := physicsCmpts.Get().(*physicssys.Component)

		physicsCmpt.SetInstantaneousVelocity(moveCmpt.Vector())

		moveCmpts.Set(moveCmpt)
		physicsCmpts.Set(physicsCmpt)

		moveCmpts.Advance()
		physicsCmpts.Advance()
	}

	// for _, e := range s.entities {
	// 	vec := e.moveCmpt.Vector()
	// 	e.physicsCmpt.SetInstantaneousVelocity(vec)
	// 	// e.moveCmpt.ResetVector()
	// }
}

func (s *System) SetMovementVector(eid entity.ID, vec mgl32.Vec2) {
	if e, exists := s.entityMap[eid]; exists {
		e.moveCmpt.SetVector(vec)
	} else {
		panic("entity does not exist")
	}
}

func (s *System) ComponentTypes() map[string]entity.IComponent {
	return map[string]entity.IComponent{
		"move": &Component{},
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	// em.RegisterComponentType("move", &Component{})
	s.entityManager = em
	s.componentQuery = em.MakeCmptQuery([]string{"move", "physics"})
}

func (s *System) EntityComponentsChanged(eid entity.ID, components []entity.IComponent) {
	var moveCmpt *Component
	var physicsCmpt *physicssys.Component

	for _, cmpt := range components {
		if ac, is := cmpt.(*Component); is {
			moveCmpt = ac
		} else if pc, is := cmpt.(*physicssys.Component); is {
			physicsCmpt = pc
		}

		if moveCmpt != nil && physicsCmpt != nil {
			break
		}
	}

	if moveCmpt != nil && physicsCmpt != nil {
		if _, exists := s.entityMap[eid]; !exists {
			s.entities = append(s.entities, entityAspect{
				id:          eid,
				moveCmpt:    moveCmpt,
				physicsCmpt: physicsCmpt,
			})

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
// 		case *Component, *positionsys.Component, *physicssys.Component:
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
