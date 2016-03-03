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
		vel := e.physicsCmpt.Velocity()

		if vel.Len() > 0.01 {
			radians := math.Atan2(float64(vel.Y()), float64(vel.X()))

			cmpt.direction = DirectionFromRadians(radians)
			e.animationCmpt.SetIsAnimating(true)

		} else {
			e.animationCmpt.SetIsAnimating(false)
		}

		if action, exists := cmpt.animations[cmpt.action]; exists {
			anim := action[cmpt.direction]
			if anim != "" {
				e.animationCmpt.SetAnimation(anim)
			}
		}
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	// no-op
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
			gameobjCmpt:   gameobjCmpt,
			animationCmpt: animationCmpt,
			physicsCmpt:   physicsCmpt,
		})

		s.entityMap[eid] = &s.entities[len(s.entities)-1]
	}
}
