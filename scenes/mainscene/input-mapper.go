package mainscene

import (
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/input"
)

type (
	InputMapper struct{}

	inputState struct {
		actions   []ICharacterInputAction
		states    map[CharacterInputState]bool
		cursorPos common.WindowPos
	}

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

func newInputState() inputState {
	return inputState{
		states:    map[CharacterInputState]bool{},
		actions:   []ICharacterInputAction{},
		cursorPos: common.WindowPos{},
	}
}

func (i inputState) Clone() inputState {
	return inputState{
		states:    i.states,
		actions:   []ICharacterInputAction{},
		cursorPos: i.cursorPos,
	}
}

func (m *InputMapper) MapInputs(state inputState, events []input.IEvent) inputState {
	for _, evt := range events {
		switch evt := evt.(type) {

		case input.KeyEvent:
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

		case input.MouseEvent:
			switch evt.MouseButton {
			case glfw.MouseButton1:
				if evt.Action == glfw.Press {
					state.actions = append(state.actions, ActionFireWeapon{WindowPos: state.cursorPos})
				}
			}

		case input.CursorEvent:
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
