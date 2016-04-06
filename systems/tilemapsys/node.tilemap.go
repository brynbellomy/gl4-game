package tilemapsys

import (
	"errors"
	_ "image/png"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/azul3d-legacy/tmx"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/systems/rendersys"
	"github.com/brynbellomy/gl4-game/systems/rendersys/shader"
	"github.com/brynbellomy/gl4-game/systems/rendersys/texture"
)

type (
	TilemapNodeFactory struct {
		shaderProgramCache *shader.ProgramCache
		textureCache       *texture.TextureCache
	}
)

func (f *TilemapNodeFactory) NodeFromConfig(config map[string]interface{}) (rendersys.INode, error) {
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

	return &TilemapNode{program: program, textureCache: f.textureCache}, nil
}

type (
	TilemapNode struct {
		tilemap               *tmx.Map
		loadedVertexSublayers []loadedVertexSublayer
		textureCache          *texture.TextureCache

		program     shader.Program // shader program
		uCamera     int32          // camera uniform
		uModel      int32          // model uniform
		uProjection int32          // projection uniform

		// texture  texture.TextureID // texture id
		size     common.Size
		position mgl32.Vec2
		rotation float32
	}
)

func (n *TilemapNode) SetTexture(tid texture.TextureID) {
	// no-op
}

func (n *TilemapNode) SetTilemap(tm *tmx.Map) {
	n.tilemap = tm
}

func (n *TilemapNode) Init() error {
	program := uint32(n.program)

	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// VAO initialization
	for _, sublayer := range VerticesFromTilemap(n.tilemap) {
		var vao, vbo uint32

		gl.GenVertexArrays(1, &vao)
		gl.BindVertexArray(vao)

		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(sublayer.vertices)*4, gl.Ptr(sublayer.vertices), gl.STATIC_DRAW)

		tex, err := n.textureCache.Load(sublayer.texture)
		if err != nil {
			return err
		}

		n.loadedVertexSublayers = append(n.loadedVertexSublayers, loadedVertexSublayer{
			vao:         vao,
			texture:     tex,
			numVertices: len(sublayer.vertices),
		})

		// the following code describes the format of the buffer data and stores this description
		// inside the VAO object (so that all we have to do is call gl.BindVertexArray() during render)
		vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
		gl.EnableVertexAttribArray(vertAttrib)
		gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

		texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
		gl.EnableVertexAttribArray(texCoordAttrib)
		gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	}

	n.uCamera = cameraUniform
	n.uModel = modelUniform
	n.uProjection = projectionUniform

	return nil
}

// func (n *TilemapNode) SetVertices(vs []float32) {
// 	n.vertices = vs
// }

func (n *TilemapNode) Render(c rendersys.RenderContext) {
	gl.UseProgram(uint32(n.program))

	trans := mgl32.Translate3D(n.position.X(), n.position.Y(), 0.0)
	rotate := mgl32.Rotate3DZ(n.rotation).Mat4()
	scale := mgl32.Scale3D(n.size.Width(), n.size.Height(), 1.0)

	model := trans.Mul4(rotate).Mul4(scale)

	gl.UniformMatrix4fv(n.uModel, 1, false, &model[0])
	gl.UniformMatrix4fv(n.uCamera, 1, false, &c.Camera[0])
	gl.UniformMatrix4fv(n.uProjection, 1, false, &c.Projection[0])

	for _, layer := range n.loadedVertexSublayers {
		gl.BindVertexArray(layer.vao)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, uint32(layer.texture))

		gl.DrawArrays(gl.TRIANGLES, 0, int32(layer.numVertices))
	}
}

func (n *TilemapNode) Destroy() error {
	// @@TODO
	// release textures, etc?
	for _, layer := range n.loadedVertexSublayers {
		gl.DeleteBuffers(1, &layer.vao)
	}
	n.loadedVertexSublayers = nil
	return nil
}

func (n *TilemapNode) SetPos(pos mgl32.Vec2) {
	n.position = pos
}

func (n *TilemapNode) SetSize(size common.Size) {
	n.size = size
}

func (n *TilemapNode) SetRotation(rotation float32) {
	n.rotation = rotation
}
