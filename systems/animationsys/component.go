package animationsys

import (
	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	Component struct {
		AtlasName      string      `config:"atlasName"`
		Animation      string      `config:"animation"`
		IsAnimating    bool        `config:"isAnimating"`
		CurrentIndex   int         `config:"currentIndex"`
		AnimationStart common.Time `config:"-"`
		FPS            int         `config:"fps"`
	}

	ComponentSlice []Component
)

func (c *Component) SetFPS(fps int) {
	c.FPS = fps
}

func (c *Component) SetAnimation(animation string) {
	c.Animation = animation
}

func (c *Component) SetIsAnimating(is bool) {
	c.IsAnimating = is
}

func (c Component) Clone() entity.IComponent {
	return c
}

func (cs ComponentSlice) Append(cmpt entity.IComponent) entity.IComponentSlice {
	return append(cs, cmpt.(Component))
}

func (cs ComponentSlice) Remove(idx int) entity.IComponentSlice {
	return append(cs[:idx], cs[idx+1:]...)
}

func (cs ComponentSlice) Get(idx int) (entity.IComponent, bool) {
	if idx >= len(cs) {
		return nil, false
	}
	return cs[idx], true
}

func (cs ComponentSlice) Set(idx int, cmpt entity.IComponent) bool {
	if idx >= len(cs) {
		return false
	}
	cs[idx] = cmpt.(Component)
	return true
}
