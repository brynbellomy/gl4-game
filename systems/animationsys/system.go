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
		entities   []entityAspect
		entityMap  map[entity.ID]*entityAspect
		atlasCache *texture.AtlasCache
	}

	entityAspect struct {
		id            entity.ID
		renderCmpt    *rendersys.Component
		animationCmpt *Component
	}
)

func New(atlasCache *texture.AtlasCache) *System {
	return &System{
		entities:   []entityAspect{},
		entityMap:  map[entity.ID]*entityAspect{},
		atlasCache: atlasCache,
	}
}

func (s *System) GetAnimation(eid entity.ID) string {
	if e, exists := s.entityMap[eid]; exists {
		return e.animationCmpt.Animation
	} else {
		panic("entity does not exist")
	}
}

func (s *System) SetAnimation(eid entity.ID, animation string, animationStart common.Time) {
	if e, exists := s.entityMap[eid]; exists {
		e.animationCmpt.Animation = animation
		e.animationCmpt.IsAnimating = true

	} else {
		panic("entity does not exist")
	}
}

func (s *System) StopAnimating(eid entity.ID) {
	if e, exists := s.entityMap[eid]; exists {
		e.animationCmpt.IsAnimating = false

	} else {
		panic("entity does not exist")
	}
}

func (s *System) Update(t common.Time) {
	for _, e := range s.entities {
		cmpt := e.animationCmpt

		if !cmpt.IsAnimating {
			continue
		}

		atlas, err := s.atlasCache.Load(cmpt.AtlasName)
		if err != nil {
			panic(err.Error())
		}

		textures := atlas.Animation(cmpt.Animation)
		if len(textures) <= 0 {
			continue
		}

		elapsedNano := t - cmpt.AnimationStart
		totalFrames := int64(math.Floor(elapsedNano.Seconds() * float64(cmpt.FPS)))
		newIndex := int(totalFrames % int64(len(textures)))

		if cmpt.CurrentIndex == 0 || newIndex != cmpt.CurrentIndex {
			cmpt.CurrentIndex = newIndex
			tex := textures[cmpt.CurrentIndex]
			e.renderCmpt.SetTexture(tex)
		}
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	em.RegisterComponentType("animation", &Component{}, nil)
}

func (s *System) ComponentsWillJoin(eid entity.ID, components []entity.IComponent) {
	var animationCmpt *Component
	var renderCmpt *rendersys.Component

	for _, cmpt := range components {
		if ac, is := cmpt.(*Component); is {
			animationCmpt = ac
		} else if rc, is := cmpt.(*rendersys.Component); is {
			renderCmpt = rc
		}

		if animationCmpt != nil && renderCmpt != nil {
			break
		}
	}

	if animationCmpt != nil {
		if renderCmpt == nil {
			panic("animation component requires render component")
		}

		s.entities = append(s.entities, entityAspect{
			id:            eid,
			animationCmpt: animationCmpt,
			renderCmpt:    renderCmpt,
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
