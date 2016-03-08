package gameobjsys

import (
	"math"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/animationsys"
	"github.com/brynbellomy/gl4-game/systems/physicssys"
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
		physicsCmpt   *physicssys.Component
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
		vel := e.physicsCmpt.GetVelocity().Add(e.physicsCmpt.GetInstantaneousVelocity())

		if vel.Len() > 0.01 {
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

func (s *System) WillJoinManager(em *entity.Manager) {
	em.RegisterComponentType("gameobj", &Component{}, nil)
}

func (s *System) ComponentsWillJoin(eid entity.ID, components []entity.IComponent) {
	var gameobjCmpt *Component
	var physicsCmpt *physicssys.Component
	var animationCmpt *animationsys.Component

	for _, cmpt := range components {
		if ac, is := cmpt.(*Component); is {
			gameobjCmpt = ac
		} else if rc, is := cmpt.(*animationsys.Component); is {
			animationCmpt = rc
		} else if pc, is := cmpt.(*physicssys.Component); is {
			physicsCmpt = pc
		}

		if gameobjCmpt != nil && animationCmpt != nil && physicsCmpt != nil {
			break
		}
	}

	if gameobjCmpt != nil {
		if animationCmpt == nil {
			panic("gameobj component requires animation component")
		} else if physicsCmpt == nil {
			panic("gameobj component requires physics component")
		}

		s.entities = append(s.entities, entityAspect{
			id:            eid,
			gameobjCmpt:   gameobjCmpt,
			animationCmpt: animationCmpt,
			physicsCmpt:   physicsCmpt,
		})

		s.entityMap[eid] = &s.entities[len(s.entities)-1]
	}
}

func (s *System) ComponentsWillLeave(eid entity.ID, components []entity.IComponent) {
	remove := false
	for _, cmpt := range components {
		switch cmpt.(type) {
		case *Component, *animationsys.Component, *physicssys.Component:
			remove = true
			break
		}
	}

	if remove {
		removedIdx := -1
		for i := range s.entities {
			if s.entities[i].id == eid {
				removedIdx = i
				break
			}
		}
		if removedIdx >= 0 {
			s.entities = append(s.entities[:removedIdx], s.entities[removedIdx+1:]...)
		}
		delete(s.entityMap, eid)
	}
}
