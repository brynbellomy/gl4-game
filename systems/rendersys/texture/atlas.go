package texture

import (
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"

	"github.com/brynbellomy/go-structomancer"

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

	return NewAtlasFromConfig(textureCache, filename, m)
}

type atlasConfig struct {
	Name       string              `config:"name"`
	Animations map[string][]string `config:"animations"`
}

var atlasStruct = structomancer.New(&atlasConfig{}, "config")

func NewAtlasFromConfig(cache *TextureCache, assetRoot string, config map[string]interface{}) (*Atlas, error) {
	c, err := atlasStruct.MapToStruct(config)
	if err != nil {
		return nil, err
	}
	cfg := c.(*atlasConfig)

	animations := make(map[string][]uint32)

	for animName, frames := range cfg.Animations {
		texs := make([]uint32, len(frames))

		for i, filename := range frames {
			tex, err := cache.Load(path.Join(assetRoot, filename))
			if err != nil {
				return nil, err
			}
			texs[i] = tex
		}
		animations[animName] = texs
	}

	return &Atlas{
		name:       cfg.Name,
		animations: animations,
	}, nil
}

func (a *Atlas) Name() string {
	return a.name
}

func (a *Atlas) Animation(name string) []uint32 {
	return a.animations[name]
}
