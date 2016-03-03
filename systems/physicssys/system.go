package physicssys

import (
	"fmt"

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

	//
	// apply acceleration / velocity
	//
	for _, e := range s.entities {
		accel := e.physicsCmpt.CurrentForces()
		e.physicsCmpt.ResetForces()

		vdelta := accel.Mul(float32(elapsed.Seconds()))

		newvel := e.physicsCmpt.Velocity().Add(vdelta)

		// friction
		// newvel = newvel.Mul(0.95)

		mag := newvel.Len()
		maxvel := e.physicsCmpt.MaxVelocity()
		if mag > 0 && maxvel < mag {
			newvel = newvel.Normalize().Mul(maxvel)
		}

		e.physicsCmpt.SetVelocity(newvel)

		// add the instantaneous velocity for the movement system
		newvel = newvel.Add(e.physicsCmpt.InstantaneousVelocity())

		newpos := e.positionCmpt.Pos().Add(newvel.Mul(float32(elapsed.Seconds())))
		e.positionCmpt.SetPos(newpos)
	}

	//
	// check for collisions
	//
	for _, entA := range s.entities {
		entitiesToCheck := s.entities // @@TODO

		for _, entB := range entitiesToCheck {
			if entA == entB {
				continue
			}
			did := s.checkCollision(entA, entB)
			if did {
				fmt.Println("collision")
			}
		}
	}

	s.previousTime = t
}

func getNormal(a, b mgl32.Vec2) mgl32.Vec2 {
	return mgl32.Vec2{
		-(a.Y() - b.Y()),
		a.X() - b.X(),
	}
}

func getMinMaxProjectedPoints(boundingBox BoundingBox, pos mgl32.Vec2, normal mgl32.Vec2) (float32, float32) {
	min := boundingBox[0].Add(pos).Dot(normal)
	max := min
	for j := 0; j < len(boundingBox); j++ {
		x := boundingBox[j].Add(pos).Dot(normal)
		if x > max {
			max = x
		} else if x < min {
			min = x
		}
	}

	return min, max
}

func (s *System) checkCollision(entA, entB entityAspect) bool {
	var minA, maxA, minB, maxB float32

	cmptA := entA.physicsCmpt
	cmptB := entB.physicsCmpt
	for i := 0; i < len(cmptA.boundingBox)-1; i++ {
		normal := getNormal(cmptA.boundingBox[i+1], cmptA.boundingBox[i])
		minA, maxA = getMinMaxProjectedPoints(cmptA.boundingBox, entA.positionCmpt.Pos(), normal)
		minB, maxB = getMinMaxProjectedPoints(cmptB.boundingBox, entB.positionCmpt.Pos(), normal)

		if maxB < minA || maxA < minB {
			// no collision between these shapes
			return false
		}
	}

	for i := 0; i < len(cmptB.boundingBox)-1; i++ {
		normal := getNormal(cmptB.boundingBox[i+1], cmptB.boundingBox[i])
		minA, maxA = getMinMaxProjectedPoints(cmptA.boundingBox, entA.positionCmpt.Pos(), normal)
		minB, maxB = getMinMaxProjectedPoints(cmptB.boundingBox, entB.positionCmpt.Pos(), normal)

		if maxB < minA || maxA < minB {
			// no collision between these shapes
			return false
		}
	}

	return true
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
