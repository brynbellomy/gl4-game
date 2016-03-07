package texture

import (
	"errors"
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"

	"github.com/brynbellomy/gl4-game/systems/assetsys"
)

type Atlas struct {
	name       string
	animations map[string][]uint32
}

func NewAtlasFromFile(filename string, fs assetsys.IFilesystem, textureCache *TextureCache) (*Atlas, error) {
	configFile, err := fs.OpenFile(path.Join(filename, "atlas.yaml"), 0, 0400)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	err = yaml.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}

	return NewAtlasFromConfig(textureCache, m)
}

func NewAtlasFromConfig(cache *TextureCache, config map[string]interface{}) (*Atlas, error) {
	name, exists := config["name"].(string)
	if !exists {
		return nil, errors.New("missing required key 'name' (or wrong type)")
	}

	anims, exists := config["animations"].(map[string]interface{})
	if !exists {
		return nil, errors.New("missing required key 'animations' (or wrong type)")
	}

	animations := make(map[string][]uint32)
	for animName, frames := range anims {
		frames, is := frames.([]interface{})
		if !is {
			return nil, errors.New("animation frames must be a list of strings")
		}

		texs := make([]uint32, len(frames))
		for i, filename := range frames {
			if filename, is := filename.(string); is {
				tex, err := cache.Load(filename)
				if err != nil {
					return nil, err
				}
				texs[i] = tex
			} else {
				return nil, errors.New("animation frames must be a list of strings")
			}
		}
		animations[animName] = texs
	}

	return &Atlas{
		name:       name,
		animations: animations,
	}, nil
}

func (a *Atlas) Name() string {
	return a.name
}

func (a *Atlas) Animation(name string) []uint32 {
	return a.animations[name]
}
