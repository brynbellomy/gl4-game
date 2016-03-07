package texture

import (
	"fmt"
	"sync"

	"github.com/brynbellomy/gl4-game/systems/assetsys"
)

type AtlasCache struct {
	mutex        sync.RWMutex
	atlasses     map[string]*Atlas
	textureCache *TextureCache
	fs           assetsys.IFilesystem
}

func NewAtlasCache(textureCache *TextureCache, fs assetsys.IFilesystem) *AtlasCache {
	return &AtlasCache{
		mutex:        sync.RWMutex{},
		atlasses:     map[string]*Atlas{},
		textureCache: textureCache,
		fs:           fs,
	}
}

func (c *AtlasCache) Load(filename string) (*Atlas, error) {
	fmt.Println("atlas cache: loading", filename)

	c.mutex.RLock()
	atlas, exists := c.atlasses[filename]
	c.mutex.RUnlock()

	if exists {
		return atlas, nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	atlas, err := NewAtlasFromFile(filename, c.fs, c.textureCache)
	if err != nil {
		return nil, err
	}

	c.atlasses[filename] = atlas
	return atlas, nil
}
