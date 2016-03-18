package positionsys

import (
	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
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
