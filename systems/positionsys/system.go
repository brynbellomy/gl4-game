package positionsys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/context"
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	System struct {
		entities []entityAspect
	}

	entityAspect struct {
		positionCmpt IComponent
	}
)

func New() *System {
	return &System{
		entities: make([]entityAspect, 0),
	}
}

func (s *System) Update(c context.IContext) {
	for _, e := range s.entities {
		curpos := e.positionCmpt.Pos()
		newpos := mgl32.Vec2{
			curpos.X() + 0.001,
			curpos.Y() + 0.001,
		}
		e.positionCmpt.SetPos(newpos)
	}
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
	}
}

type (
	Component struct {
		pos mgl32.Vec2
	}

	IComponent interface {
		Pos() mgl32.Vec2
		SetPos(pos mgl32.Vec2)
	}
)

func NewComponent(pos mgl32.Vec2) IComponent {
	return &Component{pos}
}

func (c *Component) Pos() mgl32.Vec2 {
	return c.pos
}

func (c *Component) SetPos(pos mgl32.Vec2) {
	c.pos = pos
}
