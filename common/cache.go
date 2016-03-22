package common

import (
	"io/ioutil"
	"os"
	"sync"
)

type (
	Cache struct {
		mutex       sync.RWMutex
		things      map[string]interface{}
		fs          ICacheFilesystem
		decodeBytes func(bs []byte) (interface{}, error)
	}

	ICacheFilesystem interface {
		OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
	}
)

func NewCache(fs ICacheFilesystem, decodeBytes func(bs []byte) (interface{}, error)) *Cache {
	return &Cache{
		mutex:       sync.RWMutex{},
		things:      map[string]interface{}{},
		fs:          fs,
		decodeBytes: decodeBytes,
	}
}

func (c *Cache) Load(filename string) (interface{}, error) {
	c.mutex.RLock()
	thing, exists := c.things[filename]
	c.mutex.RUnlock()

	if exists {
		return thing, nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	thing, err := c.loadThing(filename)
	if err != nil {
		return nil, err
	}

	c.things[filename] = thing
	return thing, nil
}

func (c *Cache) loadThing(filename string) (interface{}, error) {
	f, err := c.fs.OpenFile(filename, 0, 0400)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return c.decodeBytes(bytes)
}
