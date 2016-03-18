package tagsys

import (
	"fmt"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	System struct {
		entityManager  *entity.Manager
		componentQuery entity.ComponentMask
		tagCmptSet     entity.IComponentSet

		idsByTag map[string]entity.ID
	}
)

// ensure that System conforms to entity.ISystem
var _ entity.ISystem = &System{}

func New() *System {
	return &System{
		idsByTag: map[string]entity.ID{},
	}
}

func (s *System) EntityWithTag(tag string) (entity.ID, bool) {
	sl := s.tagCmptSet.Slice().(ComponentSlice)
	for idx, x := range sl {
		if tag == x.GetTag() {
			return s.tagCmptSet.IDForIndex(idx)
		}
	}
	return entity.InvalidID, false
}

func (s *System) Update(t common.Time) {
	// no-op
}

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"tag": {
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

	componentQuery, err := s.entityManager.MakeCmptQuery([]string{"tag"})
	if err != nil {
		panic(err)
	}
	s.componentQuery = componentQuery

	tagCmptSet, err := s.entityManager.GetComponentSet("tag")
	if err != nil {
		panic(err)
	}
	s.tagCmptSet = tagCmptSet
}

func (s *System) ComponentsWillJoin(eid entity.ID, cmpts []entity.IComponent) error {
	// no-op
	return nil
}

func (s *System) ComponentsWillLeave(eid entity.ID, cmpts []entity.IComponent) error {
	for _, cmpt := range cmpts {
		if cmpt, is := cmpt.(Component); is {
			tag := cmpt.GetTag()

			if s.idsByTag[tag] == eid {
				delete(s.idsByTag, tag)
			} else {
				return fmt.Errorf("tagsys.System.ComponentsWillLeave: another entity claims this tag (eid: %v, other eid: %v, tag: %v)", eid, s.idsByTag[tag], tag)
			}
		}
	}
	return nil
}
