package positionsys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	Component struct {
		Pos      mgl32.Vec2  `config:"pos"`
		Size     common.Size `config:"size"`
		Rotation float32     `config:"rotation"`
		ZIndex   int         `config:"z-index"`
	}
)

func NewComponent(pos mgl32.Vec2, size common.Size, z int, rotation float32) *Component {
	return &Component{Pos: pos, Size: size, ZIndex: z, Rotation: rotation}
}

func (c *Component) GetPos() mgl32.Vec2 {
	return c.Pos
}

func (c *Component) SetPos(pos mgl32.Vec2) {
	c.Pos = pos
}

func (c *Component) GetRotation() float32 {
	return c.Rotation
}

func (c *Component) SetRotation(rotation float32) {
	c.Rotation = rotation
}

func (c *Component) GetSize() common.Size {
	return c.Size
}

func (c *Component) SetSize(size common.Size) {
	c.Size = size
}

func (c *Component) GetZIndex() int {
	return c.ZIndex
}

func (c *Component) SetZIndex(z int) {
	c.ZIndex = z
}

func (c *Component) Clone() entity.IComponent {
	x := *c
	return &x
}
