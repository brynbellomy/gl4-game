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

		previousTime common.Time
	}

	entityAspect struct {
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
		return e.positionCmpt.Pos()
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

func (s *System) SetVelocity(eid entity.ID, vel mgl32.Vec2) {
	if e, exists := s.entityMap[eid]; exists {
		e.positionCmpt.SetVelocity(vel)
	} else {
		panic("entity does not exist")
	}
}

func (s *System) AddForce(eid entity.ID, f mgl32.Vec2) {
	if e, exists := s.entityMap[eid]; exists {
		e.positionCmpt.AddForce(f)
	} else {
		panic("entity does not exist")
	}
}

func (s *System) Update(t common.Time) {
	elapsed := t - s.previousTime

	for _, e := range s.entities {
		accel := e.positionCmpt.CurrentForces()
		e.positionCmpt.ResetForces()

		accelScaled := accel.Mul(float32(elapsed.Seconds()))

		newvel := e.positionCmpt.Velocity().Add(accelScaled)
		e.positionCmpt.SetVelocity(newvel)

		newvelScaled := newvel.Mul(float32(elapsed.Seconds()))

		newpos := e.positionCmpt.Pos().Add(newvelScaled)
		e.positionCmpt.SetPos(newpos)
	}

	s.previousTime = t
}

func (s *System) ComponentsWillJoin(eid entity.ID, components []entity.IComponent) {
	// if we find a *positionsys.Component on the entity, we keep track of it
	var positionCmpt *Component

	for _, cmpt := range components {
		if cmpt, is := cmpt.(*Component); is {
			positionCmpt = cmpt
			break
		}
	}

	if positionCmpt != nil {
		aspect := entityAspect{positionCmpt: positionCmpt}
		s.entities = append(s.entities, aspect)
		s.entityMap[eid] = &s.entities[len(s.entities)-1]
	}
}
