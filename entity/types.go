package entity

import (
	"github.com/brynbellomy/gl4-game/common"
)

type (
	ISystem interface {
		Update(t common.Time)
		ComponentsWillJoin(eid ID, components []IComponent)
	}

	ID uint64

	IComponent interface{}
)
