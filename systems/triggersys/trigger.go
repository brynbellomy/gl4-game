package triggersys

import (
	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	// Trigger struct {
	// 	ICondition
	// 	IEffect
	// }

	ICondition interface {
		GetMatches(t common.Time) ([]entity.ID, error)
		WillJoinManager(em *entity.Manager) error
		WillLeaveManager() error
	}

	IEffect interface {
		Execute(t common.Time, targets []entity.ID) error
		WillJoinManager(em *entity.Manager) error
		WillLeaveManager() error
	}
)
