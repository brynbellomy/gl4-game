package entity

import (
	"github.com/brynbellomy/gl4-game/common"
)

type (
	ISystem interface {
		WillJoinManager(em *Manager)
		ComponentTypes() map[string]CmptTypeCfg
		Update(t common.Time)
	}

	CmptTypeCfg struct {
		Cmpt      IComponent
		CmptSlice IComponentSlice
	}

	ID int64

	ComponentMask uint64
	ComponentKind uint64

	IComponent interface {
		Clone() IComponent
	}

	Entity struct {
		ID             ID
		ComponentMask  ComponentMask
		Components     []IComponent
		ComponentKinds []ComponentKind
	}
)

const InvalidID ID = -1

func (c ComponentMask) Has(n ComponentKind) bool {
	return (c & ComponentMask(n)) > 0
}

func (c ComponentMask) HasAll(n ComponentMask) bool {
	return (c & n) == n
}

func (c ComponentMask) Add(other ComponentKind) ComponentMask {
	return c | ComponentMask(other)
}
