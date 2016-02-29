package steeringbehaviors

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Constant struct {
	Vec mgl32.Vec2
}

func (c *Constant) SteeringVector() mgl32.Vec2 {
	return c.Vec
}
