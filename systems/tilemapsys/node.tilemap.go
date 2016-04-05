package tilemapsys

// import (
// 	"errors"

// 	"github.com/brynbellomy/gl4-game/systems/rendersys"
// )

// type (
// 	TilemapNodeFactory struct {
// 		spriteNodeFactory rendersys.SpriteNodeFactory
// 	}
// )

// func (f *TilemapNodeFactory) NodeFromConfig(config map[string]interface{}) (INode, error) {
// 	tilemap, exists := config["tilemap"].(string)
// 	if !exists {
// 		return nil, errors.New("missing required key 'tilemap' (or wrong type)")
// 	}

// 	program, err := f.shaderProgramCache.Load(vertex, fragment)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return NewSpriteNode(program)
// }

// type (
// 	TilemapNode struct {
//         rendersys.SpriteNode
// 	}
// )
