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
)

func New() *System {
	return &System{
		conditionFactory: common.NewConfigFactory(),
		effectFactory:    common.NewConfigFactory(),
	}
}

func (s *System) RegisterConditionType(name string, specimen ICondition) {
	s.conditionFactory.Register(name, specimen)
}

func (s *System) RegisterEffectType(name string, specimen IEffect) {
	s.effectFactory.Register(name, specimen)
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

		// @@TODO: move this to some kind of component init hook
		if triggerCmpt.TriggersConfig != nil && triggerCmpt.Triggers == nil {
			ts := make([]Trigger, len(triggerCmpt.TriggersConfig))
			for i, cfg := range triggerCmpt.TriggersConfig {
				condition, err := s.conditionFactory.Build(cfg.Condition)
				if err != nil {
					// @@TODO
					panic(err)
				}

				effect, err := s.effectFactory.Build(cfg.Effect)
				if err != nil {
					// @@TODO
					panic(err)
				}

				ts[i] = Trigger{
					Condition: condition.(ICondition),
					Effect:    effect.(IEffect),
				}

				ts[i].Condition.WillJoinManager(s.entityManager)
				ts[i].Effect.WillJoinManager(s.entityManager)
			}
			triggerCmpt.Triggers = ts

			triggerCmptSlice[triggerCmptIdxs[i]] = triggerCmpt
		}

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

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"trigger": {Component{}, ComponentSlice{}},
	}
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
}
