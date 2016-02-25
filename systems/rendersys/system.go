package rendersys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/context"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/node"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
)

type (
	System struct {
		entities   []entityAspect
		renderRoot node.INode
	}

	IRenderContext interface {
		CurrentTransform() mgl32.Mat4
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

func (s *System) SetRenderRoot(root node.INode) {
	s.renderRoot = root
}

func (s *System) Update(c context.IContext) {
	for _, ent := range s.entities {
		rnode := ent.renderCmpt.renderNode

		rnode.SetPos(ent.positionCmpt.Pos())
		rnode.SetSize(ent.positionCmpt.Size())
		rnode.SetTexture(ent.renderCmpt.Texture())
	}

	s.render(s.renderRoot)
}

func (s *System) render(node node.INode) {
	// renderCtx := NewRenderContext()
	node.Render()
	for _, child := range node.Children() {
		s.render(child)
	}
}

func (s *System) EntityWillJoin(eid entity.ID, components []entity.IComponent) {
	// if we find a *positionsys.Component on the entity, we keep track of it
	var positionCmpt positionsys.IComponent
	var renderCmpt *Component

	for _, cmpt := range components {
		if rc, is := cmpt.(*Component); is {
			renderCmpt = rc
		} else if pc, is := cmpt.(*positionsys.Component); is {
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
