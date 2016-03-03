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
	if s.previousTime == 0 {
		s.previousTime = t
		return
	}

	elapsed := t - s.previousTime

	for _, e := range s.entities {
		accel := e.physicsCmpt.CurrentForces()
		e.physicsCmpt.ResetForces()

		vdelta := accel.Mul(float32(elapsed.Seconds()))

		newvel := e.physicsCmpt.Velocity().Add(vdelta)

		// friction
		newvel = newvel.Mul(0.95)

		mag := newvel.Len()
		maxvel := e.physicsCmpt.MaxVelocity()
		if mag > 0 && maxvel < mag {
			newvel = newvel.Normalize().Mul(maxvel)
		}

		e.physicsCmpt.SetVelocity(newvel)

		newpos := e.positionCmpt.Pos().Add(newvel.Mul(float32(elapsed.Seconds())))
		e.positionCmpt.SetPos(newpos)
	}

	s.previousTime = t
}

func (s *System) WillJoinManager(em *entity.Manager) {
	// no-op
}

func (s *System) ComponentsWillJoin(eid entity.ID, components []entity.IComponent) {
	var physicsCmpt *Component
	var positionCmpt *positionsys.Component

	for _, cmpt := range components {
		if phc, is := cmpt.(*Component); is {
			physicsCmpt = phc
		} else if posc, is := cmpt.(*positionsys.Component); is {
			positionCmpt = posc
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
