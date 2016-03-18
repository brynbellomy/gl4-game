package animationsys

import (
	"math"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/rendersys"
	"github.com/brynbellomy/gl4-game/systems/rendersys/texture"
)

type (
	System struct {
		atlasCache *texture.AtlasCache

		entityManager    *entity.Manager
		componentQuery   entity.ComponentMask
		renderCmptSet    entity.IComponentSet
		animationCmptSet entity.IComponentSet
	}
)

// ensure that System conforms to entity.ISystem
var _ entity.ISystem = &System{}

func New(atlasCache *texture.AtlasCache) *System {
	return &System{
		atlasCache: atlasCache,
	}
}

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"animation": {
			Coder: common.NewCoder(common.CoderConfig{
				ConfigType: Component{},
				Tag:        "config",
				Decode:     func(x interface{}) (interface{}, error) { return x.(Component), nil },
				Encode:     func(x interface{}) (interface{}, error) { /* @@TODO */ panic("unimplemented") },
			}),
			Slice: ComponentSlice{},
		},
	}
}

func (s *System) Update(t common.Time) {
	matchIDs := s.entityManager.EntitiesMatching(s.componentQuery)
	renderCmptIdxs, err := s.renderCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}
	animationCmptIdxs, err := s.animationCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}

	animationCmptSlice := s.animationCmptSet.Slice().(ComponentSlice)
	renderCmptSlice := s.renderCmptSet.Slice().(rendersys.ComponentSlice)

	for i := 0; i < len(animationCmptIdxs); i++ {
		animationCmpt := animationCmptSlice[animationCmptIdxs[i]]
		renderCmpt := renderCmptSlice[renderCmptIdxs[i]]

		atlas, err := s.atlasCache.Load(animationCmpt.AtlasName)
		if err != nil {
			panic(err.Error())
		}

		textures := atlas.Animation(animationCmpt.Animation)
		if len(textures) <= 0 {
			panic("textures slice is empty")
		}

		if !animationCmpt.IsAnimating {
			animationCmpt.CurrentIndex = 0
			renderCmpt.SetTexture(textures[animationCmpt.CurrentIndex])

		} else {
			elapsedNano := t - animationCmpt.AnimationStart
			totalFrames := int64(math.Floor(elapsedNano.Seconds() * float64(animationCmpt.FPS)))
			newIndex := int(totalFrames % int64(len(textures)))

			if animationCmpt.CurrentIndex == 0 || newIndex != animationCmpt.CurrentIndex {
				animationCmpt.CurrentIndex = newIndex
				tex := textures[animationCmpt.CurrentIndex]
				renderCmpt.SetTexture(tex)
			}
		}

		animationCmptSlice[animationCmptIdxs[i]] = animationCmpt
		renderCmptSlice[renderCmptIdxs[i]] = renderCmpt
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	s.entityManager = em

	componentQuery, err := s.entityManager.MakeCmptQuery([]string{"render", "animation"})
	if err != nil {
		panic(err)
	}
	s.componentQuery = componentQuery

	renderCmptSet, err := s.entityManager.GetComponentSet("render")
	if err != nil {
		panic(err)
	}
	s.renderCmptSet = renderCmptSet

	animationCmptSet, err := s.entityManager.GetComponentSet("animation")
	if err != nil {
		panic(err)
	}
	s.animationCmptSet = animationCmptSet
}

func (s *System) ComponentsWillJoin(eid entity.ID, cmpts []entity.IComponent) error {
	// no-op
	return nil
}

func (s *System) ComponentsWillLeave(eid entity.ID, cmpts []entity.IComponent) error {
	// no-op
	return nil
}
