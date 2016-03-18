package physicssys

import (
	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	TouchingCondition struct {
		id             entity.ID            `config:"-"`
		entityManager  *entity.Manager      `config:"-"`
		physicsCmptSet entity.IComponentSet `config:"-"`
	}
)

func (c *TouchingCondition) WillJoinManager(em *entity.Manager, eid entity.ID) error {
	c.id = eid
	c.entityManager = em

	physicsCmptSet, err := em.GetComponentSet("physics")
	if err != nil {
		return err
	}

	c.physicsCmptSet = physicsCmptSet
	return nil
}

func (c *TouchingCondition) WillLeaveManager() error {
	c.id = entity.InvalidID
	c.entityManager = nil
	c.physicsCmptSet = nil
	return nil
}

func (c *TouchingCondition) GetMatches(t common.Time) ([]entity.ID, error) {
	pc, err := c.physicsCmptSet.Get(c.id)
	if err != nil {
		return nil, err
	}

	physicsCmpt := pc.(Component)

	matchIDs := make([]entity.ID, 0)
	for _, coll := range physicsCmpt.GetCollisions() {
		if coll.EntityA == c.id {
			matchIDs = append(matchIDs, coll.EntityB)
		} else if coll.EntityB == c.id {
			matchIDs = append(matchIDs, coll.EntityA)
		} else {
			panic("TouchingCondition: collision doesn't contain .id")
		}
	}

	return matchIDs, nil
}
