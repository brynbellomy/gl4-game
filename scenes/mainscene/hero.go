package mainscene

import (
	"path"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/node"
	"github.com/brynbellomy/gl4-game/systems/animationsys"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
	"github.com/brynbellomy/gl4-game/systems/rendersys"
	"github.com/brynbellomy/gl4-game/texture"
)

func hero(assetRoot string) ([]entity.IComponent, error) {
	heroTexture, err := texture.Load(path.Join(assetRoot, "textures/lavos/walking-down-001.png"))
	if err != nil {
		return nil, err
	}

	heroAtlas := texture.NewAtlas()
	err = heroAtlas.LoadAnimation("walking-down", []string{
		path.Join(assetRoot, "textures/lavos/walking-down-001.png"),
		path.Join(assetRoot, "textures/lavos/walking-down-002.png"),
		path.Join(assetRoot, "textures/lavos/walking-down-003.png"),
		path.Join(assetRoot, "textures/lavos/walking-down-004.png"),
	})
	if err != nil {
		return nil, err
	}

	err = heroAtlas.LoadAnimation("walking-left", []string{
		path.Join(assetRoot, "textures/lavos/walking-left-001.png"),
		path.Join(assetRoot, "textures/lavos/walking-left-002.png"),
		path.Join(assetRoot, "textures/lavos/walking-left-003.png"),
		path.Join(assetRoot, "textures/lavos/walking-left-004.png"),
	})
	if err != nil {
		return nil, err
	}

	err = heroAtlas.LoadAnimation("walking-up", []string{
		path.Join(assetRoot, "textures/lavos/walking-up-001.png"),
		path.Join(assetRoot, "textures/lavos/walking-up-002.png"),
		path.Join(assetRoot, "textures/lavos/walking-up-003.png"),
		path.Join(assetRoot, "textures/lavos/walking-up-004.png"),
	})
	if err != nil {
		return nil, err
	}

	err = heroAtlas.LoadAnimation("walking-right", []string{
		path.Join(assetRoot, "textures/lavos/walking-right-001.png"),
		path.Join(assetRoot, "textures/lavos/walking-right-002.png"),
		path.Join(assetRoot, "textures/lavos/walking-right-003.png"),
		path.Join(assetRoot, "textures/lavos/walking-right-004.png"),
	})
	if err != nil {
		return nil, err
	}

	return []entity.IComponent{
		positionsys.NewComponent(mgl32.Vec2{0, 0}, common.Size{0.2, 0.4}, 1),
		rendersys.NewComponent(node.NewSpriteNode(), heroTexture),
		animationsys.NewComponent(heroAtlas, "walking", 0, 2),
		// steeringsys.NewComponent([]steeringsys.IBehavior{}),
	}, nil
}
