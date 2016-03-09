package shader

import (
	"fmt"
	"sync"
)

type (
	ProgramCache struct {
		mutex       sync.RWMutex
		programs    map[programKey]Program
		shaderCache *ShaderCache
	}

	programKey struct{ vsfile, fsfile string }
)

func NewProgramCache(shaderCache *ShaderCache) *ProgramCache {
	return &ProgramCache{
		mutex:       sync.RWMutex{},
		programs:    map[programKey]Program{},
		shaderCache: shaderCache,
	}
}

func (c *ProgramCache) Load(vertexShaderFile, fragmentShaderFile string) (Program, error) {
	key := programKey{vertexShaderFile, fragmentShaderFile}

	c.mutex.RLock()
	program, exists := c.programs[key]
	c.mutex.RUnlock()

	if exists {
		return program, nil
	}

	fmt.Println("program cache: loading", "[", vertexShaderFile, "+", fragmentShaderFile, "]")

	c.mutex.Lock()
	defer c.mutex.Unlock()

	vs, err := c.shaderCache.Load(vertexShaderFile, VertexShader)
	if err != nil {
		return 0, err
	}

	fs, err := c.shaderCache.Load(fragmentShaderFile, FragmentShader)
	if err != nil {
		return 0, err
	}

	program, err = NewProgram(vs, fs)
	if err != nil {
		return 0, err
	}

	c.programs[key] = program
	return program, nil
}
