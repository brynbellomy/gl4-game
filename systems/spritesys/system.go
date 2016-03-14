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

func New(textureCache *texture.TextureCache) *System {
	return &System{
		textureCache: textureCache,
	}
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

		if !spriteCmpt.IsTextureLoaded {
			textureName := spriteCmpt.GetTextureName()

			var tex uint32
			if textureName != "" {
				t, err := s.textureCache.Load(textureName)
				if err != nil {
					panic(err.Error())
				}
				tex = t
			}
			spriteCmpt.SetTexture(tex)
			spriteCmpt.IsTextureLoaded = true
		}

		renderCmpt.SetTexture(spriteCmpt.GetTexture())

		spriteCmptSlice[spriteCmptIdxs[i]] = spriteCmpt
		renderCmptSlice[renderCmptIdxs[i]] = renderCmpt
	}
}

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"sprite": {Component{}, ComponentSlice{}},
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

// func (s *System) EntityComponentsChanged(eid entity.ID, components []entity.IComponent) {
// 	var spriteCmpt *Component
// 	var renderCmpt *rendersys.Component

// 	for _, cmpt := range components {
// 		if rc, is := cmpt.(*Component); is {
// 			spriteCmpt = rc
// 		} else if pc, is := cmpt.(*rendersys.Component); is {
// 			renderCmpt = pc
// 		}

// 		if renderCmpt != nil && spriteCmpt != nil {
// 			break
// 		}
// 	}

// 	if renderCmpt != nil && spriteCmpt != nil {
// 		if _, exists := s.entityMap[eid]; !exists {
// 			s.entities = append(s.entities, entityAspect{
// 				id:         eid,
// 				renderCmpt: renderCmpt,
// 				spriteCmpt: spriteCmpt,
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
