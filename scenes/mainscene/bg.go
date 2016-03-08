package mainscene

// import (
// 	"path"

// 	"github.com/go-gl/mathgl/mgl32"

// 	"github.com/brynbellomy/gl4-game/common"
// 	"github.com/brynbellomy/gl4-game/entity"
// 	"github.com/brynbellomy/gl4-game/systems/positionsys"
// 	"github.com/brynbellomy/gl4-game/systems/rendersys"
// 	"github.com/brynbellomy/gl4-game/systems/spritesys"
// )

// func bg(assetRoot string) ([]entity.IComponent, error) {
// 	// bgTexture, err := texture.Load(path.Join(assetRoot, "textures/square.png"))
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	spriteNode, err := rendersys.NewSpriteNode(
// 		path.Join(assetRoot, "shaders/default-sprite.vertex.glsl"),
// 		path.Join(assetRoot, "shaders/default-sprite.fragment.glsl"),
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return []entity.IComponent{
// 		positionsys.NewComponent(mgl32.Vec2{0, 0}, common.Size{2.0, 2.0}, 0, 0),
// 		rendersys.NewComponent(spriteNode, 0, "shaders/default-sprite.vertex.glsl", "shaders/default-sprite.fragment.glsl"),
// 		spritesys.NewComponent("textures/square.png"),
// 	}, nil
// }
