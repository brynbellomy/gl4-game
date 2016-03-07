package shader

import (
	"fmt"
	"sync"
)

// var globalShaderCache = NewShaderCache()

// func LoadShader(filename string, shaderType ShaderType) (Shader, error) {
// 	return globalShaderCache.LoadShader(filename, shaderType)
// }

type (
	ShaderCache struct {
		mutex   sync.RWMutex
		shaders map[string]Shader
	}
)

func NewShaderCache() *ShaderCache {
	return &ShaderCache{
		mutex:   sync.RWMutex{},
		shaders: map[string]Shader{},
	}
}

func (c *ShaderCache) LoadShader(filename string, shaderType ShaderType) (Shader, error) {
	fmt.Println("shader cache: loading", filename)

	c.mutex.RLock()
	shader, exists := c.shaders[filename]
	c.mutex.RUnlock()

	if exists {
		return shader, nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	shader, err := NewShader(filename, shaderType)
	if err != nil {
		return 0, err
	}

	c.shaders[filename] = shader
	return shader, nil
}
