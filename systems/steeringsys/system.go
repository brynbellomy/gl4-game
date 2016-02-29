package steeringsys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/context"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
)

type (
	System struct {
		entities  []entityAspect
		entityMap map[entity.ID]*entityAspect
	}

	entityAspect struct {
		steeringCmpt *Component
		positionCmpt *positionsys.Component
	}

	IBehavior interface {
		SteeringVector() mgl32.Vec2
	}
)

func New() *System {
	return &System{
		entities:  []entityAspect{},
		entityMap: map[entity.ID]*entityAspect{},
	}
}

func (s *System) AddBehavior(eid entity.ID, behavior IBehavior) {
	if ent, exists := s.entityMap[eid]; exists {
		ent.steeringCmpt.behaviors = append(ent.steeringCmpt.behaviors, behavior)
	} else {
		panic("entity does not exist")
	}
}

func (s *System) Update(c context.IContext) {
	for _, ent := range s.entities {
		var currentSteering mgl32.Vec2

		for _, behavior := range ent.steeringCmpt.behaviors {
			steeringVec := behavior.SteeringVector()
			currentSteering = currentSteering.Add(steeringVec)
		}

		pos := ent.positionCmpt.Pos()
		ent.positionCmpt.SetPos(pos.Add(currentSteering))
	}
}

func (s *System) ComponentsWillJoin(eid entity.ID, components []entity.IComponent) {
	var positionCmpt *positionsys.Component
	var steeringCmpt *Component

	for _, cmpt := range components {
		if rc, is := cmpt.(*Component); is {
			steeringCmpt = rc
		} else if pc, is := cmpt.(*positionsys.Component); is {
			positionCmpt = pc
		}

		if positionCmpt != nil && steeringCmpt != nil {
			break
		}
	}

	if steeringCmpt != nil {
		if positionCmpt == nil {
			panic("steering component requires position component")
		}

		aspect := entityAspect{
			positionCmpt: positionCmpt,
			steeringCmpt: steeringCmpt,
		}

		s.entities = append(s.entities, aspect)
		s.entityMap[eid] = &aspect
	}
}
