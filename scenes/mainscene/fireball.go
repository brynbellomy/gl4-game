package mainscene

// import (
// 	"path"

// 	"github.com/go-gl/mathgl/mgl32"

// 	"github.com/brynbellomy/gl4-game/common"
// 	"github.com/brynbellomy/gl4-game/entity"
// 	"github.com/brynbellomy/gl4-game/systems/animationsys"
// 	"github.com/brynbellomy/gl4-game/systems/movesys"
// 	"github.com/brynbellomy/gl4-game/systems/physicssys"
// 	"github.com/brynbellomy/gl4-game/systems/positionsys"
// 	"github.com/brynbellomy/gl4-game/systems/projectilesys"
// 	"github.com/brynbellomy/gl4-game/systems/rendersys"
// 	"github.com/brynbellomy/gl4-game/systems/rendersys/texture"
// )

// type FireballFactory struct {
// 	AssetRoot string

// 	fireballTexture uint32
// 	fireballAtlas   *texture.Atlas
// }

// func NewFireballFactory(assetRoot string) (*FireballFactory, error) {
// 	f := &FireballFactory{
// 		AssetRoot: assetRoot,
// 	}

// 	err := f.init()
// 	return f, err
// }

// func (f *FireballFactory) init() error {
// 	fireballTexture, err := texture.Load(path.Join(f.AssetRoot, "textures/fireball/flying-001.png"))
// 	if err != nil {
// 		return err
// 	}

// 	f.fireballTexture = fireballTexture

// 	return err
// }

// func (f *FireballFactory) Build(pos mgl32.Vec2, vec mgl32.Vec2) ([]entity.IComponent, error) {
// 	boundingBox := physicssys.BoundingBox{
// 		{-0.1, -0.07},
// 		{-0.1, 0.07},
// 		{0.1, 0.07},
// 		{0.1, -0.07},
// 	}

// 	spriteNode, err := rendersys.NewSpriteNode(
// 		path.Join(f.AssetRoot, "shaders/fireball.vertex.glsl"),
// 		path.Join(f.AssetRoot, "shaders/fireball.fragment.glsl"),
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return []entity.IComponent{
// 		positionsys.NewComponent(pos, common.Size{0.2, 0.14}, 2, 0),
// 		physicssys.NewComponent(mgl32.Vec2{}, 8, mgl32.Vec2{}, boundingBox, uint64(ProjectileCollider), uint64(EnemyCollider)),
// 		rendersys.NewComponent(spriteNode, f.fireballTexture),
// 		animationsys.NewComponent("fireball", "flying", true, 0, 12),
// 		movesys.NewComponent(mgl32.Vec2{0, 0}),
// 		projectilesys.NewComponent(vec, 0.01, 10, projectilesys.Firing, true),
// 	}, nil
// }
