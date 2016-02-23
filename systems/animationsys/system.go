package animationsys

import (
	"math"
	"time"

	"github.com/brynbellomy/gl4-game/context"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/rendersys"
)

type (
	System struct {
		entities []entityAspect
	}

	entityAspect struct {
		renderCmpt    *rendersys.Component
		animationCmpt *Component
	}
)

func New() *System {
	return &System{
		entities: make([]entityAspect, 0),
	}
}

func (s *System) Update(c context.IContext) {
	for _, e := range s.entities {
		cmpt := e.animationCmpt

		textures := cmpt.atlas.Animation(cmpt.animation)
		if len(textures) <= 0 {
			continue
		}

		elapsed := time.Now().Sub(cmpt.animationStart)
		totalFrames := int64(math.Floor(elapsed.Seconds() * float64(cmpt.fps)))
		newIndex := int(totalFrames % int64(len(textures)))

		if cmpt.currentIndex == 0 || newIndex != cmpt.currentIndex {
			cmpt.currentIndex = newIndex
			tex := textures[cmpt.currentIndex]
			e.renderCmpt.SetTexture(tex)
		}
	}
}

func (s *System) EntityWillJoin(eid entity.ID, components []entity.IComponent) {
	// if we find a *animationsys.Component on the entity, we keep track of it
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

		animationCmpt.animationStart = time.Now()

		aspect := entityAspect{animationCmpt: animationCmpt, renderCmpt: renderCmpt}
		s.entities = append(s.entities, aspect)
	}
}