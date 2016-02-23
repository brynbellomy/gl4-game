package rendersys

import "github.com/brynbellomy/gl4-game/node"

type (
	Component struct {
		renderNode *node.SpriteNode
		texture    uint32
	}
)

func NewComponent(renderNode *node.SpriteNode, texture uint32) *Component {
	return &Component{renderNode, texture}
}

func (c *Component) Texture() uint32 {
	return c.texture
}

func (c *Component) SetTexture(tex uint32) {
	c.texture = tex
}
