package scene

import (
	"path"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/context"
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	Scene struct {
		assetRoot  string
		projection mgl32.Mat4
		systems    []entity.ISystem
		entities   map[entity.ID][]entity.IComponent
	}

	IScene interface {
		AssetRoot() string
		AssetPath(asset string) string
		Projection() mgl32.Mat4
		Update()
		AddEntity(eid entity.ID, components []entity.IComponent)
	}

	Config struct {
		AssetRoot  string
		Projection mgl32.Mat4
		Systems    []entity.ISystem
	}
)

func New(config Config) *Scene {
	return &Scene{
		assetRoot:  config.AssetRoot,
		projection: config.Projection,
		systems:    config.Systems,
		entities:   make(map[entity.ID][]entity.IComponent),
	}
}

func (s *Scene) AssetRoot() string {
	return s.assetRoot
}

func (s *Scene) AssetPath(asset string) string {
	return path.Join(s.assetRoot, asset)
}

func (s *Scene) GetEntity(eid entity.ID) []entity.IComponent {
	return s.entities[eid]
}

func (s *Scene) AddEntity(eid entity.ID, components []entity.IComponent) {
	for _, sys := range s.systems {
		sys.EntityWillJoin(eid, components)
	}
	s.entities[eid] = components
}

func (s *Scene) Update() {
	c := context.New()

	for _, sys := range s.systems {
		sys.Update(c)
	}
}

func (s *Scene) Projection() mgl32.Mat4 {
	return s.projection
}
