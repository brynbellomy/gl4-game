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
	"github.com/brynbellomy/gl4-game/systems/rendersys/texture"
)

func skeleton(assetRoot string) ([]entity.IComponent, error) {
	skeletonTexture, err := texture.Load(path.Join(assetRoot, "textures/skeleton/walking-down-001.png"))
	if err != nil {
		return nil, err
	}

	skeletonAtlas := texture.NewAtlas()
	err = skeletonAtlas.LoadAnimation("walking-down", []string{
		path.Join(assetRoot, "textures/skeleton/walking-down-001.png"),
		path.Join(assetRoot, "textures/skeleton/walking-down-002.png"),
		path.Join(assetRoot, "textures/skeleton/walking-down-003.png"),
	})
	if err != nil {
		return nil, err
	}

	err = skeletonAtlas.LoadAnimation("walking-left", []string{
		path.Join(assetRoot, "textures/skeleton/walking-left-001.png"),
		path.Join(assetRoot, "textures/skeleton/walking-left-002.png"),
		path.Join(assetRoot, "textures/skeleton/walking-left-003.png"),
	})
	if err != nil {
		return nil, err
	}

	err = skeletonAtlas.LoadAnimation("walking-up", []string{
		path.Join(assetRoot, "textures/skeleton/walking-up-001.png"),
		path.Join(assetRoot, "textures/skeleton/walking-up-002.png"),
		path.Join(assetRoot, "textures/skeleton/walking-up-003.png"),
	})
	if err != nil {
		return nil, err
	}

	err = skeletonAtlas.LoadAnimation("walking-right", []string{
		path.Join(assetRoot, "textures/skeleton/walking-right-001.png"),
		path.Join(assetRoot, "textures/skeleton/walking-right-002.png"),
		path.Join(assetRoot, "textures/skeleton/walking-right-003.png"),
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

	boundingBox := physicssys.BoundingBox{
		{-0.1, -0.2},
		{-0.1, 0.2},
		{0.1, 0.2},
		{0.1, -0.2},
	}

	return []entity.IComponent{
		positionsys.NewComponent(mgl32.Vec2{1, 1}, common.Size{0.2, 0.4}, 1, 0),
		physicssys.NewComponent(mgl32.Vec2{}, 20, mgl32.Vec2{}, boundingBox, uint64(EnemyCollider), 0),
		rendersys.NewComponent(rendersys.NewDefaultSpriteNode(), skeletonTexture),
		animationsys.NewComponent(skeletonAtlas, "walking-down", false, 0, 2),
		gameobjsys.NewComponent(gameobjsys.Action(0), gameobjsys.Down, animationMap),
		movesys.NewComponent(mgl32.Vec2{0, 0}),
	}, nil
}
