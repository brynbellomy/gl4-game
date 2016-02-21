package scene

import (
	"github.com/brynbellomy/gl-test/input"
	"github.com/brynbellomy/gl-test/node"
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
