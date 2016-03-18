package rendersys

import (
	"fmt"
	"sort"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
	"github.com/brynbellomy/gl4-game/systems/rendersys/shader"
)

type (
	System struct {
		entityManager   *entity.Manager
		componentQuery  entity.ComponentMask
		renderCmptSet   entity.IComponentSet
		positionCmptSet entity.IComponentSet

		projection mgl32.Mat4
		camera     mgl32.Mat4

		shaderProgramCache *shader.ProgramCache
		nodeFactory        *NodeFactory
	}

	entityAspect struct {
		id           entity.ID
		renderCmpt   Component
		positionCmpt positionsys.Component
	}
)

func New(shaderProgramCache *shader.ProgramCache) *System {
	nodeFactory := NewNodeFactory()
	nodeFactory.RegisterNodeType("sprite", &SpriteNodeFactory{shaderProgramCache})

	return &System{
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

	matchIDs := s.entityManager.EntitiesMatching(s.componentQuery)
	renderCmptIdxs, err := s.renderCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}
	positionCmptIdxs, err := s.positionCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}

	renderCmptSlice := s.renderCmptSet.Slice().(ComponentSlice)
	positionCmptSlice := s.positionCmptSet.Slice().(positionsys.ComponentSlice)

	aspects := make([]entityAspect, len(renderCmptIdxs))
	for i := 0; i < len(aspects); i++ {
		aspects[i].id = matchIDs[i]
		aspects[i].renderCmpt = renderCmptSlice[renderCmptIdxs[i]]
		aspects[i].positionCmpt = positionCmptSlice[positionCmptIdxs[i]]
	}

	sort.Sort(sortableEntities(aspects))

	for _, ent := range aspects {
		rnode := ent.renderCmpt.renderNode
		rnode.SetPos(ent.positionCmpt.GetPos())
		rnode.SetSize(ent.positionCmpt.GetSize())
		rnode.SetRotation(ent.positionCmpt.GetRotation())
		rnode.SetTexture(ent.renderCmpt.Texture())
		rnode.Render(renderCtx)
	}
}

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"render": {
			Coder: common.NewCoder(common.CoderConfig{
				ConfigType: Component{},
				Tag:        "config",
				Decode:     func(x interface{}) (interface{}, error) { return x.(Component), nil },
				Encode:     func(x interface{}) (interface{}, error) { /* @@TODO */ panic("unimplemented") },
			}),
			Slice: ComponentSlice{},
		},
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	s.entityManager = em

	componentQuery, err := s.entityManager.MakeCmptQuery([]string{"position", "render"})
	if err != nil {
		panic(err)
	}
	s.componentQuery = componentQuery

	renderCmptSet, err := s.entityManager.GetComponentSet("render")
	if err != nil {
		panic(err)
	}
	s.renderCmptSet = renderCmptSet

	positionCmptSet, err := s.entityManager.GetComponentSet("position")
	if err != nil {
		panic(err)
	}
	s.positionCmptSet = positionCmptSet
}

func (s *System) ComponentsWillJoin(eid entity.ID, cmpts []entity.IComponent) error {
	for i := range cmpts {
		if cmpt, is := cmpts[i].(Component); is {
			node, err := s.nodeFactory.NodeFromConfig(cmpt.NodeType, cmpt.NodeConfig)
			if err != nil {
				return err
			}

			cmpt.renderNode = node
			cmpts[i] = cmpt
		}
	}

	return nil
}

func (s *System) ComponentsWillLeave(eid entity.ID, cmpts []entity.IComponent) error {
	for i := range cmpts {
		if cmpt, is := cmpts[i].(Component); is {
			err := cmpt.renderNode.Destroy()
			if err != nil {
				return err
			}
		}
	}

	return nil
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
