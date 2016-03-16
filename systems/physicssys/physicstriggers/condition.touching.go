package physicstriggers

import (
	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/physicssys"
)

type (
	TouchingCondition struct {
		SubjectID entity.ID `config:"subjectID"`

		entityManager  *entity.Manager      `config:"-"`
		physicsCmptSet entity.IComponentSet `config:"-"`
	}
)

func (c *TouchingCondition) WillJoinManager(em *entity.Manager) error {
	c.entityManager = em

	physicsCmptSet, err := em.GetComponentSet("physics")
	if err != nil {
		return err
	}

	c.physicsCmptSet = physicsCmptSet
	return nil
}

func (c *TouchingCondition) WillLeaveManager() error {
	c.entityManager = nil
	c.physicsCmptSet = nil
	return nil
}

func (c *TouchingCondition) GetMatches(t common.Time) ([]entity.ID, error) {
	pc, err := c.physicsCmptSet.Get(c.SubjectID)
	if err != nil {
		return nil, err
	}

	physicsCmpt := pc.(physicssys.Component)

	matchIDs := make([]entity.ID, 0)
	for _, coll := range physicsCmpt.GetCollisions() {
		if coll.EntityA == c.SubjectID {
			matchIDs = append(matchIDs, coll.EntityB)
		} else if coll.EntityB == c.SubjectID {
			matchIDs = append(matchIDs, coll.EntityA)
		} else {
			panic("TouchingCondition: collision doesn't contain .SubjectID")
		}
	}

	return matchIDs, nil
}
