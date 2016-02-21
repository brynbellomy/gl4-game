package node

import (
	_ "image/png"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl-test/common"
	"github.com/brynbellomy/gl-test/shader"
	"github.com/brynbellomy/gl-test/texture"
)

type (
	SpriteNode struct {
		*Node

		program  uint32 // shader program
		vao      uint32 // vertex array object
		uCamera  int32  // camera uniform
		uModel   int32  // model uniform
		uTexture int32  // texture uniform

		model mgl32.Mat4 // model matrix

		texture  uint32 // texture id
		size     *common.Size
		position mgl32.Vec2

		previousTime, totalTime float64
	}
)

type SpriteNodeConfig struct {
	Projection  mgl32.Mat4
	Size        *common.Size
	Position    mgl32.Vec2
	TextureFile string
}

func NewSpriteNode(config SpriteNodeConfig) *SpriteNode {
	// Configure the vertex and fragment shaders
	program, err := shader.NewProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	// projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &config.Projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{0, 0, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Load the texture
	texture, err := texture.New(config.TextureFile)
	if err != nil {
		panic(err)
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Configure the vertex data
	vertices := common.Rect(config.Size)

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
		model:    model,
		texture:  texture,
		Node: New(Config{
			Size:     config.Size,
			Position: config.Position,
		}),
	}
}

func (n *SpriteNode) Update() {
	n.Node.Update()

	n.model = mgl32.Translate3D(n.Node.Position().X(), n.Node.Position().Y(), 0.0)

	// Update
	gl.UseProgram(n.program)
	gl.UniformMatrix4fv(n.uModel, 1, false, &n.model[0])

	// camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	// gl.UniformMatrix4fv(n.uCamera, 1, false, &n.camera[0])

	// x := float32(rand.Intn(1000)) / 1000
	// n.size = &common.Size{Width: x, Height: x}

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
