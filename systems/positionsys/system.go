package positionsys

import (
	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/triggersys"
)

type (
	System struct {
		entityManager *entity.Manager
	}
)

// ensure that System conforms to entity.ISystem
var _ entity.ISystem = &System{}

func New() *System {
	return &System{}
}

func (s *System) Name() string {
	return "position"
}

func (s *System) Update(t common.Time) {
	// no-op
}

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"position": {
			Coder: common.NewCoder(common.CoderConfig{
				ConfigType: Component{},
				Tag:        "config",
				Decode:     func(x interface{}) (interface{}, error) { return x.(Component), nil },
				Encode:     func(x interface{}) (interface{}, error) { /* @@TODO */ panic("unimplemented") },
			}),
			Slice: ComponentSlice{},
		},
	}
}

func (s *System) TriggerTypes() triggersys.TriggerTypes {
	return triggersys.TriggerTypes{
		Conditions: map[string]triggersys.ConditionTypeCfg{
			"distance": {
				Coder: common.NewCoder(common.CoderConfig{
					ConfigType: &DistanceCondition{},
					Tag:        "config",
					Decode:     func(x interface{}) (interface{}, error) { return x.(*DistanceCondition), nil },
					Encode:     func(x interface{}) (interface{}, error) { return x.(*DistanceCondition), nil },
				}),
			},
		},
		Effects: map[string]triggersys.EffectTypeCfg{},
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	s.entityManager = em
}

func (s *System) ComponentsWillJoin(eid entity.ID, cmpts []entity.IComponent) error {
	// no-op
	return nil
}

func (s *System) ComponentsWillLeave(eid entity.ID, cmpts []entity.IComponent) error {
	// no-op
	return nil
}
