package positionsys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/context"
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	System struct {
		entities  []entityAspect
		entityMap map[entity.ID]*entityAspect
	}

	entityAspect struct {
		positionCmpt IComponent
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

func (s *System) Update(c context.IContext) {
	// for _, e := range s.entities {
	// 	curpos := e.positionCmpt.Pos()
	// 	newpos := mgl32.Vec2{
	// 		curpos.X() + 0.001,
	// 		curpos.Y() + 0.001,
	// 	}
	// 	e.positionCmpt.SetPos(newpos)
	// }
}

func (s *System) EntityWillJoin(eid entity.ID, components []entity.IComponent) {
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

type (
	Component struct {
		pos  mgl32.Vec2
		size common.Size
		// zorder int
	}

	IComponent interface {
		Pos() mgl32.Vec2
		SetPos(pos mgl32.Vec2)

		Size() common.Size
		SetSize(size common.Size)

		// ZOrder() int
		// SetZOrder(z int)
	}
)

func NewComponent(pos mgl32.Vec2, size common.Size) IComponent {
	return &Component{pos: pos, size: size}
}

func (c *Component) Pos() mgl32.Vec2 {
	return c.pos
}

func (c *Component) SetPos(pos mgl32.Vec2) {
	c.pos = pos
}

func (c *Component) Size() common.Size {
	return c.size
}

func (c *Component) SetSize(size common.Size) {
	c.size = size
}

// func (c *Component) ZOrder() int {
// 	return c.zorder
// }

// func (c *Component) SetZOrder(z int) {
// 	c.zorder = z
// }
