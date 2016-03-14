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

func New(atlasCache *texture.AtlasCache) *System {
	return &System{
		atlasCache: atlasCache,
	}
}

// func (s *System) GetAnimation(eid entity.ID) string {
// 	if e, exists := s.entityMap[eid]; exists {
// 		return e.animationCmpt.Animation
// 	} else {
// 		panic("entity does not exist")
// 	}
// }

// func (s *System) SetAnimation(eid entity.ID, animation string, animationStart common.Time) {
// 	if e, exists := s.entityMap[eid]; exists {
// 		e.animationCmpt.Animation = animation
// 		e.animationCmpt.IsAnimating = true

// 	} else {
// 		panic("entity does not exist")
// 	}
// }

// func (s *System) StopAnimating(eid entity.ID) {
// 	if e, exists := s.entityMap[eid]; exists {
// 		e.animationCmpt.IsAnimating = false

// 	} else {
// 		panic("entity does not exist")
// 	}
// }

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

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"animation": {Component{}, ComponentSlice{}},
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

// func (s *System) EntityComponentsChanged(eid entity.ID, components []entity.IComponent) {
// 	var animationCmpt *Component
// 	var renderCmpt *rendersys.Component

// 	for _, cmpt := range components {
// 		if ac, is := cmpt.(*Component); is {
// 			animationCmpt = ac
// 		} else if rc, is := cmpt.(*rendersys.Component); is {
// 			renderCmpt = rc
// 		}

// 		if animationCmpt != nil && renderCmpt != nil {
// 			break
// 		}
// 	}

// 	if animationCmpt != nil && renderCmpt != nil {
// 		if _, exists := s.entityMap[eid]; !exists {
// 			s.entities = append(s.entities, entityAspect{
// 				id:            eid,
// 				animationCmpt: animationCmpt,
// 				renderCmpt:    renderCmpt,
// 			})

// 			s.entityMap[eid] = &s.entities[len(s.entities)-1]
// 		}

// 	} else {
// 		if _, exists := s.entityMap[eid]; exists {
// 			idx := -1
// 			for i := range s.entities {
// 				if s.entities[i].id == eid {
// 					idx = i
// 					break
// 				}
// 			}

// 			if idx >= 0 {
// 				s.entities = append(s.entities[:idx], s.entities[idx+1:]...)
// 			}

// 			delete(s.entityMap, eid)
// 		}
// 	}
// }

// func (s *System) ComponentsWillLeave(eid entity.ID, components []entity.IComponent) {
// 	remove := false
// 	for _, cmpt := range components {
// 		switch cmpt.(type) {
// 		case *Component, *rendersys.Component:
// 			remove = true
// 			break
// 		}
// 	}

// 	if remove {
// 		removedIdx := -1
// 		for i := range s.entities {
// 			if s.entities[i].id == eid {
// 				removedIdx = i
// 				break
// 			}
// 		}
// 		if removedIdx >= 0 {
// 			s.entities = append(s.entities[:removedIdx], s.entities[removedIdx+1:]...)
// 		}
// 		delete(s.entityMap, eid)
// 	}
// }
