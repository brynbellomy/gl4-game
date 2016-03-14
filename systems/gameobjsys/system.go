package gameobjsys

import (
	"math"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/animationsys"
	"github.com/brynbellomy/gl4-game/systems/movesys"
)

type (
	System struct {
		entityManager    *entity.Manager
		componentQuery   entity.ComponentMask
		gameobjCmptSet   entity.IComponentSet
		animationCmptSet entity.IComponentSet
		moveCmptSet      entity.IComponentSet
	}
)

func New() *System {
	return &System{}
}

func (s *System) Update(t common.Time) {
	matchIDs := s.entityManager.EntitiesMatching(s.componentQuery)
	gameobjCmptIdxs, err := s.gameobjCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}
	animationCmptIdxs, err := s.animationCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}
	moveCmptIdxs, err := s.moveCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}

	gameobjCmptSlice := s.gameobjCmptSet.Slice().(ComponentSlice)
	animationCmptSlice := s.animationCmptSet.Slice().(animationsys.ComponentSlice)
	moveCmptSlice := s.moveCmptSet.Slice().(movesys.ComponentSlice)

	for i := 0; i < len(gameobjCmptIdxs); i++ {
		gameobjCmpt := gameobjCmptSlice[gameobjCmptIdxs[i]]
		animationCmpt := animationCmptSlice[animationCmptIdxs[i]]
		moveCmpt := moveCmptSlice[moveCmptIdxs[i]]

		vel := moveCmpt.Vector()

		if vel.Len() > 0 {
			radians := math.Atan2(float64(vel.Y()), float64(vel.X()))

			gameobjCmpt.Direction = DirectionFromRadians(radians)
			animationCmpt.SetIsAnimating(true)

		} else {
			animationCmpt.SetIsAnimating(false)
		}

		if action, exists := gameobjCmpt.Animations[gameobjCmpt.Action]; exists {
			anim := action[gameobjCmpt.Direction]
			if anim != "" {
				animationCmpt.SetAnimation(anim)
			}
		}

		gameobjCmptSlice[gameobjCmptIdxs[i]] = gameobjCmpt
		animationCmptSlice[animationCmptIdxs[i]] = animationCmpt
		moveCmptSlice[moveCmptIdxs[i]] = moveCmpt
	}
}

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"gameobj": {Component{}, ComponentSlice{}},
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	s.entityManager = em

	componentQuery, err := s.entityManager.MakeCmptQuery([]string{"gameobj", "animation", "move"})
	if err != nil {
		panic(err)
	}
	s.componentQuery = componentQuery

	gameobjCmptSet, err := s.entityManager.GetComponentSet("gameobj")
	if err != nil {
		panic(err)
	}
	s.gameobjCmptSet = gameobjCmptSet

	animationCmptSet, err := s.entityManager.GetComponentSet("animation")
	if err != nil {
		panic(err)
	}
	s.animationCmptSet = animationCmptSet

	moveCmptSet, err := s.entityManager.GetComponentSet("move")
	if err != nil {
		panic(err)
	}
	s.moveCmptSet = moveCmptSet
}

// func (s *System) EntityComponentsChanged(eid entity.ID, components []entity.IComponent) {
// 	var gameobjCmpt *Component
// 	var moveCmpt *movesys.Component
// 	var animationCmpt *animationsys.Component

// 	for _, cmpt := range components {
// 		if ac, is := cmpt.(*Component); is {
// 			gameobjCmpt = ac
// 		} else if rc, is := cmpt.(*animationsys.Component); is {
// 			animationCmpt = rc
// 		} else if pc, is := cmpt.(*movesys.Component); is {
// 			moveCmpt = pc
// 		}

// 		if gameobjCmpt != nil && animationCmpt != nil && moveCmpt != nil {
// 			break
// 		}
// 	}

// 	if gameobjCmpt != nil && animationCmpt != nil && moveCmpt != nil {
// 		if _, exists := s.entityMap[eid]; !exists {
// 			s.entities = append(s.entities, entityAspect{
// 				id:            eid,
// 				gameobjCmpt:   gameobjCmpt,
// 				animationCmpt: animationCmpt,
// 				moveCmpt:      moveCmpt,
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
// 		case *Component, *animationsys.Component, *movesys.Component:
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
