package inputsys

import (
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/input"
)

type (
	System struct {
		inputQueue   *input.Enqueuer
		inputState   IInputState
		inputMapper  IInputMapper
		inputHandler IInputHandler
	}

	IInputMapper interface {
		MapInputs(state IInputState, events []input.IEvent) IInputState
	}

	IInputState interface {
		Clone() IInputState
	}

	IInputHandler interface {
		HandleInputState(t common.Time, state IInputState)
		SetControlledEntity(eid entity.ID)
	}
)

func New(initialState IInputState, inputMapper IInputMapper, inputHandler IInputHandler) *System {
	return &System{
		inputQueue:   input.NewEnqueuer(),
		inputState:   initialState,
		inputMapper:  inputMapper,
		inputHandler: inputHandler,
	}
}

func (s *System) BecomeInputResponder(w *glfw.Window) {
	s.inputQueue.BecomeInputResponder(w)
}

func (s *System) SetControlledEntity(eid entity.ID) {
	s.inputHandler.SetControlledEntity(eid)
}

func (s *System) SetInputMapper(mapper IInputMapper) {
	s.inputMapper = mapper
}

func (s *System) SetInputHandler(handler IInputHandler) {
	s.inputHandler = handler
}

func (s *System) Update(t common.Time) {
	s.inputState = s.inputMapper.MapInputs(s.inputState.Clone(), s.inputQueue.FlushEvents())
	s.inputHandler.HandleInputState(t, s.inputState)
}

func (s *System) WillJoinManager(em *entity.Manager) {
	// no-op
}
