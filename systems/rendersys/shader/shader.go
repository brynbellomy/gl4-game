package shader

import (
	"fmt"
	_ "image/png"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type (
	Shader  uint32
	Program uint32

	ShaderType uint32
)

const (
	VertexShader   ShaderType = gl.VERTEX_SHADER
	FragmentShader ShaderType = gl.FRAGMENT_SHADER
)

func NewProgram(vertexShader, fragmentShader Shader) (Program, error) {
	program := gl.CreateProgram()

	gl.AttachShader(program, uint32(vertexShader))
	gl.AttachShader(program, uint32(fragmentShader))
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(uint32(vertexShader))
	gl.DeleteShader(uint32(fragmentShader))

	return Program(program), nil
}

func compileShader(source string, shaderType ShaderType) (Shader, error) {
	shader := gl.CreateShader(uint32(shaderType))

	csource, freeCsource := gl.Strs(source)
	gl.ShaderSource(shader, 1, csource, nil)
	freeCsource()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return Shader(shader), nil
}
