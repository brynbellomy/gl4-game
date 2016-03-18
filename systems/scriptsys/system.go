package scriptsys

import (
    "github.com/robertkrimen/otto"

    "github.com/brynbellomy/gl4-game/common"
    "github.com/brynbellomy/gl4-game/entity"
)

type (
    System struct {
        entityManager  *entity.Manager
        componentQuery entity.ComponentMask
        scriptCmptSet     entity.IComponentSet

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

func (s *System) Update(t common.Time) {
    matchIDs := s.entityManager.EntitiesMatching(s.componentQuery)
    scriptCmptIdxs, err := s.scriptCmptSet.Indices(matchIDs)
    if err != nil {
        panic(err)
    }

    scriptCmptSlice := s.scriptCmptSet.Slice().(ComponentSlice)

    for i := 0; i < len(scriptCmptIdxs); i++ {
        scriptCmpt := scriptCmptSlice[scriptCmptIdxs[i]]

        updateFn, err := scriptCmpt.vm.Get("update")
        if err != nil {
            // @@TODO
            panic(err)
        } else if updateFn == otto.UndefinedValue() {
            // @@TODO
            panic("function 'update' is not defined")
        }

        scriptInterface := &ScriptInterface{
            em: &ScriptEntityManager{
                entityManager: s.entityManager,
            },
        }

        _, err = updateFn.Call(otto.UndefinedValue(), t, scriptInterface)
        if err != nil {
            // @@TODO
            panic(err)
        }
    }
}

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
    return map[string]entity.CmptTypeCfg{
        "script": {
            Coder: common.NewCoder(common.CoderConfig{
                ConfigType: ComponentCfg{},
                Tag:        "config",
                Decode:     func(x interface{}) (interface{}, error) {
                    return s.initScriptComponent(x.(ComponentCfg))
                },
                Encode:     func(x interface{}) (interface{}, error) { /* @@TODO */ panic("unimplemented") },
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

    vm := otto.New()

    script, err := vm.Compile(cfg.Filename, scriptSrc)
    if err != nil {
        return Component{}, err
    }

    _, err = vm.Run(script)
    if err != nil {
        return Component{}, err
    }

    return Component{
        vm: vm,
        script: script,
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
