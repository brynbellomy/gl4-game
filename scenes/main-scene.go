package scenes

import (
	"fmt"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/input"
	"github.com/brynbellomy/gl4-game/node"
	"github.com/brynbellomy/gl4-game/scene"
	"github.com/brynbellomy/gl4-game/systems/animationsys"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
	"github.com/brynbellomy/gl4-game/systems/rendersys"
	"github.com/brynbellomy/gl4-game/texture"
)

type (
	MainScene struct {
		*scene.Scene
		window *glfw.Window

		heroID entity.ID

		inputQueue input.IEnqueuer
		inputState inputState

		positionSystem  *positionsys.System
		renderSystem    *rendersys.System
		animationSystem *animationsys.System
	}

	inputState struct {
		up, down, left, right bool
	}
)

func NewMainScene(window *glfw.Window, assetRoot string) (*MainScene, error) {
	windowWidth, windowHeight := window.GetSize()
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), 0.1, 10.0)

	positionSystem := positionsys.New()
	renderSystem := rendersys.New()
	animationSystem := animationsys.New()

	scn := scene.New(scene.Config{
		AssetRoot:  assetRoot,
		Projection: projection,
		Systems: []entity.ISystem{
			positionSystem,
			renderSystem,
			animationSystem,
		},
	})

	mainScene := &MainScene{
		Scene:           scn,
		inputQueue:      input.NewEnqueuer(),
		window:          window,
		positionSystem:  positionSystem,
		renderSystem:    renderSystem,
		animationSystem: animationSystem,
		inputState:      inputState{},
	}

	return mainScene, nil
}

func (s *MainScene) Prepare() error {
	heroSize := common.Size{0.2, 0.4}
	heroPos := mgl32.Vec2{0.0, 0.0}
	hero := node.NewSpriteNode(node.SpriteNodeConfig{
		Projection: s.Scene.Projection(),
		Size:       heroSize,
		Position:   heroPos,
	})

	heroAtlas := texture.NewAtlas()
	err := heroAtlas.LoadAnimation("walking-down", []string{
		s.Scene.AssetPath("textures/lavos/walking-down-001.png"),
		s.Scene.AssetPath("textures/lavos/walking-down-002.png"),
		s.Scene.AssetPath("textures/lavos/walking-down-003.png"),
		s.Scene.AssetPath("textures/lavos/walking-down-004.png"),
	})
	if err != nil {
		return err
	}

	err = heroAtlas.LoadAnimation("walking-left", []string{
		s.Scene.AssetPath("textures/lavos/walking-left-001.png"),
		s.Scene.AssetPath("textures/lavos/walking-left-002.png"),
		s.Scene.AssetPath("textures/lavos/walking-left-003.png"),
		s.Scene.AssetPath("textures/lavos/walking-left-004.png"),
	})
	if err != nil {
		return err
	}

	err = heroAtlas.LoadAnimation("walking-up", []string{
		s.Scene.AssetPath("textures/lavos/walking-up-001.png"),
		s.Scene.AssetPath("textures/lavos/walking-up-002.png"),
		s.Scene.AssetPath("textures/lavos/walking-up-003.png"),
		s.Scene.AssetPath("textures/lavos/walking-up-004.png"),
	})
	if err != nil {
		return err
	}

	err = heroAtlas.LoadAnimation("walking-right", []string{
		s.Scene.AssetPath("textures/lavos/walking-right-001.png"),
		s.Scene.AssetPath("textures/lavos/walking-right-002.png"),
		s.Scene.AssetPath("textures/lavos/walking-right-003.png"),
		s.Scene.AssetPath("textures/lavos/walking-right-004.png"),
	})
	if err != nil {
		return err
	}

	bgSize := common.Size{2.0, 2.0}
	bgPos := mgl32.Vec2{0.0, 0.0}

	bgTexture, err := texture.Load(s.Scene.AssetPath("textures/square.png"))
	if err != nil {
		return err
	}

	bg := node.NewSpriteNode(node.SpriteNodeConfig{
		Projection: s.Scene.Projection(),
		Size:       bgSize,
		Position:   bgPos,
	})

	renderRoot := node.New()
	s.renderSystem.SetRenderRoot(renderRoot)

	renderRoot.AddChild(bg)
	bg.AddChild(hero)

	s.Scene.AddEntity(entity.ID(0), []entity.IComponent{
		positionsys.NewComponent(mgl32.Vec2{0, 0}, bgSize),
		rendersys.NewComponent(bg, bgTexture),
	})

	s.heroID = entity.ID(1)
	s.Scene.AddEntity(s.heroID, []entity.IComponent{
		positionsys.NewComponent(mgl32.Vec2{0, 0}, heroSize),
		rendersys.NewComponent(hero, 0),
		animationsys.NewComponent(heroAtlas, "walking", 0, 2),
	})

	s.inputQueue.BecomeInputResponder(s.window)

	return nil
}

func (s *MainScene) Update() {
	s.updateInput()

	// this updates all of the systems (position, render, etc.)
	s.Scene.Update()
}

func (s *MainScene) updateInput() {
	events := s.inputQueue.FlushEvents()
	for _, evt := range events {
		switch evt := evt.(type) {
		case input.KeyEvent:
			s.handleKeyEvent(evt)
		case input.MouseEvent:
			s.handleMouseEvent(evt)
		}
	}

	const moveDist float32 = 0.01
	if s.inputState.up {
		pos := s.positionSystem.GetPos(s.heroID)
		pos[1] -= moveDist
		s.positionSystem.SetPos(s.heroID, pos)
		s.animationSystem.SetAnimation(s.heroID, "walking-up")
	}
	if s.inputState.down {
		pos := s.positionSystem.GetPos(s.heroID)
		pos[1] += moveDist
		s.positionSystem.SetPos(s.heroID, pos)
		s.animationSystem.SetAnimation(s.heroID, "walking-down")
	}
	if s.inputState.left {
		pos := s.positionSystem.GetPos(s.heroID)
		pos[0] += moveDist
		s.positionSystem.SetPos(s.heroID, pos)
		s.animationSystem.SetAnimation(s.heroID, "walking-left")
	}
	if s.inputState.right {
		pos := s.positionSystem.GetPos(s.heroID)
		pos[0] -= moveDist
		s.positionSystem.SetPos(s.heroID, pos)
		s.animationSystem.SetAnimation(s.heroID, "walking-right")
	}
}

func (s *MainScene) handleKeyEvent(evt input.KeyEvent) {
	switch evt.Key {
	case glfw.KeyUp:
		s.inputState.up = (evt.Action == glfw.Press || evt.Action == glfw.Repeat)
	case glfw.KeyDown:
		s.inputState.down = (evt.Action == glfw.Press || evt.Action == glfw.Repeat)
	case glfw.KeyLeft:
		s.inputState.left = (evt.Action == glfw.Press || evt.Action == glfw.Repeat)
	case glfw.KeyRight:
		s.inputState.right = (evt.Action == glfw.Press || evt.Action == glfw.Repeat)
	}
	// fmt.Printf("Key event: %+v\n", evt)
	fmt.Printf("%+v\n", s.inputState)
}

func (s *MainScene) handleMouseEvent(evt input.MouseEvent) {
	fmt.Printf("Mouse event: %+v\n", evt)
}
