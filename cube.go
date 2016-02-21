package main

import (
	_ "image/png"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl-test/shader"
	"github.com/brynbellomy/gl-test/texture"
)

type Cube struct {
	previousTime float64
	angle        float64
	program      uint32

	// vertex array object
	vao uint32

	// uniforms
	uCamera  int32
	uModel   int32
	uTexture int32

	// matrices
	model mgl32.Mat4

	// texture
	texture uint32
}

func NewCube() *Cube {
	// Configure the vertex and fragment shaders
	program, err := shader.NewProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 5, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Load the texture
	texture, err := texture.New("square.png")
	if err != nil {
		panic(err)
	}

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	return &Cube{
		previousTime: glfw.GetTime(),
		program:      program,
		vao:          vao,
		uCamera:      cameraUniform,
		uModel:       modelUniform,
		uTexture:     textureUniform,
		model:        model,
		texture:      texture,
		angle:        0.0,
	}
}

func (c *Cube) Render() {
	// Update
	time := glfw.GetTime()
	elapsed := time - c.previousTime
	c.previousTime = time

	c.angle += elapsed
	c.model = mgl32.HomogRotate3D(float32(c.angle), mgl32.Vec3{0, 1, 0})

	// Render
	gl.UseProgram(c.program)
	gl.UniformMatrix4fv(c.uModel, 1, false, &c.model[0])

	// camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	// gl.UniformMatrix4fv(c.uCamera, 1, false, &c.camera[0])

	gl.BindVertexArray(c.vao)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, c.texture)

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

var cubeVertices = []float32{
	//  X, Y, Z, U, V

	// // Bottom
	// -2.0, -1.0, -1.0, 0.0, 0.0,
	// 1.0, -1.0, -1.0, 1.0, 0.0,
	// -1.0, -1.0, 1.0, 0.0, 1.0,
	// 1.0, -1.0, -1.0, 1.0, 0.0,
	// 1.0, -1.0, 1.0, 1.0, 1.0,
	// -1.0, -1.0, 1.0, 0.0, 1.0,

	// // Top
	// -1.0, 1.0, -1.0, 0.0, 0.0,
	// -1.0, 1.0, 1.0, 0.0, 1.0,
	// 1.0, 1.0, -1.0, 1.0, 0.0,
	// 1.0, 1.0, -1.0, 1.0, 0.0,
	// -1.0, 1.0, 1.0, 0.0, 1.0,
	// 1.0, 1.0, 1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,

	// // Back
	// -1.0, -1.0, -1.0, 0.0, 0.0,
	// -1.0, 1.0, -1.0, 0.0, 1.0,
	// 1.0, -1.0, -1.0, 1.0, 0.0,
	// 1.0, -1.0, -1.0, 1.0, 0.0,
	// -1.0, 1.0, -1.0, 0.0, 1.0,
	// 1.0, 1.0, -1.0, 1.0, 1.0,

	// // Left
	// -1.0, -1.0, 1.0, 0.0, 1.0,
	// -1.0, 1.0, -1.0, 1.0, 0.0,
	// -1.0, -1.0, -1.0, 0.0, 0.0,
	// -1.0, -1.0, 1.0, 0.0, 1.0,
	// -1.0, 1.0, 1.0, 1.0, 1.0,
	// -1.0, 1.0, -1.0, 1.0, 0.0,

	// // Right
	// 1.0, -1.0, 1.0, 1.0, 1.0,
	// 1.0, -1.0, -1.0, 1.0, 0.0,
	// 1.0, 1.0, -1.0, 0.0, 0.0,
	// 1.0, -1.0, 1.0, 1.0, 1.0,
	// 1.0, 1.0, -1.0, 0.0, 0.0,
	// 1.0, 1.0, 1.0, 0.0, 1.0,
}
