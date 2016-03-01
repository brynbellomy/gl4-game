package rendersys

import (
	_ "image/png"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/shader"
)

type (
	SpriteNode struct {
		program     uint32 // shader program
		vao         uint32 // vertex array object
		uCamera     int32  // camera uniform
		uModel      int32  // model uniform
		uProjection int32  // texture uniform
		uTexture    int32  // texture uniform

		texture  uint32 // texture id
		size     common.Size
		position mgl32.Vec2
	}
)

func NewSpriteNode() *SpriteNode {
	// Configure the vertex and fragment shaders
	program, err := shader.NewProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	defaultProjection := mgl32.Ident4()
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &defaultProjection[0])

	// defaultCamera := mgl32.LookAtV(mgl32.Vec3{0, 0, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, -1, 0})
	defaultCamera := mgl32.Ident4()
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &defaultCamera[0])

	defaultModel := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &defaultModel[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Configure the vertex data
	vertices := common.Rect(common.Size{1.0, 1.0})

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	return &SpriteNode{
		program:  program,
		vao:      vao,
		uCamera:  cameraUniform,
		uModel:   modelUniform,
		uTexture: textureUniform,
	}
}

func (n *SpriteNode) SetPos(pos mgl32.Vec2) {
	n.position = pos
}

func (n *SpriteNode) SetSize(size common.Size) {
	n.size = size
}

func (n *SpriteNode) SetTexture(tex uint32) {
	n.texture = tex
}

func (n *SpriteNode) Render(c RenderContext) {
	gl.UseProgram(n.program)

	trans := mgl32.Translate3D(n.position.X(), n.position.Y(), 0.0)
	scale := mgl32.Scale3D(n.size.Width(), n.size.Height(), 1.0)
	model := trans.Mul4(scale)

	gl.UniformMatrix4fv(n.uModel, 1, false, &model[0])
	gl.UniformMatrix4fv(n.uCamera, 1, false, &c.Camera[0])
	gl.UniformMatrix4fv(n.uProjection, 1, false, &c.Projection[0])

	gl.BindVertexArray(n.vao)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, n.texture)

	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
}

var vertexShader = `
#version 410

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
}
` + "\x00"

var fragmentShader = `
#version 410

uniform sampler2D tex;

in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
    outputColor = texture(tex, fragTexCoord);
}
` + "\x00"
