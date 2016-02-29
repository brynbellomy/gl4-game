package rendersys

type (
	Component struct {
		renderNode INode
		texture    uint32
	}
)

func NewComponent(renderNode INode, texture uint32) *Component {
	return &Component{renderNode, texture}
}

func (c *Component) Texture() uint32 {
	return c.texture
}

func (c *Component) SetTexture(tex uint32) {
	c.texture = tex
}
