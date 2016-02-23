package main

import (
	"fmt"
	_ "image/png"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
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

const windowWidth = 1280
const windowHeight = 960

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	assetPath, err := getAssetPath()
	if err != nil {
		panic(err)
	}

	window, err := initGLFWWindow()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	s, err := initScene(window, assetPath)
	if err != nil {
		panic(err)
	}

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		s.Update()

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func initGLFWWindow() (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "xyzzy", nil, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()
	return window, nil
}

func getAssetPath() (string, error) {
	exePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}

	p, err := filepath.Abs(filepath.Dir(exePath))
	return p, err
}

func initScene(window *glfw.Window, assetRoot string) (scene.IScene, error) {
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)

	inputHandler := input.NewHandler()
	window.SetKeyCallback(inputHandler.OnKey)

	heroSize := common.Size{0.2, 0.4}
	heroPos := mgl32.Vec2{0.0, 0.0}
	hero := node.NewSpriteNode(node.SpriteNodeConfig{
		Projection: projection,
		Size:       heroSize,
		Position:   heroPos,
	})

	heroAtlas := texture.NewAtlas()

	err := heroAtlas.LoadAnimation("walking", []string{
		path.Join(assetRoot, "resources/textures/lavos/walking-down-001.png"),
		path.Join(assetRoot, "resources/textures/lavos/walking-down-002.png"),
		path.Join(assetRoot, "resources/textures/lavos/walking-down-003.png"),
		path.Join(assetRoot, "resources/textures/lavos/walking-down-004.png"),
		path.Join(assetRoot, "resources/textures/lavos/walking-left-001.png"),
		path.Join(assetRoot, "resources/textures/lavos/walking-left-002.png"),
		path.Join(assetRoot, "resources/textures/lavos/walking-left-003.png"),
		path.Join(assetRoot, "resources/textures/lavos/walking-left-004.png"),
		path.Join(assetRoot, "resources/textures/lavos/walking-up-001.png"),
		path.Join(assetRoot, "resources/textures/lavos/walking-up-002.png"),
		path.Join(assetRoot, "resources/textures/lavos/walking-up-003.png"),
		path.Join(assetRoot, "resources/textures/lavos/walking-up-004.png"),
		path.Join(assetRoot, "resources/textures/lavos/walking-right-001.png"),
		path.Join(assetRoot, "resources/textures/lavos/walking-right-002.png"),
		path.Join(assetRoot, "resources/textures/lavos/walking-right-003.png"),
		path.Join(assetRoot, "resources/textures/lavos/walking-right-004.png"),
	})

	if err != nil {
		panic(err.Error())
	}

	bgSize := common.Size{2.0, 2.0}
	bgPos := mgl32.Vec2{0.0, 0.0}

	bgTexture, err := texture.Load(path.Join(assetRoot, "square.png"))
	if err != nil {
		panic(err.Error())
	}

	bg := node.NewSpriteNode(node.SpriteNodeConfig{
		Projection: projection,
		Size:       bgSize,
		Position:   bgPos,
	})

	renderRoot := node.New()
	renderRoot.AddChild(bg)
	bg.AddChild(hero)

	s := scene.New(scene.Config{
		InputHandler: inputHandler,
		Systems: []entity.ISystem{
			positionsys.New(),
			rendersys.New(renderRoot),
			animationsys.New(),
		},
	})

	s.AddEntity(entity.ID(0), []entity.IComponent{
		positionsys.NewComponent(mgl32.Vec2{0, 0}, bgSize),
		rendersys.NewComponent(bg, bgTexture),
	})

	s.AddEntity(entity.ID(1), []entity.IComponent{
		positionsys.NewComponent(mgl32.Vec2{0, 0}, heroSize),
		rendersys.NewComponent(hero, 0),
		animationsys.NewComponent(heroAtlas, "walking", 0, 2),
	})

	return s, nil
}
