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

		velocity mgl32.Vec2
		// acceleration mgl32.Vec2

		totalCurrentForce mgl32.Vec2
	}

	IComponent interface {
		Pos() mgl32.Vec2
		SetPos(pos mgl32.Vec2)

		Size() common.Size
		SetSize(size common.Size)

		ZIndex() int
		SetZIndex(z int)

		Velocity() mgl32.Vec2
		SetVelocity(vel mgl32.Vec2)

		// Acceleration() mgl32.Vec2
		// SetAcceleration(accel mgl32.Vec2)

		AddForce(f mgl32.Vec2)
		CurrentForces() mgl32.Vec2
		ResetForces()
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

func (c *Component) Velocity() mgl32.Vec2 {
	return c.velocity
}

func (c *Component) SetVelocity(velocity mgl32.Vec2) {
	c.velocity = velocity
}

// func (c *Component) Acceleration() mgl32.Vec2 {
// 	return c.acceleration
// }

// func (c *Component) SetAcceleration(acceleration mgl32.Vec2) {
// 	c.acceleration = acceleration
// }

func (c *Component) AddForce(f mgl32.Vec2) {
	c.totalCurrentForce = c.totalCurrentForce.Add(f)
}

func (c *Component) CurrentForces() mgl32.Vec2 {
	return c.totalCurrentForce
}

func (c *Component) ResetForces() {
	c.totalCurrentForce = mgl32.Vec2{}
}
