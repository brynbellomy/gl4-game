package spritesys

import (
	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/rendersys"
	"github.com/brynbellomy/gl4-game/systems/rendersys/texture"
)

type (
	System struct {
		entities     []entityAspect
		entityMap    map[entity.ID]*entityAspect
		textureCache *texture.TextureCache
	}

	entityAspect struct {
		id         entity.ID
		spriteCmpt *Component
		renderCmpt *rendersys.Component
	}
)

func New(textureCache *texture.TextureCache) *System {
	return &System{
		entities:     []entityAspect{},
		entityMap:    map[entity.ID]*entityAspect{},
		textureCache: textureCache,
	}
}

func (s *System) Update(t common.Time) {
	for _, ent := range s.entities {
		if !ent.spriteCmpt.IsTextureLoaded {
			textureName := ent.spriteCmpt.GetTextureName()
			var tex uint32
			if textureName != "" {
				t, err := s.textureCache.Load(textureName)
				if err != nil {
					panic(err.Error())
				}
				tex = t
			}
			ent.spriteCmpt.SetTexture(tex)
			ent.spriteCmpt.IsTextureLoaded = true
		}
	}

	for _, ent := range s.entities {
		ent.renderCmpt.SetTexture(ent.spriteCmpt.GetTexture())
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	em.RegisterComponentType("sprite", &Component{}, nil)
}

func (s *System) ComponentsWillJoin(eid entity.ID, components []entity.IComponent) {
	var spriteCmpt *Component
	var renderCmpt *rendersys.Component

	for _, cmpt := range components {
		if rc, is := cmpt.(*Component); is {
			spriteCmpt = rc
		} else if pc, is := cmpt.(*rendersys.Component); is {
			renderCmpt = pc
		}

		if renderCmpt != nil && spriteCmpt != nil {
			break
		}
	}

	if spriteCmpt != nil {
		if renderCmpt == nil {
			panic("sprite component requires render component")
		}

		s.entities = append(s.entities, entityAspect{
			id:         eid,
			renderCmpt: renderCmpt,
			spriteCmpt: spriteCmpt,
		})

		s.entityMap[eid] = &s.entities[len(s.entities)-1]
	}
}

func (s *System) ComponentsWillLeave(eid entity.ID, components []entity.IComponent) {
	remove := false
	for _, cmpt := range components {
		switch cmpt.(type) {
		case *Component, *rendersys.Component:
			remove = true
			break
		}
	}

	if remove {
		removedIdx := -1
		for i := range s.entities {
			if s.entities[i].id == eid {
				removedIdx = i
				break
			}
		}
		if removedIdx >= 0 {
			s.entities = append(s.entities[:removedIdx], s.entities[removedIdx+1:]...)
		}
		delete(s.entityMap, eid)
	}
}
