package mainscene

import (
	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/systems/inputsys"
)

type (
	inputState struct {
		actions   []ICharacterInputAction
		states    map[CharacterInputState]bool
		cursorPos common.WindowPos
	}
)

func newInputState() inputsys.IInputState {
	return inputState{
		states:    map[CharacterInputState]bool{},
		actions:   []ICharacterInputAction{},
		cursorPos: common.WindowPos{},
	}
}

func (i inputState) Clone() inputsys.IInputState {
	return inputState{
		states:    i.states,
		actions:   []ICharacterInputAction{},
		cursorPos: i.cursorPos,
	}
}
