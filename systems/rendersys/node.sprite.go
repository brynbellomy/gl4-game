package rendersys

import (
	"errors"
	_ "image/png"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/systems/rendersys/shader"
	"github.com/brynbellomy/gl4-game/systems/rendersys/texture"
)

type (
	SpriteNodeFactory struct {
		shaderProgramCache *shader.ProgramCache
	}
)

func (f *SpriteNodeFactory) NodeFromConfig(config map[string]interface{}) (INode, error) {
	v64s, exists := config["vertices"].([]float64)
	var vertices []float32
	if !exists {
		// return nil, errors.New("missing required key 'vertices' (or wrong type)")
	} else {
		for i := range v64s {
			vertices[i] = float32(v64s[i])
		}
	}

	vertexShader, exists := config["vertex-shader"].(string)
	if !exists {
		return nil, errors.New("missing required key 'vertex-shader' (or wrong type)")
	}

	fragmentShader, exists := config["fragment-shader"].(string)
	if !exists {
		return nil, errors.New("missing required key 'fragment-shader' (or wrong type)")
	}

	program, err := f.shaderProgramCache.Load(vertexShader, fragmentShader)
	if err != nil {
		return nil, err
	}

	return &SpriteNode{program: program, vertices: vertices}, nil
}

type (
	SpriteNode struct {
		vertices    []float32
		program     shader.Program // shader program
		vao         uint32         // vertex array object
		uCamera     int32          // camera uniform
		uModel      int32          // model uniform
		uProjection int32          // texture uniform
		// uTexture    int32          // texture uniform

		texture  texture.TextureID // texture id
		size     common.Size
		position mgl32.Vec2
		rotation float32
	}
)

func (n *SpriteNode) Init() error {
	program := uint32(n.program)

	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// VAO initialization
	{
		var vao, vbo uint32

		gl.GenVertexArrays(1, &vao)
		gl.BindVertexArray(vao)

		// Configure the vertex data
		// vertices := common.Rect(common.Size{1.0, 1.0})

		if n.vertices == nil {
			n.vertices = common.Rect(common.Size{1.0, 1.0})
		}

		gl.GetError()

		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(n.vertices)*4, gl.Ptr(n.vertices), gl.STATIC_DRAW)

		n.vao = vao
	}

	// the following code describes the format of the buffer data and stores this description
	// inside the VAO object (so that all we have to do is call gl.BindVertexArray() during render)
	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	n.uCamera = cameraUniform
	n.uModel = modelUniform
	n.uProjection = projectionUniform

	return nil
}

func (n *SpriteNode) SetVertices(vs []float32) {
	n.vertices = vs
}

func (n *SpriteNode) Render(c RenderContext) {
	gl.UseProgram(uint32(n.program))

	trans := mgl32.Translate3D(n.position.X(), n.position.Y(), 0.0)
	rotate := mgl32.Rotate3DZ(n.rotation).Mat4()
	scale := mgl32.Scale3D(n.size.Width(), n.size.Height(), 1.0)

	model := trans.Mul4(rotate).Mul4(scale)

	gl.UniformMatrix4fv(n.uModel, 1, false, &model[0])
	gl.UniformMatrix4fv(n.uCamera, 1, false, &c.Camera[0])
	gl.UniformMatrix4fv(n.uProjection, 1, false, &c.Projection[0])

	gl.BindVertexArray(n.vao)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, uint32(n.texture))

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(n.vertices))) //6*2*3)
}

func (n *SpriteNode) Destroy() error {
	// @@TODO
	// release textures, etc?
	gl.DeleteBuffers(1, &n.vao)
	return nil
}

func (n *SpriteNode) SetPos(pos mgl32.Vec2) {
	n.position = pos
}

func (n *SpriteNode) SetSize(size common.Size) {
	n.size = size
}

func (n *SpriteNode) SetRotation(rotation float32) {
	n.rotation = rotation
}

func (n *SpriteNode) SetTexture(tex texture.TextureID) {
	n.texture = tex
}
