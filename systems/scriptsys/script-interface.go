package scriptsys

import (
	"reflect"

	"github.com/brynbellomy/go-luaconv"
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

func (em *ScriptEntityManager) MakeCmptQuery(cmptTypes *lua.LTable) lua.LNumber {
	cmptTypeStrs, err := luaconv.LuaToNativeValue(em.L, cmptTypes, reflect.TypeOf([]string{}), "")
	if err != nil {
		// @@TODO
		panic(err)
	}

	x, err := em.entityManager.MakeCmptQuery(cmptTypeStrs.Interface().([]string))
	if err != nil {
		// @@TODO?
		panic(err)
	}
	return lua.LNumber(x)
}

func (em *ScriptEntityManager) EntitiesMatching(cmptMask lua.LNumber) *lua.LTable {
	matchIDs := em.entityManager.EntitiesMatching(entity.ComponentMask(int(cmptMask)))
	luaVal, err := luaconv.NativeValueToLua(em.L, reflect.ValueOf(matchIDs), "")
	if err != nil {
		// @@TODO?
		panic(err)
	}

	return luaVal.(*lua.LTable)
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

func (cs *ScriptComponentSet) Indices(entityIDs *lua.LTable) *lua.LTable {
	eids, err := luaconv.LuaToNativeValue(cs.L, entityIDs, reflect.TypeOf([]entity.ID{}), "")
	if err != nil {
		// @@TODO?
		panic(err)
	}

	idxs, err := cs.componentSet.Indices(eids.Interface().([]entity.ID))
	if err != nil {
		// @@TODO?
		panic(err)
	}

	luaIdxs, err := luaconv.NativeValueToLua(cs.L, reflect.ValueOf(idxs), "")
	if err != nil {
		// @@TODO?
		panic(err)
	}

	return luaIdxs.(*lua.LTable)
}

func (cs *ScriptComponentSet) Slice() interface{} {
	return cs.componentSet.Slice()
}
