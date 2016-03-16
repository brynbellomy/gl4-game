package triggersys

import (
	"fmt"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	DebugEffect struct{}
)

func (e *DebugEffect) Execute(t common.Time, targets []entity.ID) error {
	fmt.Printf("DEBUG: targets = %v\n", targets)
	return nil
}

func (e *DebugEffect) WillJoinManager(em *entity.Manager) error {
	// no-op
	return nil
}

func (e *DebugEffect) WillLeaveManager() error {
	// no-op
	return nil
}
