package rendersys

import (
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	Component struct {
		renderNode INode
		texture    uint32
	}

	ComponentCfg struct {
		NodeType   string                 `config:"nodeType"`
		NodeConfig map[string]interface{} `config:"nodeConfig"`
	}

	ComponentSlice []Component
)

func (c *Component) Texture() uint32 {
	return c.texture
}

func (c *Component) SetTexture(tex uint32) {
	c.texture = tex
}

func (c Component) Clone() entity.IComponent {
	return c
}

func (cs ComponentSlice) Append(cmpt entity.IComponent) entity.IComponentSlice {
	return append(cs, cmpt.(Component))
}

func (cs ComponentSlice) Remove(idx int) entity.IComponentSlice {
	return append(cs[:idx], cs[idx+1:]...)
}

func (cs ComponentSlice) Get(idx int) (entity.IComponent, bool) {
	if idx >= len(cs) {
		return nil, false
	}
	return cs[idx], true
}

func (cs ComponentSlice) Set(idx int, cmpt entity.IComponent) bool {
	if idx >= len(cs) {
		return false
	}
	cs[idx] = cmpt.(Component)
	return true
}
