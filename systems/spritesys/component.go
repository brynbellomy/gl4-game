package spritesys

type (
	Component struct {
		textureName     string
		isTextureLoaded bool
		texture         uint32
	}
)

func NewComponent(textureName string) *Component {
	return &Component{textureName, false, 0}
}

func (c *Component) GetTexture() uint32 {
	return c.texture
}

func (c *Component) SetTexture(tex uint32) {
	c.texture = tex
}

func (c *Component) GetTextureName() string {
	return c.textureName
}

func (c *Component) SetTextureName(tex string) {
	c.textureName = tex
}
