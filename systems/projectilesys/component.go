package projectilesys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/entity"
)

type (
	Component struct {
		Heading         mgl32.Vec2      `config:"heading"`
		ExitVelocity    float32         `config:"exitVelocity"`
		Thrust          float32         `config:"thrust"`
		State           ProjectileState `config:"state"`
		RemoveOnContact bool            `config:"removeOnContact"`
	}

	ProjectileState int
)

const (
	Firing ProjectileState = iota
	Flying
	Impacting
)

func NewComponent(heading mgl32.Vec2, exitVelocity float32, thrust float32, state ProjectileState, removeOnContact bool) *Component {
	return &Component{Heading: heading, ExitVelocity: exitVelocity, Thrust: thrust, State: state, RemoveOnContact: removeOnContact}
}

func (c *Component) GetHeading() mgl32.Vec2 {
	return c.Heading
}

func (c *Component) SetHeading(heading mgl32.Vec2) {
	c.Heading = heading
}

func (c *Component) Clone() entity.IComponent {
	x := *c
	return &x
}
