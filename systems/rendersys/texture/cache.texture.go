package texture

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/brynbellomy/gl4-game/systems/assetsys"
)

type (
	TextureCache struct {
		mutex    sync.RWMutex
		textures map[string]TextureID
		fs       assetsys.IFilesystem
	}

	TextureID uint32
)

func NewTextureCache(fs assetsys.IFilesystem) *TextureCache {
	return &TextureCache{
		mutex:    sync.RWMutex{},
		textures: map[string]TextureID{},
		fs:       fs,
	}
}

func (c *TextureCache) Load(filename string) (TextureID, error) {
	fmt.Println("texture cache: loading", filename)

	c.mutex.RLock()
	t, exists := c.textures[filename]
	c.mutex.RUnlock()

	if exists {
		return t, nil

	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	t, err := c.loadTexture(filename)
	if err != nil {
		return 0, err
	}

	c.textures[filename] = t
	return t, nil
}

func (c *TextureCache) loadTexture(file string) (TextureID, error) {
	imgFile, err := c.fs.OpenFile(file, 0, 0400)
	if err != nil {
		return 0, err
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)

	return TextureID(texture), nil
}
