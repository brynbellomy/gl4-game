package mainscene

// import (
// 	"path"

// 	"github.com/go-gl/mathgl/mgl32"

// 	"github.com/brynbellomy/gl4-game/common"
// 	"github.com/brynbellomy/gl4-game/entity"
// 	"github.com/brynbellomy/gl4-game/systems/animationsys"
// 	"github.com/brynbellomy/gl4-game/systems/gameobjsys"
// 	"github.com/brynbellomy/gl4-game/systems/movesys"
// 	"github.com/brynbellomy/gl4-game/systems/physicssys"
// 	"github.com/brynbellomy/gl4-game/systems/positionsys"
// 	"github.com/brynbellomy/gl4-game/systems/rendersys"
// 	"github.com/brynbellomy/gl4-game/systems/rendersys/texture"
// )

// func hero(assetRoot string) ([]entity.IComponent, error) {
// 	heroTexture, err := texture.Load(path.Join(assetRoot, "textures/lavos/walking-down-001.png"))
// 	if err != nil {
// 		return nil, err
// 	}

// 	animationMap := map[gameobjsys.Action]map[gameobjsys.Direction]string{
// 		gameobjsys.Action(0): map[gameobjsys.Direction]string{
// 			gameobjsys.Down:  "walking-down",
// 			gameobjsys.Left:  "walking-left",
// 			gameobjsys.Up:    "walking-up",
// 			gameobjsys.Right: "walking-right",
// 		},
// 	}

// 	boundingBox := physicssys.BoundingBox{
// 		{-0.1, -0.2},
// 		{-0.1, 0.2},
// 		{0.1, 0.2},
// 		{0.1, -0.2},
// 	}

// 	spriteNode, err := rendersys.NewSpriteNode(
// 		path.Join(assetRoot, "shaders/default-sprite.vertex.glsl"),
// 		path.Join(assetRoot, "shaders/default-sprite.fragment.glsl"),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return []entity.IComponent{
// 		positionsys.NewComponent(mgl32.Vec2{0, 0}, common.Size{0.2, 0.4}, 1, 0),
// 		physicssys.NewComponent(mgl32.Vec2{}, 20, mgl32.Vec2{}, boundingBox, uint64(HeroCollider), 0),
// 		rendersys.NewComponent(spriteNode, heroTexture),
// 		animationsys.NewComponent("hero", "walking-down", false, 0, 2),
// 		gameobjsys.NewComponent(gameobjsys.Action(0), gameobjsys.Down, animationMap),
// 		movesys.NewComponent(mgl32.Vec2{0, 0}),
// 	}, nil
// }
