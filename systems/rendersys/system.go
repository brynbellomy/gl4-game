package rendersys

import (
	"sort"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
)

type (
	System struct {
		entities   []entityAspect
		projection mgl32.Mat4
		camera     mgl32.Mat4
	}

	entityAspect struct {
		renderCmpt   *Component
		positionCmpt *positionsys.Component
	}
)

func New() *System {
	return &System{
		entities: []entityAspect{},
	}
}

func (s *System) SetProjection(p mgl32.Mat4) {
	s.projection = p
}

// func (s *System) SetCameraPos(pos mgl32.Vec2) {
// 	s.camera = mgl32.LookAtV(mgl32.Vec3{pos.X(), pos.Y(), 3}, mgl32.Vec3{pos.X(), pos.Y(), 0}, mgl32.Vec3{0, -1, 0})
// }

func (s *System) SetCamera(camera mgl32.Mat4) {
	s.camera = camera
}

func (s *System) Update(t common.Time) {
	renderCtx := RenderContext{
		Projection: s.projection,
		Camera:     s.camera,
	}

	// sort entities by z-index every time entity/component list changes
	sort.Sort(sortableEntities(s.entities))

	for _, ent := range s.entities {
		rnode := ent.renderCmpt.renderNode

		rnode.SetPos(ent.positionCmpt.Pos())
		rnode.SetSize(ent.positionCmpt.Size())
		rnode.SetTexture(ent.renderCmpt.Texture())

		rnode.Render(renderCtx)
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	// no-op
}

func (s *System) ComponentsWillJoin(eid entity.ID, components []entity.IComponent) {
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

type sortableEntities []entityAspect

func (s sortableEntities) Len() int {
	return len(s)
}

func (s sortableEntities) Less(i, j int) bool {
	z1 := s[i].positionCmpt.ZIndex()
	z2 := s[j].positionCmpt.ZIndex()
	if z1 < z2 {
		return true
	} else if z1 > z2 {
		return false
	} else {
		// defer to y coordinate in the scene if the z-indices are the same
		return s[i].positionCmpt.Pos().Y() < s[j].positionCmpt.Pos().Y()
	}
}

func (s sortableEntities) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
