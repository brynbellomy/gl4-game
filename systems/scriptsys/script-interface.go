package scriptsys

import (
	"github.com/yuin/gopher-lua"

	"github.com/brynbellomy/gl4-game/entity"
)

type (
	ScriptInterface struct {
		em *ScriptEntityManager
		L  *lua.LState
	}
)

func (si *ScriptInterface) EntityManager() *ScriptEntityManager {
	return si.em
}

type (
	ScriptEntityManager struct {
		entityManager *entity.Manager
		L             *lua.LState
	}
)

func (em *ScriptEntityManager) GetSystem(name string) entity.ISystem {
	return em.entityManager.GetSystem(name)
}

func (em *ScriptEntityManager) MakeCmptQuery(cmptTypes []string) (entity.ComponentMask, error) {
	return em.entityManager.MakeCmptQuery(cmptTypes)
}

func (em *ScriptEntityManager) EntitiesMatching(cmptMask entity.ComponentMask) []entity.ID {
	return em.entityManager.EntitiesMatching(entity.ComponentMask(int(cmptMask)))
}

func (em *ScriptEntityManager) GetComponentSet(name string) *ScriptComponentSet {
	cs, err := em.entityManager.GetComponentSet(name)
	if err != nil {
		// @@TODO?
		return nil
	}

	return &ScriptComponentSet{
		componentSet: cs,
		L:            em.L,
	}
}

type (
	ScriptComponentSet struct {
		componentSet entity.IComponentSet
		L            *lua.LState
	}
)

func (cs *ScriptComponentSet) Get(eid entity.ID) (entity.IComponent, error) {
	return cs.componentSet.Get(eid)
}

func (cs *ScriptComponentSet) Index(eid entity.ID) (int, bool) {
	i, exists := cs.componentSet.Index(eid)
	return i + 1, exists
}

func (cs *ScriptComponentSet) Indices(entityIDs []entity.ID) ([]int, error) {
	idxs, err := cs.componentSet.Indices(entityIDs)
	if err != nil {
		return nil, err
	}

	// add 1 to all indices to match Lua's 1-indexed scheme
	for i := range idxs {
		idxs[i]++
	}

	return idxs, nil
}

func (cs *ScriptComponentSet) Slice() interface{} {
	return cs.componentSet.Slice()
}
