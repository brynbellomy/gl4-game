package mainscene

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/input"
	"github.com/brynbellomy/gl4-game/scene"
	"github.com/brynbellomy/gl4-game/systems/animationsys"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
	"github.com/brynbellomy/gl4-game/systems/rendersys"
)

type (
	MainScene struct {
		*scene.Scene
		window *glfw.Window

		heroID entity.ID

		inputQueue   *input.Enqueuer
		inputState   inputState
		inputMapper  InputMapper
		inputHandler InputHandler

		positionSystem  *positionsys.System
		renderSystem    *rendersys.System
		animationSystem *animationsys.System
		// heroSystem      *herosys.System
		// steeringSystem  *steeringsys.System
	}
)

func NewMainScene(window *glfw.Window, assetRoot string) (*MainScene, error) {
	positionSystem := positionsys.New()
	renderSystem := rendersys.New()
	animationSystem := animationsys.New()
	// heroSystem := herosys.New()
	// steeringSystem := steeringsys.New()

	scn := scene.New(scene.Config{
		AssetRoot: assetRoot,
		Systems: []entity.ISystem{
			positionSystem,
			renderSystem,
			animationSystem,
			// heroSystem,
			// steeringSystem,
		},
	})

	mainScene := &MainScene{
		Scene:           scn,
		window:          window,
		positionSystem:  positionSystem,
		renderSystem:    renderSystem,
		animationSystem: animationSystem,
		// steeringSystem:  steeringSystem,
		// heroSystem:      heroSystem,
		inputState: newInputState(),
		inputHandler: InputHandler{
			positionSystem:  positionSystem,
			animationSystem: animationSystem,
		},
		inputQueue: input.NewEnqueuer(),
	}

	return mainScene, nil
}

func (s *MainScene) Prepare() error {
	ww, wh := s.window.GetSize()
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(ww)/float32(wh), 0.1, 10.0)
	s.renderSystem.SetProjection(projection)

	camera := mgl32.LookAtV(mgl32.Vec3{0, 0, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, -1, 0})
	s.renderSystem.SetCamera(camera)

	{
		bgCmpts, err := bg(s.AssetRoot())
		if err != nil {
			return err
		}
		s.Scene.EntityManager().AddComponents(entity.ID(0), bgCmpts)
	}

	{
		heroCmpts, err := hero(s.AssetRoot())
		if err != nil {
			return err
		}
		s.heroID = entity.ID(1)
		s.Scene.EntityManager().AddComponents(s.heroID, heroCmpts)
	}

	// s.steeringSystem.AddBehavior(s.heroID, &steeringbehaviors.Constant{Vec: mgl32.Vec2{0.005, 0.005}})

	s.inputQueue.BecomeInputResponder(s.window)
	s.inputHandler.SetControlledEntity(s.heroID)

	return nil
}

func (s *MainScene) Update() {
	t := common.Now() //time.Now().UTC().UnixNano()

	// update input
	s.inputState = s.inputMapper.MapInputs(s.inputState.Clone(), s.inputQueue.FlushEvents())
	s.inputHandler.HandleInputState(t, s.inputState)

	s.Scene.Update(t)
}
