package steeringbehaviors

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/entity"
)

type Follow struct {
	Vec    mgl32.Vec2
	Self   entity.ID
	Target entity.ID
}

func (c *Follow) SteeringVector() mgl32.Vec2 {
	return c.Vec
}
