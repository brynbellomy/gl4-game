package scriptsys

import (
    "fmt"

    "github.com/brynbellomy/gl4-game/common"
)

type (
    ScriptCache struct {
        *common.Cache
    }
)

func NewScriptCache(fs common.ICacheFilesystem) *ScriptCache {
    return &ScriptCache{
        Cache: common.NewCache(fs, func(bs []byte) (interface{}, error) {
            return string(bs), nil
        }),
    }
}

func (c *ScriptCache) Load(filename string) (string, error) {
    fmt.Println("script cache: loading", filename)
    scr, err := c.Cache.Load(filename)
    if err != nil {
        return "", err
    }

    return scr.(string), nil
}

