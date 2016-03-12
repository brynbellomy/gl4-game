package entity

import (
	"github.com/brynbellomy/gl4-game/common"
)

type (
	ISystem interface {
		WillJoinManager(em *Manager)
		ComponentTypes() map[string]IComponent
		EntityComponentsChanged(eid ID, components []IComponent)
		Update(t common.Time)
	}

	ID int64

	ComponentKind uint64

	IComponent interface {
		Kind() ComponentKind
		SetKind(k ComponentKind)
		Clone() IComponent
	}
)

const InvalidID ID = -1

func (k ComponentKind) Kind() ComponentKind {
	return k
}

func (k *ComponentKind) SetKind(newKind ComponentKind) {
	*k = newKind
}

func (k ComponentKind) KindIndex() int {
	return int(k)
}

func (k ComponentKind) KindMask() ComponentMask {
	return ComponentMask(1 << uint64(k))
}
