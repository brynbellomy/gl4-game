package mainscene

import (
	"path"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/animationsys"
	"github.com/brynbellomy/gl4-game/systems/gameobjsys"
	"github.com/brynbellomy/gl4-game/systems/movesys"
	"github.com/brynbellomy/gl4-game/systems/physicssys"
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

	animationMap := map[gameobjsys.Action]map[gameobjsys.Direction]string{
		gameobjsys.Action(0): map[gameobjsys.Direction]string{
			gameobjsys.Down:  "walking-down",
			gameobjsys.Left:  "walking-left",
			gameobjsys.Up:    "walking-up",
			gameobjsys.Right: "walking-right",
		},
	}

	return []entity.IComponent{
		positionsys.NewComponent(mgl32.Vec2{0, 0}, common.Size{0.2, 0.4}, 1),
		physicssys.NewComponent(mgl32.Vec2{}, 20, mgl32.Vec2{}),
		rendersys.NewComponent(rendersys.NewSpriteNode(), heroTexture),
		animationsys.NewComponent(heroAtlas, "walking-down", 0, 2),
		gameobjsys.NewComponent(gameobjsys.Action(0), gameobjsys.Down, animationMap),
		movesys.NewComponent(mgl32.Vec2{0, 0}),
	}, nil
}
