package projectilesys

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/movesys"
	"github.com/brynbellomy/gl4-game/systems/physicssys"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
)

type (
	System struct {
		entities      []entityAspect
		entityMap     map[entity.ID]*entityAspect
		entityManager *entity.Manager
	}

	entityAspect struct {
		id             entity.ID
		projectileCmpt *Component
		positionCmpt   *positionsys.Component
		moveCmpt       *movesys.Component
		physicsCmpt    *physicssys.Component
	}
)

func New() *System {
	return &System{
		entities:      make([]entityAspect, 0),
		entityMap:     make(map[entity.ID]*entityAspect),
		entityManager: nil,
	}
}

func (s *System) GetHeading(eid entity.ID) mgl32.Vec2 {
	if e, exists := s.entityMap[eid]; exists {
		return e.projectileCmpt.GetHeading()
	} else {
		panic("entity does not exist")
	}
}

func (s *System) SetHeading(eid entity.ID, pos mgl32.Vec2) {
	if e, exists := s.entityMap[eid]; exists {
		e.projectileCmpt.SetHeading(pos)
	} else {
		panic("entity does not exist")
	}
}

func (s *System) Update(t common.Time) {
	// check for collisions first
	for _, ent := range s.entities {
		if len(ent.physicsCmpt.GetCollisions()) > 0 {
			ent.projectileCmpt.State = Impacting
		}
	}

	for _, ent := range s.entities {
		switch ent.projectileCmpt.State {
		case Firing:
			headingNorm := ent.projectileCmpt.GetHeading().Normalize()
			ent.physicsCmpt.SetVelocity(headingNorm.Mul(ent.projectileCmpt.ExitVelocity))
			ent.physicsCmpt.AddForce(headingNorm.Mul(ent.projectileCmpt.Thrust))

			// only stay in the Firing state for the first frame
			ent.projectileCmpt.State = Flying

		case Flying:
			force := ent.projectileCmpt.GetHeading().Normalize().Mul(ent.projectileCmpt.Thrust)
			ent.physicsCmpt.AddForce(force)

			v := ent.physicsCmpt.GetVelocity()
			theta := float32(math.Atan2(float64(v.Y()), float64(v.X())))
			ent.positionCmpt.SetRotation(theta)

		case Impacting:
			ent.physicsCmpt.SetVelocity(mgl32.Vec2{0, 0})
			if ent.projectileCmpt.RemoveOnContact {
				s.entityManager.RemoveEntity(ent.id)
			}
		}
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	s.entityManager = em
	em.RegisterComponentType("projectile", &Component{})
}

func (s *System) ComponentsWillJoin(eid entity.ID, components []entity.IComponent) {
	var projectileCmpt *Component
	var moveCmpt *movesys.Component
	var positionCmpt *positionsys.Component
	var physicsCmpt *physicssys.Component

	for _, cmpt := range components {
		if ac, is := cmpt.(*Component); is {
			projectileCmpt = ac
		} else if rc, is := cmpt.(*positionsys.Component); is {
			positionCmpt = rc
		} else if pc, is := cmpt.(*movesys.Component); is {
			moveCmpt = pc
		} else if physc, is := cmpt.(*physicssys.Component); is {
			physicsCmpt = physc
		}

		if projectileCmpt != nil && positionCmpt != nil && moveCmpt != nil {
			break
		}
	}

	if projectileCmpt != nil {
		if positionCmpt == nil {
			panic("projectile component requires position component")
		} else if moveCmpt == nil {
			panic("projectile component requires move component")
		} else if physicsCmpt == nil {
			panic("projectile component requires physics component")
		}

		s.entities = append(s.entities, entityAspect{
			id:             eid,
			projectileCmpt: projectileCmpt,
			positionCmpt:   positionCmpt,
			moveCmpt:       moveCmpt,
			physicsCmpt:    physicsCmpt,
		})

		s.entityMap[eid] = &s.entities[len(s.entities)-1]
	}
}

func (s *System) ComponentsWillLeave(eid entity.ID, components []entity.IComponent) {
	remove := false
	for _, cmpt := range components {
		switch cmpt.(type) {
		case *Component, *movesys.Component, *positionsys.Component, *physicssys.Component:
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
