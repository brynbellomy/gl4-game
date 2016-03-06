package movesys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/physicssys"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
)

type (
	System struct {
		entities  []entityAspect
		entityMap map[entity.ID]*entityAspect
	}

	entityAspect struct {
		id           entity.ID
		moveCmpt     *Component
		positionCmpt *positionsys.Component
		physicsCmpt  *physicssys.Component
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
		// e.physicsCmpt.AddForce(e.moveCmpt.Vector())
		vec := e.moveCmpt.Vector()
		// if vec.X() != 0 || vec.Y() != 0 {
		e.physicsCmpt.SetInstantaneousVelocity(vec)
		// }
		e.moveCmpt.ResetVector()
	}
}

func (s *System) SetMovementVector(eid entity.ID, vec mgl32.Vec2) {
	if e, exists := s.entityMap[eid]; exists {
		e.moveCmpt.SetVector(vec)
	} else {
		panic("entity does not exist")
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	// no-op
}

func (s *System) ComponentsWillJoin(eid entity.ID, components []entity.IComponent) {
	var moveCmpt *Component
	var physicsCmpt *physicssys.Component
	var positionCmpt *positionsys.Component

	for _, cmpt := range components {
		if ac, is := cmpt.(*Component); is {
			moveCmpt = ac
		} else if rc, is := cmpt.(*positionsys.Component); is {
			positionCmpt = rc
		} else if pc, is := cmpt.(*physicssys.Component); is {
			physicsCmpt = pc
		}

		if moveCmpt != nil && positionCmpt != nil && physicsCmpt != nil {
			break
		}
	}

	if moveCmpt != nil {
		if positionCmpt == nil {
			panic("move component requires position component")
		} else if physicsCmpt == nil {
			panic("move component requires physics component")
		}

		s.entities = append(s.entities, entityAspect{
			id:           eid,
			moveCmpt:     moveCmpt,
			positionCmpt: positionCmpt,
			physicsCmpt:  physicsCmpt,
		})

		s.entityMap[eid] = &s.entities[len(s.entities)-1]
	}
}

func (s *System) ComponentsWillLeave(eid entity.ID, components []entity.IComponent) {
	remove := false
	for _, cmpt := range components {
		switch cmpt.(type) {
		case *Component, *positionsys.Component, *physicssys.Component:
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
