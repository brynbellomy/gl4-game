package mainscene

import (
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/input"
	"github.com/brynbellomy/gl4-game/scene"
	"github.com/brynbellomy/gl4-game/systems/animationsys"
	"github.com/brynbellomy/gl4-game/systems/gameobjsys"
	"github.com/brynbellomy/gl4-game/systems/physicssys"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
	"github.com/brynbellomy/gl4-game/systems/rendersys"
)

type (
	MainScene struct {
		*scene.Scene
		window *glfw.Window

		heroID   entity.ID
		cameraID entity.ID

		inputQueue   *input.Enqueuer
		inputState   inputState
		inputMapper  InputMapper
		inputHandler InputHandler

		positionSystem  *positionsys.System
		physicsSystem   *physicssys.System
		renderSystem    *rendersys.System
		animationSystem *animationsys.System
		gameobjSystem   *gameobjsys.System
	}
)

func NewMainScene(window *glfw.Window, assetRoot string) (*MainScene, error) {
	positionSystem := positionsys.New()
	physicsSystem := physicssys.New()
	renderSystem := rendersys.New()
	animationSystem := animationsys.New()
	gameobjSystem := gameobjsys.New()

	scn := scene.New(scene.Config{
		AssetRoot: assetRoot,
		Systems: []entity.ISystem{
			positionSystem,
			physicsSystem,
			renderSystem,
			animationSystem,
			gameobjSystem,
		},
	})

	mainScene := &MainScene{
		Scene:  scn,
		window: window,

		positionSystem:  positionSystem,
		physicsSystem:   physicsSystem,
		renderSystem:    renderSystem,
		animationSystem: animationSystem,
		gameobjSystem:   gameobjSystem,

		inputState: newInputState(),
		inputHandler: InputHandler{
			physicsSystem: physicsSystem,
			gameobjSystem: gameobjSystem,
		},
		inputQueue: input.NewEnqueuer(),
	}

	return mainScene, nil
}

func (s *MainScene) Prepare() error {
	ww, wh := s.window.GetSize()
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(ww)/float32(wh), 0.1, 10.0)
	s.renderSystem.SetProjection(projection)

	{
		s.cameraID = entity.ID(0)
		s.Scene.EntityManager().AddComponents(s.cameraID, []entity.IComponent{
			positionsys.NewComponent(mgl32.Vec2{0, 0}, common.Size{0, 0}, 0),
		})
		s.renderSystem.SetCameraPos(mgl32.Vec2{0, 0})
	}

	{
		bgCmpts, err := bg(s.AssetRoot())
		if err != nil {
			return err
		}
		s.Scene.EntityManager().AddComponents(entity.ID(1), bgCmpts)
	}

	{
		heroCmpts, err := hero(s.AssetRoot())
		if err != nil {
			return err
		}
		s.heroID = entity.ID(2)
		s.Scene.EntityManager().AddComponents(s.heroID, heroCmpts)
	}

	go func() {
		time.Sleep(2 * time.Second)
		s.physicsSystem.AddForce(s.heroID, mgl32.Vec2{10, 10})
	}()

	s.inputQueue.BecomeInputResponder(s.window)
	s.inputHandler.SetControlledEntity(s.heroID)

	return nil
}

func (s *MainScene) Update() {
	t := common.Now()

	// update input
	s.inputState = s.inputMapper.MapInputs(s.inputState.Clone(), s.inputQueue.FlushEvents())
	s.inputHandler.HandleInputState(t, s.inputState)

	s.gameobjSystem.Update(t)
	s.physicsSystem.Update(t)
	s.positionSystem.Update(t)
	s.animationSystem.Update(t)

	s.renderSystem.SetCameraPos(s.positionSystem.GetPos(s.heroID))
	s.renderSystem.Update(t)
}
