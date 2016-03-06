package inputsys

import (
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/brynbellomy/gl4-game/common"
)

type (
	Enqueuer struct {
		queuedEvents []IEvent
	}

	IEvent interface{}
)

func NewEnqueuer() *Enqueuer {
	return &Enqueuer{
		queuedEvents: []IEvent{},
	}
}

func (h *Enqueuer) BecomeInputResponder(window *glfw.Window) {
	window.SetKeyCallback(h.OnKey)
	window.SetMouseButtonCallback(h.OnMouseButton)
	window.SetCursorPosCallback(h.OnMouseMove)
}

func (h *Enqueuer) OnKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	h.queuedEvents = append(h.queuedEvents, KeyEvent{key, scancode, action, mods})
}

func (h *Enqueuer) OnMouseButton(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	h.queuedEvents = append(h.queuedEvents, MouseEvent{button, action, mod})
}

func (h *Enqueuer) OnMouseMove(w *glfw.Window, xpos float64, ypos float64) {
	h.queuedEvents = append(h.queuedEvents, CursorEvent{
		Pos: common.WindowPos{xpos, ypos},
	})
}

func (h *Enqueuer) FlushEvents() []IEvent {
	evts := h.queuedEvents
	h.queuedEvents = []IEvent{}
	return evts
}
