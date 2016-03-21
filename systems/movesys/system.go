package movesys

import (
	// "github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/physicssys"
)

type (
	System struct {
		entityManager  *entity.Manager
		componentQuery entity.ComponentMask

		moveCmptSet    entity.IComponentSet
		physicsCmptSet entity.IComponentSet
	}
)

// ensure that System conforms to entity.ISystem
var _ entity.ISystem = &System{}

func New() *System {
	return &System{}
}

func (s *System) Name() string {
	return "move"
}

func (s *System) Update(t common.Time) {
	matchIDs := s.entityManager.EntitiesMatching(s.componentQuery)

	moveCmptIdxs, err := s.moveCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}
	physCmptIdxs, err := s.physicsCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}

	moveCmptSlice := s.moveCmptSet.Slice().(ComponentSlice)
	physicsCmptSlice := s.physicsCmptSet.Slice().(physicssys.ComponentSlice)

	for i := 0; i < len(moveCmptIdxs); i++ {
		vec := moveCmptSlice[moveCmptIdxs[i]].Vector()
		physicsCmptSlice[physCmptIdxs[i]].SetInstantaneousVelocity(vec)
	}
}

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"move": {
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

	componentQuery, err := em.MakeCmptQuery([]string{"move", "physics"})
	if err != nil {
		panic(err)
	}
	s.componentQuery = componentQuery

	moveCmptSet, err := em.GetComponentSet("move")
	if err != nil {
		panic(err)
	}
	s.moveCmptSet = moveCmptSet

	physicsCmptSet, err := em.GetComponentSet("physics")
	if err != nil {
		panic(err)
	}
	s.physicsCmptSet = physicsCmptSet
}

func (s *System) ComponentsWillJoin(eid entity.ID, cmpts []entity.IComponent) error {
	// no-op
	return nil
}

func (s *System) ComponentsWillLeave(eid entity.ID, cmpts []entity.IComponent) error {
	// no-op
	return nil
}
