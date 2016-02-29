package rendersys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
)

type (
	INode interface {
		Render(c RenderContext)

		SetPos(p mgl32.Vec2)
		SetSize(s common.Size)
		SetTexture(tex uint32)
	}
)
