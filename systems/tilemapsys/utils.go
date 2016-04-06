package tilemapsys

import (
	"github.com/azul3d-legacy/tmx"

	"github.com/brynbellomy/gl4-game/systems/rendersys/texture"
)

type (
	loadedVertexSublayer struct {
		vao         uint32
		texture     texture.TextureID
		numVertices int
	}

	vertexSublayer struct {
		vertices []float32
		texture  string
	}
)

func VerticesFromTilemap(tilemap *tmx.Map) []vertexSublayer {
	mapWidthPx := float32(tilemap.TileWidth * tilemap.Width)
	mapHeightPx := float32(tilemap.TileHeight * tilemap.Height)

	// initialize the render node's vertex data from the tilemap

	// tiles := make(sortableTiles, 0)
	// for i, layer := range tilemap.Layers {
	// 	for tileCoord, tileGid := range layer.Tiles {
	// 		tiles = append(tiles, sortableTile{layer: i, coord: tileCoord, gid: tileGid})
	// 	}
	// }

	// sort.Sort(tiles)

	allSublayers := make([]vertexSublayer, 0)

	for _, layer := range tilemap.Layers {
		sublayers := map[*tmx.Tileset]vertexSublayer{}

		for tileCoord, tileGid := range layer.Tiles {
			tileset := tilemap.FindTileset(tileGid)

			// calculate X, Y, Z coords
			var minX, maxX, minY, maxY float32
			{
				minX = (float32(mapWidthPx) - float32((tileCoord.X-1)*tilemap.TileWidth)) * (1.0 / float32(mapWidthPx))
				maxX = (float32(mapWidthPx) - float32((tileCoord.X)*tilemap.TileWidth)) * (1.0 / float32(mapWidthPx))
				minY = float32((tileCoord.Y-1)*tilemap.TileHeight) * (1.0 / float32(mapHeightPx))
				maxY = float32((tileCoord.Y)*tilemap.TileHeight) * (1.0 / float32(mapHeightPx))
			}

			// calculate texture coords (U, V)
			var minU, maxU, minV, maxV float32
			{
				imgWidth, imgHeight := tileset.Image.Width, tileset.Image.Height
				imgRect := tilemap.TilesetRect(tileset, imgWidth, imgHeight, false, tileGid)
				maxU = float32(imgRect.Min.X) / float32(imgWidth)
				minU = float32(imgRect.Max.X) / float32(imgWidth)
				minV = float32(imgRect.Min.Y) / float32(imgHeight)
				maxV = float32(imgRect.Max.Y) / float32(imgHeight)

				// minU, minV, maxU, maxV = 0, 0, 1, 1
			}

			vs := []float32{
				//  X, Y, Z, U, V

				// Front
				minX, minY, 0.0, maxU, minV,
				maxX, minY, 0.0, minU, minV,
				minX, maxY, 0.0, maxU, maxV,

				maxX, minY, 0.0, minU, minV,
				maxX, maxY, 0.0, minU, maxV,
				minX, maxY, 0.0, maxU, maxV,
			}

			sublayer := sublayers[tileset]
			sublayer.vertices = append(sublayer.vertices, vs...)
			sublayer.texture = tileset.Image.Source
			sublayers[tileset] = sublayer
		}

		for _, sublayer := range sublayers {
			allSublayers = append(allSublayers, sublayer)
		}
	}

	return allSublayers
}

type (
	sortableTile struct {
		layer int
		coord tmx.Coord
		gid   uint32
	}

	sortableTiles []sortableTile
)

func (st sortableTiles) Len() int {
	return len(st)
}

func (st sortableTiles) Swap(i, j int) {
	st[i], st[j] = st[j], st[i]
}

func (st sortableTiles) Less(i, j int) bool {
	if st[i].layer < st[j].layer {
		return true
	} else if st[i].layer > st[j].layer {
		return false
	}

	if st[i].coord.Y < st[j].coord.Y {
		return true
	} else if st[i].coord.Y > st[j].coord.Y {
		return false
	}

	return st[i].coord.X < st[j].coord.X
}
