package physicssys

import "github.com/go-gl/mathgl/mgl32"

type (
	Component struct {
		velocity    mgl32.Vec2
		maxVelocity float32

		totalCurrentForce     mgl32.Vec2
		instantaneousVelocity mgl32.Vec2

		boundingBox BoundingBox
		collisions  []Collision

		collisionMask, collidesWith uint64
	}
)

func NewComponent(velocity mgl32.Vec2, maxVelocity float32, totalCurrentForce mgl32.Vec2, boundingBox BoundingBox, collisionMask uint64, collidesWith uint64) *Component {
	return &Component{
		velocity:          velocity,
		maxVelocity:       maxVelocity,
		totalCurrentForce: totalCurrentForce,
		boundingBox:       boundingBox,
		collisionMask:     collisionMask,
		collidesWith:      collidesWith,
		collisions:        []Collision{},
	}
}

func (c *Component) Velocity() mgl32.Vec2 {
	return c.velocity
}

func (c *Component) SetVelocity(velocity mgl32.Vec2) {
	c.velocity = velocity
}

func (c *Component) MaxVelocity() float32 {
	return c.maxVelocity
}

func (c *Component) SetMaxVelocity(maxVelocity float32) {
	c.maxVelocity = maxVelocity
}

func (c *Component) AddForce(f mgl32.Vec2) {
	c.totalCurrentForce = c.totalCurrentForce.Add(f)
}

func (c *Component) SetInstantaneousVelocity(v mgl32.Vec2) {
	c.instantaneousVelocity = v
}

func (c *Component) CurrentForces() mgl32.Vec2 {
	return c.totalCurrentForce
}

func (c *Component) InstantaneousVelocity() mgl32.Vec2 {
	return c.instantaneousVelocity
}

func (c *Component) ResetForces() {
	c.totalCurrentForce = mgl32.Vec2{}
}

func (c *Component) Collisions() []Collision {
	return c.collisions
}
