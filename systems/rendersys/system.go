package rendersys

import (
	"github.com/brynbellomy/gl4-game/context"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/node"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
)

type (
	System struct {
		entities []entityAspect
	}

	entityAspect struct {
		renderCmpt   *Component
		positionCmpt positionsys.IComponent
	}
)

func New() *System {
	return &System{
		entities: []entityAspect{},
	}
}

func (s *System) Update(c context.IContext) {
	for _, ent := range s.entities {
		node := ent.renderCmpt.renderNode

		node.SetPos(ent.positionCmpt.Pos())
		// node.SetSize()

		node.Render()
	}
}

func (s *System) EntityWillJoin(eid entity.ID, components []entity.IComponent) {
	// if we find a *positionsys.Component on the entity, we keep track of it
	var positionCmpt positionsys.IComponent
	var renderCmpt *Component

	for _, cmpt := range components {
		if rc, is := cmpt.(*Component); is {
			renderCmpt = rc
		} else if pc, is := cmpt.(positionsys.IComponent); is {
			positionCmpt = pc
		}

		if positionCmpt != nil && renderCmpt != nil {
			break
		}
	}

	if renderCmpt != nil {
		if positionCmpt == nil {
			panic("render component requires position component")
		}

		aspect := entityAspect{
			positionCmpt: positionCmpt,
			renderCmpt:   renderCmpt,
		}

		s.entities = append(s.entities, aspect)
	}
}

type (
	Component struct {
		renderNode *node.SpriteNode
	}

	IComponent interface {
		RenderNode() node.IRenderNode
	}
)

func NewComponent(renderNode *node.SpriteNode) IComponent {
	return &Component{renderNode}
}

func (c *Component) RenderNode() node.IRenderNode {
	return c.renderNode
}
