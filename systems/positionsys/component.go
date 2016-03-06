package positionsys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
)

type (
	Component struct {
		pos      mgl32.Vec2
		size     common.Size
		rotation float32
		zindex   int
	}
)

func NewComponent(pos mgl32.Vec2, size common.Size, z int, rotation float32) *Component {
	return &Component{pos: pos, size: size, zindex: z, rotation: rotation}
}

func (c *Component) Pos() mgl32.Vec2 {
	return c.pos
}

func (c *Component) SetPos(pos mgl32.Vec2) {
	c.pos = pos
}

func (c *Component) Rotation() float32 {
	return c.rotation
}

func (c *Component) SetRotation(rotation float32) {
	c.rotation = rotation
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
