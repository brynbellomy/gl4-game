package animationsys

import (
	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/texture"
)

type (
	Component struct {
		atlas               *texture.Atlas
		isAnimating         bool
		animation           string
		animationHasChanged bool
		currentIndex        int
		animationStart      common.Time
		fps                 int
	}
)

func NewComponent(atlas *texture.Atlas, animation string, currentIndex int, fps int) *Component {
	return &Component{
		atlas:          atlas,
		isAnimating:    false,
		animation:      animation,
		currentIndex:   currentIndex,
		animationStart: common.Time(0),
		fps:            fps,
	}
}
