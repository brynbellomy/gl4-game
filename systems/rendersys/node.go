package rendersys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/systems/rendersys/texture"
)

type (
	INode interface {
		Init() error
		Render(c RenderContext)
		Destroy() error

		SetPos(p mgl32.Vec2)
		SetSize(s common.Size)
		SetRotation(r float32)
		SetTexture(tex texture.TextureID)
	}
)
