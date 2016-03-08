package rendersys

import "errors"

type (
	NodeFactory struct {
		factories map[string]INodeFactory
	}

	INodeFactory interface {
		NodeFromConfig(config map[string]interface{}) (INode, error)
	}
)

func NewNodeFactory() *NodeFactory {
	return &NodeFactory{
		factories: map[string]INodeFactory{},
	}
}

func (f *NodeFactory) RegisterNodeType(name string, fac INodeFactory) {
	f.factories[name] = fac
}

func (f *NodeFactory) NodeFromConfig(nodeType string, cfg map[string]interface{}) (INode, error) {
	fac, exists := f.factories[nodeType]
	if !exists {
		return nil, errors.New("no such node factory '" + nodeType + "'")
	}

	return fac.NodeFromConfig(cfg)
}
