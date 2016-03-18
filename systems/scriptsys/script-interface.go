package scriptsys

import (
"github.com/brynbellomy/gl4-game/entity"
)

type (
    ScriptInterface struct {
        em *ScriptEntityManager
    }
)

func (si *ScriptInterface) EntityManager() *ScriptEntityManager {
    return si.em
}

type (
    ScriptEntityManager struct {
        entityManager *entity.Manager
    }
)

func (em *ScriptEntityManager) MakeCmptQuery(cmptTypes []string) (uint64, error) {
    x, err := em.entityManager.MakeCmptQuery(cmptTypes)
    if err != nil {
        return 0, err
    }
    return uint64(x), nil
}

func (em *ScriptEntityManager) EntitiesMatching(cmptMask uint64) []entity.ID {
    return em.entityManager.EntitiesMatching(entity.ComponentMask(cmptMask))
}