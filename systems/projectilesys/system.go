package projectilesys

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/physicssys"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
)

type (
	System struct {
		entityManager  *entity.Manager
		componentQuery entity.ComponentMask

		positionCmptSet   entity.IComponentSet
		physicsCmptSet    entity.IComponentSet
		projectileCmptSet entity.IComponentSet
	}
)

func New() *System {
	return &System{}
}

// func (s *System) GetHeading(eid entity.ID) mgl32.Vec2 {
// 	if e, exists := s.entityMap[eid]; exists {
// 		return e.projectileCmpt.GetHeading()
// 	} else {
// 		panic("entity does not exist")
// 	}
// }

// func (s *System) SetHeading(eid entity.ID, pos mgl32.Vec2) {
// 	if e, exists := s.entityMap[eid]; exists {
// 		e.projectileCmpt.SetHeading(pos)
// 	} else {
// 		panic("entity does not exist")
// 	}
// }

func (s *System) Update(t common.Time) {
	matchIDs := s.entityManager.EntitiesMatching(s.componentQuery)
	positionCmptIdxs, err := s.positionCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}
	physicsCmptIdxs, err := s.physicsCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}
	projectileCmptIdxs, err := s.projectileCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}

	positionCmptSlice := s.positionCmptSet.Slice().(positionsys.ComponentSlice)
	physicsCmptSlice := s.physicsCmptSet.Slice().(physicssys.ComponentSlice)
	projectileCmptSlice := s.projectileCmptSet.Slice().(ComponentSlice)

	// check for collisions first
	for i := 0; i < len(projectileCmptIdxs); i++ {
		if len(physicsCmptSlice[physicsCmptIdxs[i]].GetCollisions()) > 0 {
			projectileCmptSlice[projectileCmptIdxs[i]].State = Impacting
		}
	}

	for i := 0; i < len(projectileCmptIdxs); i++ {
		projectileCmpt := projectileCmptSlice[projectileCmptIdxs[i]]
		physicsCmpt := physicsCmptSlice[physicsCmptIdxs[i]]
		positionCmpt := positionCmptSlice[positionCmptIdxs[i]]

		switch projectileCmpt.State {
		case Firing:
			headingNorm := projectileCmpt.GetHeading().Normalize()
			physicsCmpt.SetVelocity(headingNorm.Mul(projectileCmpt.ExitVelocity))
			physicsCmpt.AddForce(headingNorm.Mul(projectileCmpt.Thrust))

			// only stay in the Firing state for the first frame
			projectileCmpt.State = Flying

		case Flying:
			force := projectileCmpt.GetHeading().Normalize().Mul(projectileCmpt.Thrust)
			physicsCmpt.AddForce(force)

			v := physicsCmpt.GetVelocity()
			theta := float32(math.Atan2(float64(v.Y()), float64(v.X())))
			positionCmpt.SetRotation(theta)

		case Impacting:
			physicsCmpt.SetVelocity(mgl32.Vec2{0, 0})
			if projectileCmpt.RemoveOnContact {
				s.entityManager.RemoveEntity(matchIDs[i])
			}
		}

		projectileCmptSlice[projectileCmptIdxs[i]] = projectileCmpt
		physicsCmptSlice[physicsCmptIdxs[i]] = physicsCmpt
		positionCmptSlice[positionCmptIdxs[i]] = positionCmpt
	}
}

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"projectile": {Component{}, ComponentSlice{}},
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	s.entityManager = em

	componentQuery, err := s.entityManager.MakeCmptQuery([]string{"position", "physics", "projectile"})
	if err != nil {
		panic(err)
	}
	s.componentQuery = componentQuery

	projectileCmptSet, err := s.entityManager.GetComponentSet("projectile")
	if err != nil {
		panic(err)
	}
	s.projectileCmptSet = projectileCmptSet

	positionCmptSet, err := s.entityManager.GetComponentSet("position")
	if err != nil {
		panic(err)
	}
	s.positionCmptSet = positionCmptSet

	physicsCmptSet, err := s.entityManager.GetComponentSet("physics")
	if err != nil {
		panic(err)
	}
	s.physicsCmptSet = physicsCmptSet
}

// func (s *System) EntityComponentsChanged(eid entity.ID, components []entity.IComponent) {
// 	var projectileCmpt *Component
// 	var moveCmpt *movesys.Component
// 	var positionCmpt *positionsys.Component
// 	var physicsCmpt *physicssys.Component

// 	for _, cmpt := range components {
// 		if ac, is := cmpt.(*Component); is {
// 			projectileCmpt = ac
// 		} else if rc, is := cmpt.(*positionsys.Component); is {
// 			positionCmpt = rc
// 		} else if pc, is := cmpt.(*movesys.Component); is {
// 			moveCmpt = pc
// 		} else if physc, is := cmpt.(*physicssys.Component); is {
// 			physicsCmpt = physc
// 		}

// 		if projectileCmpt != nil && positionCmpt != nil && moveCmpt != nil {
// 			break
// 		}
// 	}

// 	if projectileCmpt != nil && positionCmpt != nil && moveCmpt != nil {
// 		if _, exists := s.entityMap[eid]; !exists {
// 			s.entities = append(s.entities, entityAspect{
// 				id:             eid,
// 				projectileCmpt: projectileCmpt,
// 				positionCmpt:   positionCmpt,
// 				moveCmpt:       moveCmpt,
// 				physicsCmpt:    physicsCmpt,
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
// 		case *Component, *movesys.Component, *positionsys.Component, *physicssys.Component:
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
