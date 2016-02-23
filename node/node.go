package node

type (
	Node struct {
		children []INode

		// previousTime, totalTime float64
	}

	INode interface {
		Render()
		Children() []INode
		AddChild(child INode)
	}
)

func New() *Node {
	return &Node{
		children: []INode{},
	}
}

// func (n *Node) Position() mgl32.Vec2 {
// 	return n.position
// }

// func (n *Node) Size() *common.Size {
// 	return n.size
// }

func (n *Node) Children() []INode {
	return n.children
}

func (n *Node) AddChild(child INode) {
	n.children = append(n.children, child)
}

func (n *Node) Render() {
	// {
	// 	t := glfw.GetTime()
	// 	elapsed := t - n.previousTime
	// 	n.previousTime = t
	// 	n.totalTime += elapsed
	// 	n.position = mgl32.Vec2{
	// 		float32(math.Sin(n.totalTime)),
	// 		float32(math.Cos(n.totalTime)),
	// 	}
	// }

	for _, child := range n.Children() {
		child.Render()
	}
}
