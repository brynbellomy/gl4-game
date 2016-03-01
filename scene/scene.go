package scene

import (
	"path"

	"github.com/brynbellomy/gl4-game/entity"
)

type (
	Scene struct {
		assetRoot     string
		entityManager entity.Manager
	}

	Config struct {
		AssetRoot string
		Systems   []entity.ISystem
	}
)

func New(config Config) *Scene {
	return &Scene{
		assetRoot:     config.AssetRoot,
		entityManager: entity.NewManager(config.Systems),
	}
}

func (s *Scene) AssetRoot() string {
	return s.assetRoot
}

func (s *Scene) AssetPath(asset string) string {
	return path.Join(s.assetRoot, asset)
}

func (s *Scene) EntityManager() *entity.Manager {
	return &s.entityManager
}
