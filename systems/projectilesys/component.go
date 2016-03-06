package projectilesys

import "github.com/go-gl/mathgl/mgl32"

type (
	Component struct {
		heading         mgl32.Vec2
		exitVelocity    float32
		acceleration    float32
		state           ProjectileState
		removeOnContact bool
	}

	ProjectileState int
)

const (
	Firing ProjectileState = iota
	Flying
	Impacting
)

func NewComponent(heading mgl32.Vec2, exitVelocity float32, acceleration float32, state ProjectileState, removeOnContact bool) *Component {
	return &Component{heading: heading, exitVelocity: exitVelocity, acceleration: acceleration, state: state, removeOnContact: removeOnContact}
}

func (c *Component) Heading() mgl32.Vec2 {
	return c.heading
}

func (c *Component) SetHeading(heading mgl32.Vec2) {
	c.heading = heading
}
