package input

import (
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/brynbellomy/gl4-game/common"
)

type (
	KeyEvent struct {
		Key      glfw.Key
		Scancode int
		Action   glfw.Action
		Mods     glfw.ModifierKey
	}

	MouseEvent struct {
		MouseButton glfw.MouseButton
		Action      glfw.Action
		Mods        glfw.ModifierKey
	}

	CursorEvent struct {
		Pos common.WindowPos
	}
)
