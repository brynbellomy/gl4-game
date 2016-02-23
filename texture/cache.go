package texture

import (
	"fmt"
	"sync"
)

var globalTextureCache = NewCache()

func Load(filename string) (uint32, error) {
	return globalTextureCache.Load(filename)
}

type Cache struct {
	mutex    sync.RWMutex
	textures map[string]uint32
}

func NewCache() *Cache {
	return &Cache{
		mutex:    sync.RWMutex{},
		textures: map[string]uint32{},
	}
}

func (c *Cache) Load(filename string) (uint32, error) {
	fmt.Println("texture cache: loading", filename)

	c.mutex.RLock()
	t, exists := c.textures[filename]
	c.mutex.RUnlock()

	if exists {
		return t, nil

	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	t, err := New(filename)
	if err != nil {
		return 0, err
	}

	c.textures[filename] = t
	return t, nil
}
