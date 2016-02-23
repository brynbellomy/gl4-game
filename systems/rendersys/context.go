package rendersys

import (
	"github.com/go-gl/mathgl/mgl32"
)

type (
	RenderContext struct {
		currentTransform mgl32.Mat4
	}
)

func NewRenderContext() IRenderContext {
	return &RenderContext{}
}

func (c *RenderContext) CurrentTransform() mgl32.Mat4 {
	return c.currentTransform
}
