package gameobjsys

import (
	"math"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/animationsys"
	"github.com/brynbellomy/gl4-game/systems/movesys"
)

type (
	System struct {
		entityManager    *entity.Manager
		componentQuery   entity.ComponentMask
		gameobjCmptSet   entity.IComponentSet
		animationCmptSet entity.IComponentSet
		moveCmptSet      entity.IComponentSet
	}
)

// ensure that System conforms to entity.ISystem
var _ entity.ISystem = &System{}

func New() *System {
	return &System{}
}

func (s *System) Name() string {
	return "gameobj"
}

func (s *System) Update(t common.Time) {
	matchIDs := s.entityManager.EntitiesMatching(s.componentQuery)
	gameobjCmptIdxs, err := s.gameobjCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}
	animationCmptIdxs, err := s.animationCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}
	moveCmptIdxs, err := s.moveCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}

	gameobjCmptSlice := s.gameobjCmptSet.Slice().(ComponentSlice)
	animationCmptSlice := s.animationCmptSet.Slice().(animationsys.ComponentSlice)
	moveCmptSlice := s.moveCmptSet.Slice().(movesys.ComponentSlice)

	for i := 0; i < len(gameobjCmptIdxs); i++ {
		gameobjCmpt := gameobjCmptSlice[gameobjCmptIdxs[i]]
		animationCmpt := animationCmptSlice[animationCmptIdxs[i]]
		moveCmpt := moveCmptSlice[moveCmptIdxs[i]]

		vel := moveCmpt.Vector()

		if vel.Len() > 0 {
			radians := math.Atan2(float64(vel.Y()), float64(vel.X()))

			gameobjCmpt.Direction = DirectionFromRadians(radians)
			animationCmpt.SetIsAnimating(true)

		} else {
			animationCmpt.SetIsAnimating(false)
		}

		if action, exists := gameobjCmpt.Animations[gameobjCmpt.Action]; exists {
			anim := action[gameobjCmpt.Direction]
			if anim != "" {
				animationCmpt.SetAnimation(anim)
			}
		}

		gameobjCmptSlice[gameobjCmptIdxs[i]] = gameobjCmpt
		animationCmptSlice[animationCmptIdxs[i]] = animationCmpt
		moveCmptSlice[moveCmptIdxs[i]] = moveCmpt
	}
}

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"gameobj": {
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

	componentQuery, err := s.entityManager.MakeCmptQuery([]string{"gameobj", "animation", "move"})
	if err != nil {
		panic(err)
	}
	s.componentQuery = componentQuery

	gameobjCmptSet, err := s.entityManager.GetComponentSet("gameobj")
	if err != nil {
		panic(err)
	}
	s.gameobjCmptSet = gameobjCmptSet

	animationCmptSet, err := s.entityManager.GetComponentSet("animation")
	if err != nil {
		panic(err)
	}
	s.animationCmptSet = animationCmptSet

	moveCmptSet, err := s.entityManager.GetComponentSet("move")
	if err != nil {
		panic(err)
	}
	s.moveCmptSet = moveCmptSet
}

func (s *System) ComponentsWillJoin(eid entity.ID, cmpts []entity.IComponent) error {
	// no-op
	return nil
}

func (s *System) ComponentsWillLeave(eid entity.ID, cmpts []entity.IComponent) error {
	// no-op
	return nil
}
