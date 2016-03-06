package mainscene

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/animationsys"
	"github.com/brynbellomy/gl4-game/systems/gameobjsys"
	"github.com/brynbellomy/gl4-game/systems/inputsys"
	"github.com/brynbellomy/gl4-game/systems/movesys"
	"github.com/brynbellomy/gl4-game/systems/physicssys"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
	"github.com/brynbellomy/gl4-game/systems/projectilesys"
	"github.com/brynbellomy/gl4-game/systems/rendersys"
)

type (
	MainScene struct {
		window     *glfw.Window
		projection mgl32.Mat4

		assetRoot     string
		entityManager *entity.Manager

		heroID   entity.ID
		cameraID entity.ID

		inputSystem  *inputsys.System
		inputHandler *InputHandler

		positionSystem   *positionsys.System
		physicsSystem    *physicssys.System
		renderSystem     *rendersys.System
		animationSystem  *animationsys.System
		gameobjSystem    *gameobjsys.System
		moveSystem       *movesys.System
		projectileSystem *projectilesys.System

		fireballFactory *FireballFactory
	}
)

func NewMainScene(window *glfw.Window, assetRoot string) (*MainScene, error) {
	fireballFactory, err := NewFireballFactory(assetRoot)
	if err != nil {
		return nil, err
	}

	var (
		positionSystem   = positionsys.New()
		physicsSystem    = physicssys.New()
		renderSystem     = rendersys.New()
		animationSystem  = animationsys.New()
		gameobjSystem    = gameobjsys.New()
		moveSystem       = movesys.New()
		projectileSystem = projectilesys.New()

		inputMapper  = &InputMapper{}
		inputHandler = &InputHandler{
			moveSystem:     moveSystem,
			positionSystem: positionSystem,
			gameobjSystem:  gameobjSystem,
		}

		inputSystem = inputsys.New(newInputState(), inputMapper, inputHandler)

		entityManager = entity.NewManager([]entity.ISystem{
			positionSystem,
			physicsSystem,
			renderSystem,
			animationSystem,
			gameobjSystem,
			moveSystem,
			projectileSystem,
		})

		mainScene = &MainScene{
			window: window,

			assetRoot:     assetRoot,
			entityManager: entityManager,

			positionSystem:   positionSystem,
			physicsSystem:    physicsSystem,
			renderSystem:     renderSystem,
			animationSystem:  animationSystem,
			gameobjSystem:    gameobjSystem,
			moveSystem:       moveSystem,
			projectileSystem: projectileSystem,

			inputSystem:  inputSystem,
			inputHandler: inputHandler,

			fireballFactory: fireballFactory,
		}
	)

	return mainScene, nil
}

func (s *MainScene) Prepare() error {
	ww, wh := s.window.GetSize()

	s.projection = mgl32.Perspective(mgl32.DegToRad(45.0), float32(ww)/float32(wh), 0.1, 10.0)
	s.renderSystem.SetProjection(s.projection)

	{
		bgCmpts, err := bg(s.assetRoot)
		if err != nil {
			return err
		}
		bgID := s.entityManager.NewEntityID()
		s.entityManager.AddComponents(bgID, bgCmpts)
	}

	{
		heroCmpts, err := hero(s.assetRoot)
		if err != nil {
			return err
		}
		s.heroID = s.entityManager.NewEntityID()
		s.entityManager.AddComponents(s.heroID, heroCmpts)
	}

	{
		skeletonCmpts, err := skeleton(s.assetRoot)
		if err != nil {
			return err
		}
		skeletonID := s.entityManager.NewEntityID()
		s.entityManager.AddComponents(skeletonID, skeletonCmpts)
	}

	s.cameraID = s.heroID

	s.inputSystem.BecomeInputResponder(s.window)
	s.inputSystem.SetControlledEntity(s.heroID)
	s.inputHandler.onFireWeapon = s.onFireWeapon

	s.physicsSystem.OnCollision(func(c physicssys.Collision) {
		// fmt.Printf("collision ~> %+v\n", c)
	})

	return nil
}

func (s *MainScene) getWorldPos(windowPos common.WindowPos) (mgl32.Vec2, error) {
	depthBuf := make([]float32, 1)
	gl.ReadPixels(
		int32(windowPos.X()),
		int32(windowPos.Y()),
		1,
		1,
		gl.DEPTH_COMPONENT,
		gl.FLOAT,
		gl.Ptr(depthBuf),
	)

	// the model matrix of the world is always set to identity, so the "model * view" parameter is just the camera transform
	modelview := s.getCamera()

	windowWidth, windowHeight := s.window.GetSize()
	flippedY := float64(windowHeight) - windowPos.Y() // we have to flip the Y value because our coordinate system is oriented differently from OpenGL's

	worldPos, err := mgl32.UnProject(
		mgl32.Vec3{float32(windowPos.X()), float32(flippedY), depthBuf[0]}, // win coords
		modelview,    // modelview
		s.projection, // projection
		0,            // initialX
		0,            // initialY
		windowWidth,  // width
		windowHeight, // height
	)

	if err != nil {
		return mgl32.Vec2{}, err
	}

	return mgl32.Vec2{worldPos.X(), worldPos.Y()}, nil
}

// @@TODO
// @@TODO
// @@TODO
func (s *MainScene) getCameraPos() mgl32.Vec2 {
	return s.positionSystem.GetPos(s.cameraID)
}

func (s *MainScene) getCamera() mgl32.Mat4 {
	pos := s.getCameraPos()
	return mgl32.LookAtV(mgl32.Vec3{pos.X(), pos.Y(), 3}, mgl32.Vec3{pos.X(), pos.Y(), 0}, mgl32.Vec3{0, -1, 0})
}

func (s *MainScene) onFireWeapon(controlledEntity entity.ID, x ActionFireWeapon) {
	targetPos, err := s.getWorldPos(x.WindowPos)
	if err != nil {
		panic(err)
	}

	pos := s.positionSystem.GetPos(controlledEntity)
	vec := targetPos.Sub(pos)

	fireball, err := s.fireballFactory.Build(pos, vec)
	if err != nil {
		// @@TODO
		panic(err)
	}

	s.entityManager.AddComponents(s.entityManager.NewEntityID(), fireball)
}

func (s *MainScene) Update() {
	t := common.Now()

	s.inputSystem.Update(t)

	s.gameobjSystem.Update(t)
	s.projectileSystem.Update(t)
	s.moveSystem.Update(t)
	s.physicsSystem.Update(t)
	s.positionSystem.Update(t)
	s.animationSystem.Update(t)

	s.renderSystem.SetCamera(s.getCamera())
	s.renderSystem.Update(t)

	s.entityManager.CullEntities()
}
