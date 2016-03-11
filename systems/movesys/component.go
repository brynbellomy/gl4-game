package movesys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/entity"
)

type (
	Component struct {
		Vec mgl32.Vec2 `config:"vector"`

		entity.ComponentKind `config:"-"`
	}
)

func NewComponent(vec mgl32.Vec2) *Component {
	return &Component{Vec: vec}
}

func (c *Component) Vector() mgl32.Vec2 {
	return c.Vec
}

func (c *Component) SetVector(vec mgl32.Vec2) {
	c.Vec = vec
}

func (c *Component) ResetVector() {
	c.Vec = mgl32.Vec2{0, 0}
}

func (c *Component) Clone() entity.IComponent {
	x := *c
	return &x
}
