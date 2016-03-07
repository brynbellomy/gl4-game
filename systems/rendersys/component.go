package rendersys

import (
	"github.com/brynbellomy/gl4-game/systems/rendersys/shader"
)

type (
	Component struct {
		renderNode INode
		texture    uint32

		vertexShaderFile    string
		fragmentShaderFile  string
		shaderProgramLoaded bool
		shaderProgram       shader.Program
	}
)

func NewComponent(renderNode INode, vertexShaderFile string, fragmentShaderFile string) *Component {
	return &Component{
		renderNode:         renderNode,
		texture:            0,
		vertexShaderFile:   vertexShaderFile,
		fragmentShaderFile: fragmentShaderFile,
	}
}

func (c *Component) Texture() uint32 {
	return c.texture
}

func (c *Component) SetTexture(tex uint32) {
	c.texture = tex
}

func (c *Component) ShaderProgram() shader.Program {
	return c.shaderProgram
}

func (c *Component) SetShaderProgram(sp shader.Program) {
	c.shaderProgram = sp
}

// func (c *Component) TextureName() string {
// 	return c.textureName
// }

// func (c *Component) SetTextureName(tex string) {
// 	c.textureName = tex
// }
