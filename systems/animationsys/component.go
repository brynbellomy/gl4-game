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

func (c *Component) Clone() entity.IComponent {
	x := *c
	return &x
}
