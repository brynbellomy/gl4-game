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
		entities  []entityAspect
		entityMap map[entity.ID]*entityAspect
	}

	entityAspect struct {
		id            entity.ID
		gameobjCmpt   *Component
		animationCmpt *animationsys.Component
		moveCmpt      *movesys.Component
	}
)

func New() *System {
	return &System{
		entities:  make([]entityAspect, 0),
		entityMap: make(map[entity.ID]*entityAspect),
	}
}

func (s *System) Update(t common.Time) {
	for _, e := range s.entities {
		cmpt := e.gameobjCmpt
		// vel := e.moveCmpt.GetVelocity().Add(e.moveCmpt.GetInstantaneousVelocity())
		vel := e.moveCmpt.Vector()

		if vel.Len() > 0 {
			radians := math.Atan2(float64(vel.Y()), float64(vel.X()))

			cmpt.Direction = DirectionFromRadians(radians)
			e.animationCmpt.SetIsAnimating(true)

		} else {
			e.animationCmpt.SetIsAnimating(false)
		}

		if action, exists := cmpt.Animations[cmpt.Action]; exists {
			anim := action[cmpt.Direction]
			if anim != "" {
				e.animationCmpt.SetAnimation(anim)
			}
		}
	}
}

func (s *System) ComponentTypes() map[string]entity.IComponent {
	return map[string]entity.IComponent{
		"gameobj": &Component{},
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	// em.RegisterComponentType("gameobj", &Component{})
}

func (s *System) EntityComponentsChanged(eid entity.ID, components []entity.IComponent) {
	var gameobjCmpt *Component
	var moveCmpt *movesys.Component
	var animationCmpt *animationsys.Component

	for _, cmpt := range components {
		if ac, is := cmpt.(*Component); is {
			gameobjCmpt = ac
		} else if rc, is := cmpt.(*animationsys.Component); is {
			animationCmpt = rc
		} else if pc, is := cmpt.(*movesys.Component); is {
			moveCmpt = pc
		}

		if gameobjCmpt != nil && animationCmpt != nil && moveCmpt != nil {
			break
		}
	}

	if gameobjCmpt != nil && animationCmpt != nil && moveCmpt != nil {
		if _, exists := s.entityMap[eid]; !exists {
			s.entities = append(s.entities, entityAspect{
				id:            eid,
				gameobjCmpt:   gameobjCmpt,
				animationCmpt: animationCmpt,
				moveCmpt:      moveCmpt,
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
