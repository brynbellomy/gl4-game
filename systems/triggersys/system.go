package triggersys

import (
	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	System struct {
		entityManager  *entity.Manager
		componentQuery entity.ComponentMask
		triggerCmptSet entity.IComponentSet

		conditionFactory *common.ConfigFactory
		effectFactory    *common.ConfigFactory
	}

	ITriggerProviderSystem interface {
		TriggerTypes() TriggerTypes
	}

	TriggerTypes struct {
		Conditions map[string]ConditionTypeCfg
		Effects    map[string]EffectTypeCfg
	}

	ConditionTypeCfg struct {
		Coder *common.Coder
	}

	EffectTypeCfg struct {
		Coder *common.Coder
	}
)

// ensure that System conforms to entity.ISystem
var _ entity.ISystem = &System{}

func New() *System {
	return &System{
		conditionFactory: common.NewConfigFactory(),
		effectFactory:    common.NewConfigFactory(),
	}
}

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"trigger": {
			Coder: common.NewCoder(common.CoderConfig{
				ConfigType: ComponentCfg{},
				Tag:        "config",
				Decode: func(x interface{}) (interface{}, error) {
					return s.initTriggerCmpt(x.(ComponentCfg))
				},
				Encode: func(x interface{}) (interface{}, error) {
					// @@TODO
					panic("unimplemented")
				},
			}),
			Slice: ComponentSlice{},
		},
	}
}

func (s *System) TriggerTypes() TriggerTypes {
	return TriggerTypes{
		Conditions: map[string]ConditionTypeCfg{},
		Effects: map[string]EffectTypeCfg{
			"debug": {
				Coder: common.NewCoder(common.CoderConfig{
					ConfigType: &DebugEffect{},
					Tag:        "config",
					Decode:     func(x interface{}) (interface{}, error) { return x.(*DebugEffect), nil },
					Encode:     func(x interface{}) (interface{}, error) { return x.(*DebugEffect), nil },
				}),
			},
		},
	}
}

func (s *System) initTriggerCmpt(cfg ComponentCfg) (Component, error) {
	ts := make([]Trigger, len(cfg.TriggersConfig))

	for i, triggerCfg := range cfg.TriggersConfig {
		condition, err := s.conditionFactory.Decode(triggerCfg.Condition)
		if err != nil {
			// @@TODO
			panic(err)
		}

		effect, err := s.effectFactory.Decode(triggerCfg.Effect)
		if err != nil {
			// @@TODO
			panic(err)
		}

		ts[i] = Trigger{
			Condition: condition.(ICondition),
			Effect:    effect.(IEffect),
		}
	}

	return Component{Triggers: ts}, nil
}

func (s *System) WillJoinManager(em *entity.Manager) {
	s.entityManager = em

	componentQuery, err := s.entityManager.MakeCmptQuery([]string{"trigger"})
	if err != nil {
		panic(err)
	}
	s.componentQuery = componentQuery

	triggerCmptSet, err := s.entityManager.GetComponentSet("trigger")
	if err != nil {
		panic(err)
	}
	s.triggerCmptSet = triggerCmptSet

	for _, sys := range s.entityManager.Systems() {
		if sys, is := sys.(ITriggerProviderSystem); is {
			types := sys.TriggerTypes()

			for name, c := range types.Conditions {
				s.conditionFactory.Register(name, c.Coder)
			}

			for name, c := range types.Effects {
				s.effectFactory.Register(name, c.Coder)
			}
		}
	}
}

func (s *System) ComponentsWillJoin(eid entity.ID, cmpts []entity.IComponent) error {
	for _, cmpt := range cmpts {
		if cmpt, is := cmpt.(Component); is {
			for i := range cmpt.Triggers {
				err := cmpt.Triggers[i].Condition.WillJoinManager(s.entityManager, eid)
				if err != nil {
					return err
				}

				err = cmpt.Triggers[i].Effect.WillJoinManager(s.entityManager, eid)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *System) ComponentsWillLeave(eid entity.ID, cmpts []entity.IComponent) error {
	for _, cmpt := range cmpts {
		if cmpt, is := cmpt.(Component); is {
			for i := range cmpt.Triggers {
				err := cmpt.Triggers[i].Condition.WillLeaveManager()
				if err != nil {
					return err
				}

				err = cmpt.Triggers[i].Effect.WillLeaveManager()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *System) Update(t common.Time) {
	matchIDs := s.entityManager.EntitiesMatching(s.componentQuery)
	triggerCmptIdxs, err := s.triggerCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}

	triggerCmptSlice := s.triggerCmptSet.Slice().(ComponentSlice)

	for i := 0; i < len(triggerCmptIdxs); i++ {
		triggerCmpt := triggerCmptSlice[triggerCmptIdxs[i]]

		for j := range triggerCmpt.Triggers {
			matches, err := triggerCmpt.Triggers[j].Condition.GetMatches(t)
			if err != nil {
				// @@TODO
				panic(err)
			}
			err = triggerCmpt.Triggers[j].Effect.Execute(t, matches)
			if err != nil {
				// @@TODO
				panic(err)
			}
		}
	}
}
