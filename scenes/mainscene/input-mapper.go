package mainscene

import (
	"fmt"

	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/brynbellomy/gl4-game/input"
)

type (
	InputMapper struct{}

	inputState struct {
		actions []CharacterInputAction
		states  map[CharacterInputState]bool
	}

	CharacterInputAction int
	CharacterInputState  int
)

const (
	ActionFireWeapon CharacterInputAction = iota
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
		states:  map[CharacterInputState]bool{},
		actions: []CharacterInputAction{},
	}
}

func (i inputState) Clone() inputState {
	return inputState{
		states:  i.states,
		actions: []CharacterInputAction{},
	}
}

func (m *InputMapper) MapInputs(state inputState, events []input.IEvent) inputState {
	for _, evt := range events {
		switch evt := evt.(type) {

		case input.KeyEvent:
			switch evt.Key {
			case glfw.KeyUp:
				state.states[StateUp] = (evt.Action == glfw.Press || evt.Action == glfw.Repeat)
			case glfw.KeyDown:
				state.states[StateDown] = (evt.Action == glfw.Press || evt.Action == glfw.Repeat)
			case glfw.KeyLeft:
				state.states[StateLeft] = (evt.Action == glfw.Press || evt.Action == glfw.Repeat)
			case glfw.KeyRight:
				state.states[StateRight] = (evt.Action == glfw.Press || evt.Action == glfw.Repeat)

			case glfw.KeyLeftShift:
				state.states[StateSprint] = (evt.Action == glfw.Press || evt.Action == glfw.Repeat)

			case glfw.KeyZ:
				if evt.Action == glfw.Press {
					state.actions = append(state.actions, ActionFireWeapon)
				}
			}

		case input.MouseEvent:
			fmt.Printf("Mouse event: %+v\n", evt)

		default:
			panic("unknown input event type")
		}
	}
	return state
}
