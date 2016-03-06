package mainscene

import (
	"path"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
	"github.com/brynbellomy/gl4-game/systems/rendersys"
	"github.com/brynbellomy/gl4-game/systems/rendersys/texture"
)

func bg(assetRoot string) ([]entity.IComponent, error) {
	bgTexture, err := texture.Load(path.Join(assetRoot, "textures/square.png"))
	if err != nil {
		return nil, err
	}

	return []entity.IComponent{
		positionsys.NewComponent(mgl32.Vec2{0, 0}, common.Size{2.0, 2.0}, 0, 0),
		rendersys.NewComponent(rendersys.NewDefaultSpriteNode(), bgTexture),
	}, nil
}
