package physicssys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
	"github.com/brynbellomy/gl4-game/systems/triggersys"
)

type (
	System struct {
		entityManager  *entity.Manager
		componentQuery entity.ComponentMask

		previousTime common.Time
		onCollision  func(c Collision)

		positionCmptSet entity.IComponentSet
		physicsCmptSet  entity.IComponentSet
	}
)

// ensure that System conforms to entity.ISystem
var _ entity.ISystem = &System{}

func New() *System {
	return &System{}
}

func (s *System) OnCollision(fn func(c Collision)) {
	s.onCollision = fn
}

func (s *System) Update(t common.Time) {
	if s.previousTime == 0 {
		s.previousTime = t
		return
	}

	elapsed := t - s.previousTime

	matchIDs := s.entityManager.EntitiesMatching(s.componentQuery)
	positionCmptIdxs, err := s.positionCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}
	physCmptIdxs, err := s.physicsCmptSet.Indices(matchIDs)
	if err != nil {
		panic(err)
	}

	positionCmptSlice := s.positionCmptSet.Slice().(positionsys.ComponentSlice)
	physCmptSlice := s.physicsCmptSet.Slice().(ComponentSlice)

	//
	// apply acceleration / velocity
	//
	for i := 0; i < len(physCmptIdxs); i++ {
		physCmpt := physCmptSlice[physCmptIdxs[i]]
		posCmpt := positionCmptSlice[positionCmptIdxs[i]]

		accel := physCmpt.CurrentForces()
		physCmpt.ResetForces()

		vdelta := accel.Mul(float32(elapsed.Seconds()))

		newvel := physCmpt.GetVelocity().Add(vdelta)

		// friction
		// newvel = newvel.Mul(0.95)

		mag := newvel.Len()
		maxvel := physCmpt.GetMaxVelocity()
		if mag > 0 && maxvel < mag {
			newvel = newvel.Normalize().Mul(maxvel)
		}

		physCmpt.SetVelocity(newvel)

		// add the instantaneous velocity for the movement system
		newvel = newvel.Add(physCmpt.GetInstantaneousVelocity())

		newpos := posCmpt.GetPos().Add(newvel.Mul(float32(elapsed.Seconds())))
		posCmpt.SetPos(newpos)

		// take this opportunity to (unrelatedly) clear the collisions slice before step 2 (viz., check for collisions)
		physCmpt.ResetCollisions()

		physCmptSlice[physCmptIdxs[i]] = physCmpt
		positionCmptSlice[positionCmptIdxs[i]] = posCmpt
	}

	//
	// check for collisions
	//
	for i := 0; i < len(physCmptIdxs); i++ {
		physCmptA := physCmptSlice[physCmptIdxs[i]]
		posCmptA := positionCmptSlice[positionCmptIdxs[i]]

		for j := i + 1; j < len(physCmptIdxs); j++ {
			physCmptB := physCmptSlice[physCmptIdxs[j]]
			posCmptB := positionCmptSlice[positionCmptIdxs[j]]

			did := s.checkCollision(physCmptA, physCmptB, posCmptA, posCmptB)
			if did {
				c := Collision{matchIDs[i], matchIDs[j]}
				physCmptA.AddCollision(c)
				physCmptB.AddCollision(c)
				s.onCollision(c)
			}

			physCmptSlice[physCmptIdxs[j]] = physCmptB
			positionCmptSlice[positionCmptIdxs[j]] = posCmptB
		}
		physCmptSlice[physCmptIdxs[i]] = physCmptA
		positionCmptSlice[positionCmptIdxs[i]] = posCmptA
	}

	s.previousTime = t
}

func (s *System) checkCollision(physCmptA, physCmptB Component, posCmptA, posCmptB positionsys.Component) bool {
	var minA, maxA, minB, maxB float32

	if physCmptA.CollisionMask&physCmptB.CollidesWith == 0 && physCmptB.CollisionMask&physCmptA.CollidesWith == 0 {
		return false
	}

	for i := 0; i < len(physCmptA.GetBoundingBox())-1; i++ {
		normal := getNormal(physCmptA.GetBoundingBox()[i+1], physCmptA.GetBoundingBox()[i])
		minA, maxA = getMinMaxProjectedPoints(physCmptA.GetBoundingBox(), posCmptA.GetPos(), normal)
		minB, maxB = getMinMaxProjectedPoints(physCmptB.GetBoundingBox(), posCmptB.GetPos(), normal)

		if maxB < minA || maxA < minB {
			// no collision between these shapes
			return false
		}
	}

	for i := 0; i < len(physCmptB.GetBoundingBox())-1; i++ {
		normal := getNormal(physCmptB.GetBoundingBox()[i+1], physCmptB.GetBoundingBox()[i])
		minA, maxA = getMinMaxProjectedPoints(physCmptA.GetBoundingBox(), posCmptA.GetPos(), normal)
		minB, maxB = getMinMaxProjectedPoints(physCmptB.GetBoundingBox(), posCmptB.GetPos(), normal)

		if maxB < minA || maxA < minB {
			// no collision between these shapes
			return false
		}
	}

	return true
}

func getNormal(a, b mgl32.Vec2) mgl32.Vec2 {
	return mgl32.Vec2{
		-(a.Y() - b.Y()),
		a.X() - b.X(),
	}
}

func getMinMaxProjectedPoints(boundingBox BoundingBox, pos mgl32.Vec2, normal mgl32.Vec2) (float32, float32) {
	min := boundingBox[0].Add(pos).Dot(normal)
	max := min
	for j := 0; j < len(boundingBox); j++ {
		x := boundingBox[j].Add(pos).Dot(normal)
		if x > max {
			max = x
		} else if x < min {
			min = x
		}
	}

	return min, max
}

func (s *System) ComponentTypes() map[string]entity.CmptTypeCfg {
	return map[string]entity.CmptTypeCfg{
		"physics": {
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
			"touching": {
				Coder: common.NewCoder(common.CoderConfig{
					ConfigType: &TouchingCondition{},
					Tag:        "config",
					Decode:     func(x interface{}) (interface{}, error) { return x.(*TouchingCondition), nil },
					Encode:     func(x interface{}) (interface{}, error) { return x.(*TouchingCondition), nil },
				}),
			},
		},
		Effects: map[string]triggersys.EffectTypeCfg{},
	}
}

func (s *System) WillJoinManager(em *entity.Manager) {
	s.entityManager = em

	componentQuery, err := s.entityManager.MakeCmptQuery([]string{"physics", "position"})
	if err != nil {
		panic(err)
	}
	s.componentQuery = componentQuery

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
