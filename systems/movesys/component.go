package movesys

import (
	"github.com/go-gl/mathgl/mgl32"
)

type (
	Component struct {
		// kind MovementKind
		vec mgl32.Vec2
	}

	// MovementKind int
)

func NewComponent(vec mgl32.Vec2) *Component {
	return &Component{vec: vec}
}

func (c *Component) Vector() mgl32.Vec2 {
	return c.vec
}

func (c *Component) SetVector(vec mgl32.Vec2) {
	c.vec = vec
}

// const (
// 	ConstantMovement MovementKind = iota
// 	GoalMovement
// )
