package scene

import (
	"github.com/brynbellomy/gl4-game/context"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/input"
)

type (
	Scene struct {
		inputHandler input.IHandler
		systems      []entity.ISystem
	}

	IScene interface {
		Update()
		AddEntity(eid entity.ID, components []entity.IComponent)
	}

	Config struct {
		InputHandler input.IHandler
		Systems      []entity.ISystem
	}
)

func New(config Config) IScene {
	return &Scene{
		inputHandler: config.InputHandler,
		systems:      config.Systems,
	}
}

func (s *Scene) AddEntity(eid entity.ID, components []entity.IComponent) {
	for _, sys := range s.systems {
		sys.EntityWillJoin(eid, components)
	}
}

func (s *Scene) Update() {
	c := context.New()

	if s.inputHandler != nil {
		s.inputHandler.Update()
	}

	for _, sys := range s.systems {
		sys.Update(c)
	}
}
