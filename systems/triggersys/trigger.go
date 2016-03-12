package triggersys

import (
	"github.com/brynbellomy/gl4-game/common"
	// "github.com/brynbellomy/gl4-game/entity"
)

type (
	Trigger struct {
		ICondition
		IEffect
	}

	ICondition interface {
		ConditionValue(t common.Time) bool
	}

	IEffect interface {
		Execute(t common.Time)
	}
)

func (t *Trigger) Update(t common.Time) {
	// @@TODO
	for eid, cmpts := range entities {
	}
}
