package spritesys

import (
	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/rendersys"
	"github.com/brynbellomy/gl4-game/systems/rendersys/texture"
)

type (
	System struct {
		textureCache *texture.TextureCache

		entityManager  *entity.Manager
		componentQuery entity.ComponentMask

		renderCmptSet entity.IComponentSet
		spriteCmptSet entity.IComponentSet
	}
)

// ensure that System conforms to entity.ISystem
var _ entity.ISystem = &System{}

func New(textureCache *texture.TextureCache) *System {
	return &System{
		textureCache: textureCache,
	}
}

func (s *System) Name() string {
	return "sprite"
}

func (s *System) Update(t common.Time) {
	matchIDs := s.entityManager.EntitiesMatching(s.componentQuery)
	renderCmptIdxs, err := s.renderCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}
	spriteCmptIdxs, err := s.spriteCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}

	renderCmptSlice := s.renderCmptSet.Slice().(rendersys.ComponentSlice)
	spriteCmptSlice := s.spriteCmptSet.Slice().(ComponentSlice)

	for i := 0; i < len(spriteCmptIdxs); i++ {
		spriteCmpt := spriteCmptSlice[spriteCmptIdxs[i]]
		renderCmpt := renderCmptSlice[renderCmptIdxs[i]]

		renderCmpt.SetTexture(spriteCmpt.GetTexture())

		spriteCmptSlice[spriteCmptIdxs[i]] = spriteCmpt
		renderCmptSlice[renderCmptIdxs[i]] = renderCmpt
	}
}

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"sprite": {
			Slice: ComponentSlice{},
			Coder: common.NewCoder(common.CoderConfig{
				ConfigType: ComponentCfg{},
				Tag:        "config",
				Decode: func(x interface{}) (interface{}, error) {
					cfg := x.(ComponentCfg)
					t, err := s.textureCache.Load(cfg.TextureName)
					if err != nil {
						return nil, err
					}
					return Component{Texture: t}, nil
				},
				Encode: func(x interface{}) (interface{}, error) { /* @@TODO */ panic("unimplemented") },
			}),
		},
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	s.entityManager = em

	componentQuery, err := s.entityManager.MakeCmptQuery([]string{"sprite", "render"})
	if err != nil {
		panic(err)
	}
	s.componentQuery = componentQuery

	renderCmptSet, err := s.entityManager.GetComponentSet("render")
	if err != nil {
		panic(err)
	}
	s.renderCmptSet = renderCmptSet

	spriteCmptSet, err := s.entityManager.GetComponentSet("sprite")
	if err != nil {
		panic(err)
	}
	s.spriteCmptSet = spriteCmptSet
}

func (s *System) ComponentsWillJoin(eid entity.ID, cmpts []entity.IComponent) error {
	// var spriteCmptIdx = -1

	// for i := range cmpts {
	// 	if _, is := cmpts[i].(Component); is {
	// 		spriteCmptIdx = i
	// 		break
	// 	}
	// }

	// if spriteCmptIdx != -1 {
	// 	spriteCmpt := cmpts[spriteCmptIdx].(Component)
	// 	textureName := spriteCmpt.GetTextureName()

	// 	var tex texture.TextureID
	// 	if textureName != "" {
	// 		t, err := s.textureCache.Load(textureName)
	// 		if err != nil {
	// 			panic(err.Error())
	// 		}
	// 		tex = t
	// 	}
	// 	spriteCmpt.SetTexture(tex)
	// 	cmpts[spriteCmptIdx] = spriteCmpt
	// }

	return nil
}

func (s *System) ComponentsWillLeave(eid entity.ID, cmpts []entity.IComponent) error {
	// no-op
	return nil
}
