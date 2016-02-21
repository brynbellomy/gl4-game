package node

import (
	_ "image/png"
	"math"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
)

type (
	Node struct {
		size     *common.Size
		position mgl32.Vec2
		children []INode

		previousTime, totalTime float64
	}

	INode interface {
		Update()
		Children() []INode
		AddChild(child INode)
		Position() mgl32.Vec2
		Size() *common.Size
	}
)

type Config struct {
	Size     *common.Size
	Position mgl32.Vec2
}

func New(config Config) *Node {
	return &Node{
		size:     config.Size,
		position: config.Position,
		children: []INode{},
	}
}

func (n *Node) Position() mgl32.Vec2 {
	return n.position
}

func (n *Node) Size() *common.Size {
	return n.size
}

func (n *Node) Children() []INode {
	return n.children
}

func (n *Node) AddChild(child INode) {
	n.children = append(n.children, child)
}

func (n *Node) Update() {
	// Update
	{
		t := glfw.GetTime()
		elapsed := t - n.previousTime
		n.previousTime = t
		n.totalTime += elapsed
		n.position = mgl32.Vec2{
			float32(math.Sin(n.totalTime)),
			float32(math.Cos(n.totalTime)),
		}
	}

	for _, child := range n.Children() {
		child.Update()
	}
}
