package texture

import "fmt"

type Atlas struct {
	textures map[string][]uint32
}

func NewAtlas() *Atlas {
	return &Atlas{
		textures: map[string][]uint32{},
	}
}

func (a *Atlas) Animation(name string) []uint32 {
	return a.textures[name]
}

func (a *Atlas) LoadAnimation(name string, filenames []string) error {
    fmt.Println("Loading animation", name, "...")
	textures := make([]uint32, len(filenames))

	for i, filename := range filenames {
		tex, err := globalTextureCache.Load(filename)
		if err != nil {
			return err
		}

		textures[i] = tex
	}

	a.textures[name] = textures

    fmt.Println("Loaded animation", name, "=", textures)

    return nil
}
