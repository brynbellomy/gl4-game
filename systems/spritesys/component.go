package spritesys

import (
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	Component struct {
		TextureName     string `config:"texture"`
		IsTextureLoaded bool   `config:"-"`
		Texture         uint32 `config:"-"`
	}

	ComponentSlice []Component
)

func NewComponent(textureName string) *Component {
	return &Component{
		TextureName:     textureName,
		IsTextureLoaded: false,
		Texture:         0,
	}
}

func (c *Component) GetTexture() uint32 {
	return c.Texture
}

func (c *Component) SetTexture(tex uint32) {
	c.Texture = tex
}

func (c *Component) GetTextureName() string {
	return c.TextureName
}

func (c *Component) SetTextureName(tex string) {
	c.TextureName = tex
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
