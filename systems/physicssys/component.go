package physicssys

import "github.com/go-gl/mathgl/mgl32"

type (
	Component struct {
		velocity mgl32.Vec2
		// acceleration mgl32.Vec2

		totalCurrentForce mgl32.Vec2
		// totalInstantaneousVelocity mgl32.Vec2
	}
)

func NewComponent() *Component {
	return &Component{}
}

func (c *Component) Velocity() mgl32.Vec2 {
	return c.velocity
}

func (c *Component) SetVelocity(velocity mgl32.Vec2) {
	c.velocity = velocity
}

func (c *Component) AddForce(f mgl32.Vec2) {
	c.totalCurrentForce = c.totalCurrentForce.Add(f)
}

// func (c *Component) AddInstantaneousVelocity(v mgl32.Vec2) {
// 	c.totalInstantaneousVelocity = c.totalInstantaneousVelocity.Add(v)
// }

func (c *Component) CurrentForces() mgl32.Vec2 {
	return c.totalCurrentForce
}

// func (c *Component) CurrentInstantaneousVelocity() mgl32.Vec2 {
// 	return c.totalInstantaneousVelocity
// }

func (c *Component) ResetForces() {
	c.totalCurrentForce = mgl32.Vec2{}
}

// func (c *Component) ResetInstantaneousVelocity() {
// 	c.totalInstantaneousVelocity = mgl32.Vec2{}
// }
