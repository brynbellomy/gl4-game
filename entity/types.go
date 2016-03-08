package entity

import (
	"github.com/brynbellomy/gl4-game/common"
)

type (
	ISystem interface {
		WillJoinManager(em *Manager)
		ComponentsWillJoin(eid ID, components []IComponent)
		ComponentsWillLeave(eid ID, components []IComponent)
		Update(t common.Time)
	}

	ID int64

	IComponent interface{}
)

const InvalidID ID = -1
