package shader

import (
	"fmt"
	"sync"
)

// var globalProgramCache = NewProgramCache()

// func LoadProgram(vertexShaderFile, fragmentShaderFile string) (Program, error) {
// 	return globalProgramCache.LoadProgram(vertexShaderFile, fragmentShaderFile)
// }

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

func (c *ProgramCache) LoadProgram(vertexShaderFile, fragmentShaderFile string) (Program, error) {
	fmt.Println("program cache: loading", "[", vertexShaderFile, "+", fragmentShaderFile, "]")

	key := programKey{vertexShaderFile, fragmentShaderFile}

	c.mutex.RLock()
	program, exists := c.programs[key]
	c.mutex.RUnlock()

	if exists {
		return program, nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	vs, err := c.shaderCache.LoadShader(vertexShaderFile, VertexShader)
	if err != nil {
		return 0, err
	}

	fs, err := c.shaderCache.LoadShader(fragmentShaderFile, FragmentShader)
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
