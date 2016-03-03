package mainscene

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/input"
	"github.com/brynbellomy/gl4-game/systems/animationsys"
	"github.com/brynbellomy/gl4-game/systems/gameobjsys"
	"github.com/brynbellomy/gl4-game/systems/movesys"
	"github.com/brynbellomy/gl4-game/systems/physicssys"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
	"github.com/brynbellomy/gl4-game/systems/projectilesys"
	"github.com/brynbellomy/gl4-game/systems/rendersys"
)

type (
	MainScene struct {
		window                    *glfw.Window
		projection                mgl32.Mat4
		windowWidth, windowHeight int

		assetRoot     string
		entityManager entity.Manager

		heroID   entity.ID
		cameraID entity.ID

		inputQueue   *input.Enqueuer
		inputState   inputState
		inputMapper  InputMapper
		inputHandler InputHandler

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
	positionSystem := positionsys.New()
	physicsSystem := physicssys.New()
	renderSystem := rendersys.New()
	animationSystem := animationsys.New()
	gameobjSystem := gameobjsys.New()
	moveSystem := movesys.New()
	projectileSystem := projectilesys.New()

	entityManager := entity.NewManager([]entity.ISystem{
		positionSystem,
		physicsSystem,
		renderSystem,
		animationSystem,
		gameobjSystem,
		moveSystem,
		projectileSystem,
	})

	mainScene := &MainScene{
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

		inputState: newInputState(),
		inputHandler: InputHandler{
			moveSystem:     moveSystem,
			positionSystem: positionSystem,
			gameobjSystem:  gameobjSystem,
		},
		inputQueue: input.NewEnqueuer(),

		fireballFactory: NewFireballFactory(assetRoot),
	}

	return mainScene, nil
}

func (s *MainScene) Prepare() error {
	ww, wh := s.window.GetSize()
	s.windowWidth = ww
	s.windowHeight = wh

	s.projection = mgl32.Perspective(mgl32.DegToRad(45.0), float32(ww)/float32(wh), 0.1, 10.0)
	s.renderSystem.SetProjection(s.projection)

	{
		s.cameraID = s.entityManager.NewEntityID()
		s.entityManager.AddComponents(s.cameraID, []entity.IComponent{
			positionsys.NewComponent(mgl32.Vec2{0, 0}, common.Size{0, 0}, 0),
		})
		s.renderSystem.SetCameraPos(mgl32.Vec2{0, 0})
	}

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

	// {
	// 	fireballCmpts, err := fireball(s.assetRoot, mgl32.Vec2{0, 0}, mgl32.Vec2{1, 1})
	// 	if err != nil {
	// 		return err
	// 	}
	// 	fireballID := s.entityManager.NewEntityID()
	// 	s.entityManager.AddComponents(fireballID, fireballCmpts)
	// }

	s.inputQueue.BecomeInputResponder(s.window)
	s.inputHandler.SetControlledEntity(s.heroID)

	s.inputHandler.onFireWeapon = s.onFireWeapon

	return nil
}

func (s *MainScene) getWorldPos(windowPos mgl32.Vec2) (mgl32.Vec2, error) {
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

	cameraPos := s.getCameraPos()
	viewMatrix := mgl32.Ident4().Mul4(mgl32.Translate3D(cameraPos.X(), -cameraPos.Y(), 0))

	flippedY := float32(s.windowHeight) - windowPos.Y()
	worldPos, err := mgl32.UnProject(
		mgl32.Vec3{windowPos.X(), flippedY, depthBuf[0]}, // win coords
		viewMatrix,     // modelview
		s.projection,   // projection
		0,              // initialX
		0,              // initialY
		s.windowWidth,  // width
		s.windowHeight, // height
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
	return s.positionSystem.GetPos(s.heroID)
}

func (s *MainScene) onFireWeapon(controlledEntity entity.ID, x ActionFireWeapon) {
	targetPos, err := s.getWorldPos(x.Target)
	if err != nil {
		panic(err)
	}

	pos := s.positionSystem.GetPos(controlledEntity)
	vec := pos.Sub(targetPos)

	fireball, err := s.fireballFactory.Build(pos, vec)
	if err != nil {
		// @@TODO
		panic(err)
	}

	s.entityManager.AddComponents(s.entityManager.NewEntityID(), fireball)
}

func (s *MainScene) Update() {
	t := common.Now()

	// update input
	s.inputState = s.inputMapper.MapInputs(s.inputState.Clone(), s.inputQueue.FlushEvents())
	s.inputHandler.HandleInputState(t, s.inputState)

	s.gameobjSystem.Update(t)
	s.projectileSystem.Update(t)
	s.moveSystem.Update(t)
	s.physicsSystem.Update(t)
	s.positionSystem.Update(t)
	s.animationSystem.Update(t)

	s.renderSystem.SetCameraPos(s.positionSystem.GetPos(s.heroID))
	s.renderSystem.Update(t)
}
