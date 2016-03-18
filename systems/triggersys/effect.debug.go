package triggersys

import (
	"fmt"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	DebugEffect struct {
		eid entity.ID
	}
)

func (e *DebugEffect) Execute(t common.Time, targets []entity.ID) error {
	if len(targets) > 0 {
		fmt.Printf("DEBUG: entity = %v, targets = %v\n", e.eid, targets)
	}
	return nil
}

func (e *DebugEffect) WillJoinManager(em *entity.Manager, eid entity.ID) error {
	e.eid = eid
	return nil
}

func (e *DebugEffect) WillLeaveManager() error {
	e.eid = entity.InvalidID
	return nil
}
