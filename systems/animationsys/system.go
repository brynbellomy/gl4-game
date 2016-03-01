package animationsys

import (
	"math"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/rendersys"
)

type (
	System struct {
		entities  []entityAspect
		entityMap map[entity.ID]*entityAspect
	}

	entityAspect struct {
		renderCmpt    *rendersys.Component
		animationCmpt *Component
	}
)

func New() *System {
	return &System{
		entities:  make([]entityAspect, 0),
		entityMap: make(map[entity.ID]*entityAspect),
	}
}

func (s *System) GetAnimation(eid entity.ID) string {
	if e, exists := s.entityMap[eid]; exists {
		return e.animationCmpt.animation
	} else {
		panic("entity does not exist")
	}
}

func (s *System) SetAnimation(eid entity.ID, animation string, animationStart common.Time) {
	if e, exists := s.entityMap[eid]; exists {
		e.animationCmpt.animation = animation
		e.animationCmpt.isAnimating = true

	} else {
		panic("entity does not exist")
	}
}

func (s *System) StopAnimating(eid entity.ID) {
	if e, exists := s.entityMap[eid]; exists {
		e.animationCmpt.isAnimating = false

	} else {
		panic("entity does not exist")
	}
}

func (s *System) Update(t common.Time) {
	for _, e := range s.entities {
		cmpt := e.animationCmpt

		if !cmpt.isAnimating {
			continue
		}

		textures := cmpt.atlas.Animation(cmpt.animation)
		if len(textures) <= 0 {
			continue
		}

		elapsedNano := t - cmpt.animationStart
		totalFrames := int64(math.Floor(elapsedNano.Seconds() * float64(cmpt.fps)))
		newIndex := int(totalFrames % int64(len(textures)))

		if cmpt.currentIndex == 0 || newIndex != cmpt.currentIndex {
			cmpt.currentIndex = newIndex
			tex := textures[cmpt.currentIndex]
			e.renderCmpt.SetTexture(tex)
		}
	}
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
			animationCmpt: animationCmpt,
			renderCmpt:    renderCmpt,
		})

		s.entityMap[eid] = &s.entities[len(s.entities)-1]
	}
}
