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

func (s *System) Update(t common.Time) {
	// no-op
}

func (s *System) WillJoinManager(em *entity.Manager) {
	// no-op
}

func (s *System) ComponentsWillJoin(eid entity.ID, components []entity.IComponent) {
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
