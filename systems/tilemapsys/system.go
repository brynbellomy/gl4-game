package tilemapsys

import (
	"sort"

	"github.com/azul3d-legacy/tmx"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/rendersys"
)

type (
	System struct {
		tilemapCache *TilemapCache
	}
)

// ensure that System conforms to entity.ISystem
var _ entity.ISystem = &System{}

func New(tilemapCache *TilemapCache) *System {
	return &System{tilemapCache}
}

func (s *System) Name() string {
	return "tilemap"
}

func (s *System) Update(t common.Time) {
}

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"tilemap": {
			Slice: ComponentSlice{},
			Coder: common.NewCoder(common.CoderConfig{
				ConfigType: ComponentCfg{},
				Tag:        "config",
				Decode: func(x interface{}) (interface{}, error) {
					cfg := x.(ComponentCfg)

					tilemap, err := s.tilemapCache.Load(cfg.Tilemap)
					if err != nil {
						return nil, err
					}

					return Component{Tilemap: tilemap}, nil
				},
				Encode: func(x interface{}) (interface{}, error) { /* @@TODO */ panic("unimplemented") },
			}),
		},
	}
}

// func (s *System) RenderNodeFactories() map[string]rendersys.INodeFactory {
// 	return map[string]rendersys.INodeFactory{
// 		"tilemap": &TilemapNodeFactory{
// 			rendersys.SpriteNodeFactory{s.shaderProgramCache},
// 		},
// 	}
// }

func (s *System) WillJoinManager(em *entity.Manager) {
	// no-op
}

func (s *System) ComponentsWillJoin(eid entity.ID, cmpts []entity.IComponent) error {
	var renderCmptIdx = -1
	var tilemapCmptIdx = -1
	for i := range cmpts {
		if _, is := cmpts[i].(rendersys.Component); is {
			renderCmptIdx = i
		} else if _, is := cmpts[i].(Component); is {
			tilemapCmptIdx = i
		}

		if renderCmptIdx != -1 && tilemapCmptIdx != -1 {
			break
		}
	}

	if renderCmptIdx != -1 && tilemapCmptIdx != -1 {
		renderCmpt := cmpts[renderCmptIdx].(rendersys.Component)
		tilemapCmpt := cmpts[tilemapCmptIdx].(Component)

		tilemap := tilemapCmpt.Tilemap
		mapWidthPx := float32(tilemap.TileWidth * tilemap.Width)
		mapHeightPx := float32(tilemap.TileHeight * tilemap.Height)

		// initialize the render node's vertex data from the tilemap

		tiles := make(sortableTiles, 0)
		for i, layer := range tilemap.Layers {
			for tileCoord, tileGid := range layer.Tiles {
				tiles = append(tiles, sortableTile{layer: i, coord: tileCoord, gid: tileGid})
			}
		}

		sort.Sort(tiles)

		vertices := []float32{}
		for _, tile := range tiles {
			tileset := tilemap.FindTileset(tile.gid)

			// calculate X, Y, Z coords
			var minX, maxX, minY, maxY float32
			{
				minX = (float32(mapWidthPx) - float32((tile.coord.X-1)*tilemap.TileWidth)) * (1.0 / float32(mapWidthPx))
				maxX = (float32(mapWidthPx) - float32((tile.coord.X)*tilemap.TileWidth)) * (1.0 / float32(mapWidthPx))
				minY = float32((tile.coord.Y-1)*tilemap.TileHeight) * (1.0 / float32(mapHeightPx))
				maxY = float32((tile.coord.Y)*tilemap.TileHeight) * (1.0 / float32(mapHeightPx))
			}

			// calculate texture coords (U, V)
			var minU, maxU, minV, maxV float32
			{
				imgWidth, imgHeight := tileset.Image.Width, tileset.Image.Height
				imgRect := tilemap.TilesetRect(tileset, imgWidth, imgHeight, false, tile.gid)
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

			vertices = append(vertices, vs...)
		}

		renderCmpt.RenderNode().(*rendersys.SpriteNode).SetVertices(vertices)

		cmpts[renderCmptIdx] = renderCmpt
		cmpts[tilemapCmptIdx] = tilemapCmpt
	}

	return nil
}

func (s *System) ComponentsWillLeave(eid entity.ID, cmpts []entity.IComponent) error {
	// no-op
	return nil
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
