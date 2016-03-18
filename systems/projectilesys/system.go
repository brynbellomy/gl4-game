package projectilesys

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/physicssys"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
)

type (
	System struct {
		entityManager  *entity.Manager
		componentQuery entity.ComponentMask

		positionCmptSet   entity.IComponentSet
		physicsCmptSet    entity.IComponentSet
		projectileCmptSet entity.IComponentSet
	}
)

// ensure that System conforms to entity.ISystem
var _ entity.ISystem = &System{}

func New() *System {
	return &System{}
}

func (s *System) Update(t common.Time) {
	matchIDs := s.entityManager.EntitiesMatching(s.componentQuery)
	positionCmptIdxs, err := s.positionCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}
	physicsCmptIdxs, err := s.physicsCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}
	projectileCmptIdxs, err := s.projectileCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}

	positionCmptSlice := s.positionCmptSet.Slice().(positionsys.ComponentSlice)
	physicsCmptSlice := s.physicsCmptSet.Slice().(physicssys.ComponentSlice)
	projectileCmptSlice := s.projectileCmptSet.Slice().(ComponentSlice)

	// check for collisions first
	for i := 0; i < len(projectileCmptIdxs); i++ {
		if len(physicsCmptSlice[physicsCmptIdxs[i]].GetCollisions()) > 0 {
			projectileCmptSlice[projectileCmptIdxs[i]].State = Impacting
		}
	}

	for i := 0; i < len(projectileCmptIdxs); i++ {
		projectileCmpt := projectileCmptSlice[projectileCmptIdxs[i]]
		physicsCmpt := physicsCmptSlice[physicsCmptIdxs[i]]
		positionCmpt := positionCmptSlice[positionCmptIdxs[i]]

		switch projectileCmpt.State {
		case Firing:
			headingNorm := projectileCmpt.GetHeading().Normalize()
			physicsCmpt.SetVelocity(headingNorm.Mul(projectileCmpt.ExitVelocity))
			physicsCmpt.AddForce(headingNorm.Mul(projectileCmpt.Thrust))

			// only stay in the Firing state for the first frame
			projectileCmpt.State = Flying

		case Flying:
			force := projectileCmpt.GetHeading().Normalize().Mul(projectileCmpt.Thrust)
			physicsCmpt.AddForce(force)

			v := physicsCmpt.GetVelocity()
			theta := float32(math.Atan2(float64(v.Y()), float64(v.X())))
			positionCmpt.SetRotation(theta)

		case Impacting:
			physicsCmpt.SetVelocity(mgl32.Vec2{0, 0})
			if projectileCmpt.RemoveOnContact {
				s.entityManager.RemoveEntity(matchIDs[i])
			}
		}

		projectileCmptSlice[projectileCmptIdxs[i]] = projectileCmpt
		physicsCmptSlice[physicsCmptIdxs[i]] = physicsCmpt
		positionCmptSlice[positionCmptIdxs[i]] = positionCmpt
	}
}

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"projectile": {
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

	componentQuery, err := s.entityManager.MakeCmptQuery([]string{"position", "physics", "projectile"})
	if err != nil {
		panic(err)
	}
	s.componentQuery = componentQuery

	projectileCmptSet, err := s.entityManager.GetComponentSet("projectile")
	if err != nil {
		panic(err)
	}
	s.projectileCmptSet = projectileCmptSet

	positionCmptSet, err := s.entityManager.GetComponentSet("position")
	if err != nil {
		panic(err)
	}
	s.positionCmptSet = positionCmptSet

	physicsCmptSet, err := s.entityManager.GetComponentSet("physics")
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
