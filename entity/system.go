package entity

import (
	"github.com/brynbellomy/gl4-game/context"
)

type (
	ISystem interface {
		Update(c context.IContext)
		EntityWillJoin(eid ID, components []IComponent)
	}

	ID uint64

	IComponent interface{}
)
