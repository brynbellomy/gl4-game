package movesys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/entity"
)

type (
	Component struct {
		Vec mgl32.Vec2 `config:"vector"`
	}

	ComponentSlice []Component
)

func NewComponent(vec mgl32.Vec2) *Component {
	return &Component{Vec: vec}
}

func (c *Component) Vector() mgl32.Vec2 {
	return c.Vec
}

func (c *Component) SetVector(vec mgl32.Vec2) {
	c.Vec = vec
}

func (c *Component) ResetVector() {
	c.Vec = mgl32.Vec2{0, 0}
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
