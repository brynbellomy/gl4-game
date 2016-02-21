package input

import (
	"fmt"

	"github.com/go-gl/glfw/v3.1/glfw"
)

type (
	Handler struct {
		queuedEvents []*inputEvent
	}

	IHandler interface {
		Update()
	}

	inputEvent struct {
		key      glfw.Key
		scancode int
		action   glfw.Action
		mods     glfw.ModifierKey
	}
)

func NewHandler() *Handler {
	return &Handler{
		queuedEvents: []*inputEvent{},
	}
}

func (h *Handler) OnKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	h.queuedEvents = append(h.queuedEvents, &inputEvent{key, scancode, action, mods})
}

func (h *Handler) Update() {
	for _, e := range h.queuedEvents {
		h.handleEvent(e)
	}

	h.queuedEvents = []*inputEvent{}
}

func (h *Handler) handleEvent(e *inputEvent) {
	fmt.Println("KEY:", e.key, e.action)
}
