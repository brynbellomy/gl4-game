package animationsys

import "github.com/brynbellomy/gl4-game/common"

type (
	Component struct {
		AtlasName      string      `config:"atlasName"`
		Animation      string      `config:"animation"`
		IsAnimating    bool        `config:"isAnimating"`
		CurrentIndex   int         `config:"currentIndex"`
		AnimationStart common.Time `config:"animationStart"`
		FPS            int         `config:"fps"`
	}
)

func NewComponent(atlasName string, animation string, isAnimating bool, currentIndex int, fps int) *Component {
	return &Component{
		AtlasName:      atlasName,
		IsAnimating:    isAnimating,
		Animation:      animation,
		CurrentIndex:   currentIndex,
		AnimationStart: common.Time(0),
		FPS:            fps,
	}
}

func (c *Component) SetFPS(fps int) {
	c.FPS = fps
}

func (c *Component) SetAnimation(animation string) {
	c.Animation = animation
}

func (c *Component) SetIsAnimating(is bool) {
	c.IsAnimating = is
}
