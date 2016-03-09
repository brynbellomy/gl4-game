package rendersys

import (
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	Component struct {
		renderNode INode  `config:"-"`
		texture    uint32 `config:"-"`

		NodeType   string                 `config:"nodeType"`
		NodeConfig map[string]interface{} `config:"nodeConfig"`
	}
)

func (c *Component) Texture() uint32 {
	return c.texture
}

func (c *Component) SetTexture(tex uint32) {
	c.texture = tex
}

func (c *Component) Clone() entity.IComponent {
	x := *c
	return &x
}
