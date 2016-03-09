package shader

import (
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/brynbellomy/gl4-game/systems/assetsys"
)

type (
	ShaderCache struct {
		mutex   sync.RWMutex
		shaders map[string]Shader
		fs      assetsys.IFilesystem
	}
)

func NewShaderCache(fs assetsys.IFilesystem) *ShaderCache {
	return &ShaderCache{
		mutex:   sync.RWMutex{},
		shaders: map[string]Shader{},
		fs:      fs,
	}
}

func (c *ShaderCache) Load(filename string, shaderType ShaderType) (Shader, error) {
	fmt.Println("shader cache: loading", filename)

	c.mutex.RLock()
	shader, exists := c.shaders[filename]
	c.mutex.RUnlock()

	if exists {
		return shader, nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	shader, err := c.loadShader(filename, shaderType)
	if err != nil {
		return 0, err
	}

	c.shaders[filename] = shader
	return shader, nil
}

func (c *ShaderCache) loadShader(filename string, shaderType ShaderType) (Shader, error) {
	f, err := c.fs.OpenFile(filename, 0, 0400)
	if err != nil {
		return 0, err
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return 0, err
	}

	bytes = append(bytes, '\x00')

	return compileShader(string(bytes), shaderType)
}
