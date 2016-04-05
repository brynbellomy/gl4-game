package tilemapsys_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/brynbellomy/gl4-game/systems/assetsys"
	"github.com/brynbellomy/gl4-game/systems/tilemapsys"
)

func TestLoad(T *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		T.Error(err)
	}

	assetRoot := path.Join(cwd, "..", "..", "resources", "tilemaps")
	fs := assetsys.NewDefaultFilesystem(assetRoot)
	cache := tilemapsys.NewTilemapCache(fs)

	theMap, err := cache.Load("islands/monkey.tmx")
	if err != nil {
		T.Error(err)
	}

	// j, err := json.MarshalIndent(theMap.Layers[0].Tiles, "", "\t")
	// if err != nil {
	// 	T.Error(err)
	// }

	for _, layer := range theMap.Layers {
		for coord, tileGid := range layer.Tiles {
			fmt.Println(coord, "~>", tileGid)
			ts := theMap.FindTileset(tileGid)
			rect := theMap.TilesetRect(ts, ts.Image.Width, ts.Image.Height, false, tileGid)
			fmt.Println(ts.Image.Source)
            imageCoordsToTextureCoords(rect.)
			fmt.Println(rect)
		}
	}
}

func imageCoordsToTextureCoords(imgX, imgY, imgWidth, imgHeight float32) (float32, float32) {
	return imgX / imgWidth, imgY / imgHeight
}
