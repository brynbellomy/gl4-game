package physicssys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/entity"
)

type (
	Component struct {
		Velocity              mgl32.Vec2 `config:"velocity"`
		MaxVelocity           float32    `config:"maxVelocity"`
		TotalCurrentForce     mgl32.Vec2 `config:"totalCurrentForce"`
		InstantaneousVelocity mgl32.Vec2 `config:"instantaneousVelocity"`

		BoundingBox   BoundingBox `config:"boundingBox"`
		Collisions    []Collision `config:"collisions"`
		CollisionMask uint64      `config:"collisionMask"`
		CollidesWith  uint64      `config:"collidesWith"`

		entity.ComponentKind `config:"-"`
	}

	ComponentSlice []Component
)

func NewComponent(velocity mgl32.Vec2, maxVelocity float32, totalCurrentForce mgl32.Vec2, boundingBox BoundingBox, collisionMask uint64, collidesWith uint64) *Component {
	return &Component{
		Velocity:              velocity,
		MaxVelocity:           maxVelocity,
		TotalCurrentForce:     totalCurrentForce,
		InstantaneousVelocity: mgl32.Vec2{0, 0},
		BoundingBox:           boundingBox,
		CollisionMask:         collisionMask,
		CollidesWith:          collidesWith,
		Collisions:            []Collision{},
	}
}

func (c *Component) GetVelocity() mgl32.Vec2 {
	return c.Velocity
}

func (c *Component) SetVelocity(velocity mgl32.Vec2) {
	c.Velocity = velocity
}

func (c *Component) GetMaxVelocity() float32 {
	return c.MaxVelocity
}

func (c *Component) SetMaxVelocity(maxVelocity float32) {
	c.MaxVelocity = maxVelocity
}

func (c *Component) AddForce(f mgl32.Vec2) {
	c.TotalCurrentForce = c.TotalCurrentForce.Add(f)
}

func (c *Component) SetInstantaneousVelocity(v mgl32.Vec2) {
	c.InstantaneousVelocity = v
}

func (c *Component) CurrentForces() mgl32.Vec2 {
	return c.TotalCurrentForce
}

func (c *Component) GetInstantaneousVelocity() mgl32.Vec2 {
	return c.InstantaneousVelocity
}

func (c *Component) ResetForces() {
	c.TotalCurrentForce = mgl32.Vec2{}
}

func (c *Component) GetCollisions() []Collision {
	return c.Collisions
}

func (c *Component) AddCollision(coll Collision) {
	c.Collisions = append(c.Collisions, coll)
}

func (c *Component) ResetCollisions() {
	c.Collisions = []Collision{}
}

func (c *Component) GetBoundingBox() BoundingBox {
	return c.BoundingBox
}

func (c Component) Clone() entity.IComponent {
	return c
}

func (cs ComponentSlice) Append(cmpt entity.IComponent) entity.IComponentSlice {
	return append(cs, cmpt.(Component))
}

func (cs ComponentSlice) Remove(idx int) entity.IComponentSlice {
	return append(cs[:idx], cs[idx+1:]...)
}

func (cs ComponentSlice) Get(idx int) (entity.IComponent, bool) {
	if idx >= len(cs) {
		return nil, false
	}
	return cs[idx], true
}

func (cs ComponentSlice) Set(idx int, cmpt entity.IComponent) bool {
	if idx >= len(cs) {
		return false
	}
	cs[idx] = cmpt.(Component)
	return true
}
