package main

import (
	"fmt"
	_ "image/png"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/input"
	"github.com/brynbellomy/gl4-game/node"
	"github.com/brynbellomy/gl4-game/scene"
)

const windowWidth = 1280
const windowHeight = 960

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)

	scene := scene.New(node.Config{
		Size: &common.Size{Width: 2.0, Height: 2.0},
	})

	bg := node.NewSpriteNode(node.SpriteNodeConfig{
		Projection:  projection,
		Size:        &common.Size{Width: 2.0, Height: 2.0},
		Position:    mgl32.Vec2{0.0, 0.0},
		TextureFile: "square.png",
	})

	hero := node.NewSpriteNode(node.SpriteNodeConfig{
		Projection:  projection,
		Size:        &common.Size{Width: 0.2, Height: 0.4},
		Position:    mgl32.Vec2{0.0, 0.0},
		TextureFile: "resources/textures/lavos/walking-down-001.png",
	})
	scene.AddChild(bg)
	bg.AddChild(hero)

	inputHandler := input.NewHandler()
	window.SetKeyCallback(inputHandler.OnKey)
	scene.SetInputHandler(inputHandler)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		scene.Update()

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
