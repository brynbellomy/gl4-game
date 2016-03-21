package scriptsys

import (
	// "github.com/layeh/gopher-luar"
	"reflect"

	"github.com/brynbellomy/go-luaconv"
	"github.com/yuin/gopher-lua"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	System struct {
		entityManager  *entity.Manager
		componentQuery entity.ComponentMask
		scriptCmptSet  entity.IComponentSet

		scriptCache *ScriptCache
	}
)

// ensure that System conforms to entity.ISystem
var _ entity.ISystem = &System{}

func New(scriptCache *ScriptCache) *System {
	return &System{
		scriptCache: scriptCache,
	}
}

func (s *System) Name() string {
	return "script"
}

func (s *System) Update(t common.Time) {
	matchIDs := s.entityManager.EntitiesMatching(s.componentQuery)
	scriptCmptIdxs, err := s.scriptCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err) // @@TODO
	}

	scriptCmptSlice := s.scriptCmptSet.Slice().(ComponentSlice)

	for i := 0; i < len(scriptCmptIdxs); i++ {
		scriptCmpt := scriptCmptSlice[scriptCmptIdxs[i]]

		scriptInterface := &ScriptInterface{
			L: scriptCmpt.L,
			em: &ScriptEntityManager{
				L:             scriptCmpt.L,
				entityManager: s.entityManager,
			},
		}

		luaT, err := luaconv.Wrap(scriptCmpt.L, reflect.ValueOf(t))
		if err != nil {
			panic(err) // @@TODO
		}

		luaScriptCtx, err := luaconv.Wrap(scriptCmpt.L, reflect.ValueOf(scriptInterface))
		if err != nil {
			panic(err) // @@TODO
		}

		luaOwnID, err := luaconv.Wrap(scriptCmpt.L, reflect.ValueOf(matchIDs[i]))
		if err != nil {
			panic(err) // @@TODO
		}

		err = scriptCmpt.L.CallByParam(lua.P{
			Fn:      scriptCmpt.L.GetGlobal("update"),
			NRet:    0,
			Protect: false,
		}, luaT, luaScriptCtx, luaOwnID)

		if err != nil {
			panic(err) // @@TODO
		}
	}
}

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"script": {
			Coder: common.NewCoder(common.CoderConfig{
				ConfigType: ComponentCfg{},
				Tag:        "config",
				Decode: func(x interface{}) (interface{}, error) {
					return s.initScriptComponent(x.(ComponentCfg))
				},
				Encode: func(x interface{}) (interface{}, error) { /* @@TODO */ panic("unimplemented") },
			}),
			Slice: ComponentSlice{},
		},
	}
}

func (s *System) initScriptComponent(cfg ComponentCfg) (Component, error) {
	scriptSrc, err := s.scriptCache.Load(cfg.Filename)
	if err != nil {
		return Component{}, err
	}

	L := lua.NewState()

	err = L.DoString(scriptSrc)
	if err != nil {
		return Component{}, err
	}

	return Component{
		L: L,
	}, nil
}

func (s *System) WillJoinManager(em *entity.Manager) {
	s.entityManager = em

	componentQuery, err := s.entityManager.MakeCmptQuery([]string{"script"})
	if err != nil {
		panic(err)
	}
	s.componentQuery = componentQuery

	scriptCmptSet, err := s.entityManager.GetComponentSet("script")
	if err != nil {
		panic(err)
	}
	s.scriptCmptSet = scriptCmptSet
}

func (s *System) ComponentsWillJoin(eid entity.ID, cmpts []entity.IComponent) error {
	// no-op
	return nil
}

func (s *System) ComponentsWillLeave(eid entity.ID, cmpts []entity.IComponent) error {
	// no-op
	return nil
}
