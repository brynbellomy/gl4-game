package positionsys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
)

type (
	Component struct {
		pos    mgl32.Vec2
		size   common.Size
		zindex int
	}

	IComponent interface {
		Pos() mgl32.Vec2
		SetPos(pos mgl32.Vec2)

		Size() common.Size
		SetSize(size common.Size)

		ZIndex() int
		SetZIndex(z int)
	}
)

func NewComponent(pos mgl32.Vec2, size common.Size, z int) IComponent {
	return &Component{pos: pos, size: size, zindex: z}
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

func (c *Component) ZIndex() int {
	return c.zindex
}

func (c *Component) SetZIndex(z int) {
	c.zindex = z
}
