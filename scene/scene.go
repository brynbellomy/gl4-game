package scene

import (
	"github.com/brynbellomy/gl4-game/input"
	"github.com/brynbellomy/gl4-game/node"
)

type (
	Scene struct {
		*node.Node
		inputHandler input.IHandler
	}

	IScene interface {
		SetInputHandler(inputHandler input.IHandler)
	}
)

func New(config node.Config) *Scene {
	return &Scene{
		Node: node.New(config),
	}
}

func (s *Scene) SetInputHandler(inputHandler input.IHandler) {
	s.inputHandler = inputHandler
}
