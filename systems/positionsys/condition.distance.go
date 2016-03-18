package positionsys

import (
    "fmt"

    "github.com/brynbellomy/gl4-game/common"
    "github.com/brynbellomy/gl4-game/entity"
)

type (
    DistanceCondition struct {
        Distance float32               `config:"distance"`
        Mode     DistanceConditionMode `config:"mode"`

        id              entity.ID            `config:"-"`
        entityManager   *entity.Manager      `config:"-"`
        positionCmptSet entity.IComponentSet `config:"-"`
    }

    DistanceConditionMode string
)

const (
    DistanceConditionModeGreater DistanceConditionMode = "greater than"
    DistanceConditionModeLess    DistanceConditionMode = "less than"
)

func (c *DistanceCondition) WillJoinManager(em *entity.Manager, eid entity.ID) error {
    c.id = eid
    c.entityManager = em

    positionCmptSet, err := em.GetComponentSet("position")
    if err != nil {
        return err
    }

    c.positionCmptSet = positionCmptSet
    return nil
}

func (c *DistanceCondition) WillLeaveManager() error {
    c.id = entity.InvalidID
    c.entityManager = nil
    c.positionCmptSet = nil
    return nil
}

func (c *DistanceCondition) GetMatches(t common.Time) ([]entity.ID, error) {
    pc, err := c.positionCmptSet.Get(c.id)
    if err != nil {
        return nil, err
    }

    positionCmpt := pc.(Component)
    pos := positionCmpt.GetPos()

    matchIDs := make([]entity.ID, 0)
    for i, cmpt := range c.positionCmptSet.Slice().(ComponentSlice) {
        otherPos := cmpt.GetPos()

        dist := otherPos.Sub(pos).Len()
        if dist < c.Distance && c.Mode == DistanceConditionModeLess ||
            dist > c.Distance && c.Mode == DistanceConditionModeGreater {

            eid, ok := c.positionCmptSet.IDForIndex(i)
            if !ok {
                return nil, fmt.Errorf("positionsys.DistanceCondition.GetMatches: cannot get positionCmptSet.IDForIndex(%v)", i)
            }

            matchIDs = append(matchIDs, eid)
        }
    }

    return matchIDs, nil
}
