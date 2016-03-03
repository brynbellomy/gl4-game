package projectilesys

import (
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
		return e.projectileCmpt.Heading()
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
	for _, ent := range s.entities {
		switch ent.projectileCmpt.state {
		case Firing:
			headingNorm := ent.projectileCmpt.Heading().Normalize()
			vel := headingNorm.Mul(ent.projectileCmpt.exitVelocity)
			ent.physicsCmpt.SetVelocity(vel)

			// ent.moveCmpt.SetVector(headingNorm.Mul(ent.projectileCmpt.acceleration))
			ent.physicsCmpt.AddForce(headingNorm.Mul(ent.projectileCmpt.acceleration))

			// only stay in the Firing state for the first frame
			ent.projectileCmpt.state = Flying

		case Flying:
			force := ent.projectileCmpt.Heading().Normalize().Mul(ent.projectileCmpt.acceleration)
			// ent.moveCmpt.SetVector(force)
			ent.physicsCmpt.AddForce(force)

		case Impacting:
			ent.physicsCmpt.SetVelocity(mgl32.Vec2{0, 0})
		}
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	s.entityManager = em
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
			projectileCmpt: projectileCmpt,
			positionCmpt:   positionCmpt,
			moveCmpt:       moveCmpt,
			physicsCmpt:    physicsCmpt,
		})

		s.entityMap[eid] = &s.entities[len(s.entities)-1]
	}
}
