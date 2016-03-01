package physicssys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
)

type (
	System struct {
		entities  []entityAspect
		entityMap map[entity.ID]*entityAspect

		previousTime common.Time
	}

	entityAspect struct {
		physicsCmpt  *Component
		positionCmpt *positionsys.Component
	}
)

func New() *System {
	return &System{
		entities:  make([]entityAspect, 0),
		entityMap: make(map[entity.ID]*entityAspect),
	}
}

func (s *System) AddForce(eid entity.ID, f mgl32.Vec2) {
	if e, exists := s.entityMap[eid]; exists {
		e.physicsCmpt.AddForce(f)
	} else {
		panic("entity does not exist")
	}
}

func (s *System) Update(t common.Time) {
	elapsed := t - s.previousTime

	for _, e := range s.entities {
		accel := e.physicsCmpt.CurrentForces()
		e.physicsCmpt.ResetForces()

		accelScaled := accel.Mul(float32(elapsed.Seconds()))

		newvel := e.physicsCmpt.Velocity().Add(accelScaled)

		// friction
		newvel = newvel.Mul(0.95)

		e.physicsCmpt.SetVelocity(newvel)

		newpos := e.positionCmpt.Pos().Add(newvel.Mul(float32(elapsed.Seconds())))
		e.positionCmpt.SetPos(newpos)
	}

	s.previousTime = t
}

func (s *System) ComponentsWillJoin(eid entity.ID, components []entity.IComponent) {
	var physicsCmpt *Component
	var positionCmpt *positionsys.Component

	for _, cmpt := range components {
		if ac, is := cmpt.(*Component); is {
			physicsCmpt = ac
		} else if rc, is := cmpt.(*positionsys.Component); is {
			positionCmpt = rc
		}

		if physicsCmpt != nil && positionCmpt != nil {
			break
		}
	}

	if physicsCmpt != nil {
		if positionCmpt == nil {
			panic("physics component requires position component")
		}

		s.entities = append(s.entities, entityAspect{
			physicsCmpt:  physicsCmpt,
			positionCmpt: positionCmpt,
		})

		s.entityMap[eid] = &s.entities[len(s.entities)-1]
	}
}
