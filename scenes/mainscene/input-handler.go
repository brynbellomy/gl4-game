package mainscene

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/gameobjsys"
	"github.com/brynbellomy/gl4-game/systems/inputsys"
	"github.com/brynbellomy/gl4-game/systems/movesys"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
)

type InputHandler struct {
	moveSystem       *movesys.System
	positionSystem   *positionsys.System
	gameobjSystem    *gameobjsys.System
	controlledEntity entity.ID
	onFireWeapon     func(controlledEntity entity.ID, x ActionFireWeapon)
}

func NewInputHandler(moveSystem *movesys.System, positionSystem *positionsys.System, gameobjSystem *gameobjsys.System) *InputHandler {
	return &InputHandler{
		moveSystem:       moveSystem,
		positionSystem:   positionSystem,
		gameobjSystem:    gameobjSystem,
		controlledEntity: entity.InvalidID,
	}
}

func (h *InputHandler) SetControlledEntity(eid entity.ID) {
	h.controlledEntity = eid
}

func (h *InputHandler) HandleInputState(t common.Time, st inputsys.IInputState) {
	state := st.(inputState)

	var accelAmt float32 = 1
	if state.states[StateSprint] {
		accelAmt *= 2
	}

	var totalAccel mgl32.Vec2
	if state.states[StateUp] {
		totalAccel = totalAccel.Add(mgl32.Vec2{0.0, -accelAmt})
	}

	if state.states[StateDown] {
		totalAccel = totalAccel.Add(mgl32.Vec2{0.0, accelAmt})
	}

	if state.states[StateLeft] {
		totalAccel = totalAccel.Add(mgl32.Vec2{accelAmt, 0.0})
	}

	if state.states[StateRight] {
		totalAccel = totalAccel.Add(mgl32.Vec2{-accelAmt, 0.0})
	}

	if h.controlledEntity != entity.InvalidID {
		h.moveSystem.SetMovementVector(h.controlledEntity, totalAccel)
	}

	for _, x := range state.actions {
		switch x := x.(type) {
		case ActionFireWeapon:
			if h.onFireWeapon != nil {
				h.onFireWeapon(h.controlledEntity, x)
			}
		}
	}
}
