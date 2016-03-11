package rendersys

import (
	"sort"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
	"github.com/brynbellomy/gl4-game/systems/rendersys/shader"
)

type (
	System struct {
		entities   []entityAspect
		entityMap  map[entity.ID]*entityAspect
		projection mgl32.Mat4
		camera     mgl32.Mat4

		shaderProgramCache *shader.ProgramCache
		nodeFactory        *NodeFactory
	}

	entityAspect struct {
		id           entity.ID
		renderCmpt   *Component
		positionCmpt *positionsys.Component
	}
)

func New(shaderProgramCache *shader.ProgramCache) *System {
	nodeFactory := NewNodeFactory()
	nodeFactory.RegisterNodeType("sprite", &SpriteNodeFactory{shaderProgramCache})

	return &System{
		entities:           []entityAspect{},
		entityMap:          map[entity.ID]*entityAspect{},
		shaderProgramCache: shaderProgramCache,
		nodeFactory:        nodeFactory,
	}
}

func (s *System) SetProjection(p mgl32.Mat4) {
	s.projection = p
}

func (s *System) SetCamera(camera mgl32.Mat4) {
	s.camera = camera
}

func (s *System) Update(t common.Time) {
	renderCtx := RenderContext{
		Projection: s.projection,
		Camera:     s.camera,
	}

	for _, ent := range s.entities {
		rnode := ent.renderCmpt.renderNode
		rnode.SetPos(ent.positionCmpt.GetPos())
		rnode.SetSize(ent.positionCmpt.GetSize())
		rnode.SetRotation(ent.positionCmpt.GetRotation())
		rnode.SetTexture(ent.renderCmpt.Texture())
		rnode.Render(renderCtx)
	}
}

func (s *System) ComponentTypes() map[string]entity.IComponent {
	return map[string]entity.IComponent{
		"render": &Component{},
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	// em.RegisterComponentType("render", &Component{})
}

func (s *System) EntityComponentsChanged(eid entity.ID, components []entity.IComponent) {
	var positionCmpt *positionsys.Component
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

	if renderCmpt != nil && positionCmpt != nil {
		if _, exists := s.entityMap[eid]; !exists {
			// initialize the render node on the `renderCmpt`
			node, err := s.nodeFactory.NodeFromConfig(renderCmpt.NodeType, renderCmpt.NodeConfig)
			if err != nil {
				// @@TODO
				panic(err.Error())
			}

			renderCmpt.renderNode = node

			s.entities = append(s.entities, entityAspect{
				id:           eid,
				positionCmpt: positionCmpt,
				renderCmpt:   renderCmpt,
			})

			s.entityMap[eid] = &s.entities[len(s.entities)-1]
		}

	} else {
		if _, exists := s.entityMap[eid]; exists {
			idx := -1
			for i := range s.entities {
				if s.entities[i].id == eid {
					idx = i
					break
				}
			}

			if idx >= 0 {
				s.entities = append(s.entities[:idx], s.entities[idx+1:]...)
			}

			delete(s.entityMap, eid)
		}
	}

	// sort entities by z-index every time entity/component list changes
	sort.Sort(sortableEntities(s.entities))
}

type sortableEntities []entityAspect

func (s sortableEntities) Len() int {
	return len(s)
}

func (s sortableEntities) Less(i, j int) bool {
	z1 := s[i].positionCmpt.GetZIndex()
	z2 := s[j].positionCmpt.GetZIndex()
	if z1 < z2 {
		return true
	} else if z1 > z2 {
		return false
	} else {
		// defer to y coordinate in the scene if the z-indices are the same
		return s[i].positionCmpt.GetPos().Y() < s[j].positionCmpt.GetPos().Y()
	}
}

func (s sortableEntities) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
