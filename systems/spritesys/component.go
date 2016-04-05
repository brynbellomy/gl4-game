package spritesys

import (
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/rendersys/texture"
)

type (
	Component struct {
		Texture texture.TextureID
	}

	ComponentCfg struct {
		TextureName string `config:"texture"`
	}

	ComponentSlice []Component
)

func (c *Component) GetTexture() texture.TextureID {
	return c.Texture
}

func (c *Component) SetTexture(tex texture.TextureID) {
	c.Texture = tex
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
