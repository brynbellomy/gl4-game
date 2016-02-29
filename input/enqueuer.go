package input

import "github.com/go-gl/glfw/v3.1/glfw"

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
}

func (h *Enqueuer) OnKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	h.queuedEvents = append(h.queuedEvents, KeyEvent{key, scancode, action, mods})
}

func (h *Enqueuer) OnMouseButton(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	h.queuedEvents = append(h.queuedEvents, MouseEvent{button, action, mod})
}

func (h *Enqueuer) FlushEvents() []IEvent {
	evts := h.queuedEvents
	h.queuedEvents = []IEvent{}
	return evts
}
