package tilemapsys

import (
	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/rendersys"
	"github.com/brynbellomy/gl4-game/systems/rendersys/shader"
	"github.com/brynbellomy/gl4-game/systems/rendersys/texture"
)

type (
	System struct {
		tilemapCache       *TilemapCache
		shaderProgramCache *shader.ProgramCache
		textureCache       *texture.TextureCache
	}
)

// ensure that System conforms to entity.ISystem
var _ entity.ISystem = &System{}

func New(tilemapCache *TilemapCache, shaderProgramCache *shader.ProgramCache, textureCache *texture.TextureCache) *System {
	return &System{tilemapCache, shaderProgramCache, textureCache}
}

func (s *System) Name() string {
	return "tilemap"
}

func (s *System) Update(t common.Time) {
	// no-op
}

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"tilemap": {
			Slice: ComponentSlice{},
			Coder: common.NewCoder(common.CoderConfig{
				ConfigType: ComponentCfg{},
				Tag:        "config",
				Decode: func(x interface{}) (interface{}, error) {
					cfg := x.(ComponentCfg)

					tilemap, err := s.tilemapCache.Load(cfg.Tilemap)
					if err != nil {
						return nil, err
					}

					return Component{Tilemap: tilemap}, nil
				},
				Encode: func(x interface{}) (interface{}, error) { /* @@TODO */ panic("unimplemented") },
			}),
		},
	}
}

func (s *System) RenderNodeFactories() map[string]rendersys.INodeFactory {
	return map[string]rendersys.INodeFactory{
		"tilemap": &TilemapNodeFactory{
			shaderProgramCache: s.shaderProgramCache,
			textureCache:       s.textureCache,
		},
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	// no-op
}

func (s *System) ComponentsWillJoin(eid entity.ID, cmpts []entity.IComponent) error {
	var renderCmptIdx = -1
	var tilemapCmptIdx = -1
	for i := range cmpts {
		if _, is := cmpts[i].(rendersys.Component); is {
			renderCmptIdx = i
		} else if _, is := cmpts[i].(Component); is {
			tilemapCmptIdx = i
		}

		if renderCmptIdx != -1 && tilemapCmptIdx != -1 {
			break
		}
	}

	if renderCmptIdx != -1 && tilemapCmptIdx != -1 {
		renderCmpt := cmpts[renderCmptIdx].(rendersys.Component)
		tilemapCmpt := cmpts[tilemapCmptIdx].(Component)

		tilemap := tilemapCmpt.Tilemap

		renderCmpt.RenderNode().(*TilemapNode).SetTilemap(tilemap)

		cmpts[renderCmptIdx] = renderCmpt
		cmpts[tilemapCmptIdx] = tilemapCmpt
	}

	return nil
}

func (s *System) ComponentsWillLeave(eid entity.ID, cmpts []entity.IComponent) error {
	// no-op
	return nil
}
