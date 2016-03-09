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
)

func NewComponent(textureName string) *Component {
	return &Component{textureName, false, 0}
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

func (c *Component) Clone() entity.IComponent {
	x := *c
	return &x
}
