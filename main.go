package main

import (
    "flag"
    "runtime/pprof"
	"fmt"
	_ "image/png"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/brynbellomy/go-speedometer"

	"github.com/brynbellomy/gl4-game/scenes/mainscene"
)

const windowWidth = 1280
const windowHeight = 960

var cpuprofile = flag.String("cpuprof", "", "write cpu profile to file")

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
    flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

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
	gl.DepthFunc(gl.LEQUAL)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	scn, err := mainscene.NewMainScene(window, assetPath)
	if err != nil {
		panic(err)
	}

	err = scn.Prepare()
	if err != nil {
		panic(err)
	}

	sp := speedometer.New()
	sp.Start()

	go func() {
		for {
			fmt.Println("fps:", sp.Speed().PerSecond())
			time.Sleep(5 * time.Second)
		}
	}()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		scn.Update()

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
		sp.Incr(1)
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
	// window, err := glfw.CreateWindow(windowWidth, windowHeight, "xyzzy", glfw.GetPrimaryMonitor(), nil)
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
	if err != nil {
		return "", err
	}

	return path.Join(p, "resources"), nil
}
