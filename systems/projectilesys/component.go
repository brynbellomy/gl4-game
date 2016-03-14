package projectilesys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/entity"
)

type (
	Component struct {
		Heading         mgl32.Vec2      `config:"heading"`
		ExitVelocity    float32         `config:"exitVelocity"`
		Thrust          float32         `config:"thrust"`
		State           ProjectileState `config:"state"`
		RemoveOnContact bool            `config:"removeOnContact"`
	}

	ProjectileState int

	ComponentSlice []Component
)

const (
	Firing ProjectileState = iota
	Flying
	Impacting
)

func (c *Component) GetHeading() mgl32.Vec2 {
	return c.Heading
}

func (c *Component) SetHeading(heading mgl32.Vec2) {
	c.Heading = heading
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
