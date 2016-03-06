package mainscene

import (
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/systems/inputsys"
)

type (
	InputMapper struct{}

	CharacterInputState int

	ICharacterInputAction interface{}
)

const (
	StateUp CharacterInputState = iota
	StateDown
	StateLeft
	StateRight
	StateSprint
)

func (m *InputMapper) MapInputs(st inputsys.IInputState, events []inputsys.IEvent) inputsys.IInputState {
	state := st.(inputState)

	for _, evt := range events {
		switch evt := evt.(type) {

		case inputsys.KeyEvent:
			switch evt.Key {
			case glfw.KeyW:
				state.states[StateUp] = (evt.Action == glfw.Press || evt.Action == glfw.Repeat)
			case glfw.KeyS:
				state.states[StateDown] = (evt.Action == glfw.Press || evt.Action == glfw.Repeat)
			case glfw.KeyA:
				state.states[StateLeft] = (evt.Action == glfw.Press || evt.Action == glfw.Repeat)
			case glfw.KeyD:
				state.states[StateRight] = (evt.Action == glfw.Press || evt.Action == glfw.Repeat)

			case glfw.KeyLeftShift:
				state.states[StateSprint] = (evt.Action == glfw.Press || evt.Action == glfw.Repeat)
			}

		case inputsys.MouseEvent:
			switch evt.MouseButton {
			case glfw.MouseButton1:
				if evt.Action == glfw.Press {
					state.actions = append(state.actions, ActionFireWeapon{WindowPos: state.cursorPos})
				}
			}

		case inputsys.CursorEvent:
			state.cursorPos = evt.Pos

		default:
			panic("unknown input event type")
		}
	}
	return state
}

type (
	ActionFireWeapon struct {
		WindowPos common.WindowPos
	}
)
