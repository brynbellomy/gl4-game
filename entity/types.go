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

// func (k ComponentKind) Kind() ComponentKind {
// 	return k
// }

// func (k *ComponentKind) SetKind(newKind ComponentKind) {
// 	*k = newKind
// }

// func (k ComponentKind) KindIndex() int {
// 	return int(k)
// }

// func (k ComponentKind) KindMask() ComponentMask {
// 	return ComponentMask(1 << uint64(k))
// }
