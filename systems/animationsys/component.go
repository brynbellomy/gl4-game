package animationsys

import (
	"time"

	"github.com/brynbellomy/gl4-game/texture"
)

type (
	Component struct {
		atlas               *texture.Atlas
		animation           string
		animationHasChanged bool
		currentIndex        int
		animationStart      time.Time
		fps                 int
	}
)

func NewComponent(atlas *texture.Atlas, animation string, currentIndex int, fps int) *Component {
	return &Component{
		atlas:          atlas,
		animation:      animation,
		currentIndex:   currentIndex,
		animationStart: time.Time{},
		fps:            fps,
	}
}
