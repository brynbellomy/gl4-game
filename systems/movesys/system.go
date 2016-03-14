package movesys

import (
	// "github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/physicssys"
)

type (
	System struct {
		entityManager  *entity.Manager
		componentQuery entity.ComponentMask

		moveCmptSet    entity.IComponentSet
		physicsCmptSet entity.IComponentSet
	}
)

func New() *System {
	return &System{}
}

func (s *System) Update(t common.Time) {
	matchIDs := s.entityManager.EntitiesMatching(s.componentQuery)

	moveCmptIdxs, err := s.moveCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}
	physCmptIdxs, err := s.physicsCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}

	moveCmptSlice := s.moveCmptSet.Slice().(ComponentSlice)
	physicsCmptSlice := s.physicsCmptSet.Slice().(physicssys.ComponentSlice)

	for i := 0; i < len(moveCmptIdxs); i++ {
		vec := moveCmptSlice[moveCmptIdxs[i]].Vector()
		physicsCmptSlice[physCmptIdxs[i]].SetInstantaneousVelocity(vec)
	}
}

// func (s *System) SetMovementVector(eid entity.ID, vec mgl32.Vec2) {
// 	if e, exists := s.entityMap[eid]; exists {
// 		e.moveCmpt.SetVector(vec)
// 	} else {
// 		panic("entity does not exist")
// 	}
// }

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"move": {Component{}, ComponentSlice{}},
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	s.entityManager = em

	componentQuery, err := em.MakeCmptQuery([]string{"move", "physics"})
	if err != nil {
		panic(err)
	}
	s.componentQuery = componentQuery

	moveCmptSet, err := em.GetComponentSet("move")
	if err != nil {
		panic(err)
	}
	s.moveCmptSet = moveCmptSet

	physicsCmptSet, err := em.GetComponentSet("physics")
	if err != nil {
		panic(err)
	}
	s.physicsCmptSet = physicsCmptSet
}

// func (s *System) EntityComponentsChanged(eid entity.ID, components []entity.IComponent) {
// 	var moveCmpt *Component
// 	var physicsCmpt *physicssys.Component

// 	for _, cmpt := range components {
// 		if ac, is := cmpt.(*Component); is {
// 			moveCmpt = ac
// 		} else if pc, is := cmpt.(*physicssys.Component); is {
// 			physicsCmpt = pc
// 		}

// 		if moveCmpt != nil && physicsCmpt != nil {
// 			break
// 		}
// 	}

// 	if moveCmpt != nil && physicsCmpt != nil {
// 		if _, exists := s.entityMap[eid]; !exists {
// 			s.entities = append(s.entities, entityAspect{
// 				id:          eid,
// 				moveCmpt:    moveCmpt,
// 				physicsCmpt: physicsCmpt,
// 			})

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
