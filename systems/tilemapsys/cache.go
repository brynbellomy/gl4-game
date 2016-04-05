package tilemapsys

import (
	"fmt"

	"github.com/azul3d-legacy/tmx"

	"github.com/brynbellomy/gl4-game/common"
)

type (
	TilemapCache struct {
		*common.Cache
	}
)

func NewTilemapCache(fs common.ICacheFilesystem) *TilemapCache {
	return &TilemapCache{
		Cache: common.NewCache(fs, func(bs []byte) (interface{}, error) {
			// r := strings.NewReader(string(bs))
			return tmx.Parse(bs)
		}),
	}
}

func (c *TilemapCache) Load(filename string) (*tmx.Map, error) {
	fmt.Println("tilemap cache: loading", filename)
	theMap, err := c.Cache.Load(filename)
	if err != nil {
		return nil, err
	}

	return theMap.(*tmx.Map), nil
}
