package mainscene

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/animationsys"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
)

type InputHandler struct {
	positionSystem   *positionsys.System
	animationSystem  *animationsys.System
	controlledEntity entity.ID
}

func (h *InputHandler) SetControlledEntity(eid entity.ID) {
	h.controlledEntity = eid
}

func (h *InputHandler) HandleInputState(t common.Time, state inputState) {
	var accelAmt float32 = 1
	if state.states[StateSprint] {
		accelAmt *= 2
	}

	if state.states[StateUp] {
		h.positionSystem.AddForce(h.controlledEntity, mgl32.Vec2{0.0, -accelAmt})
		h.animationSystem.SetAnimation(h.controlledEntity, "walking-up", t)
	}

	if state.states[StateDown] {
		h.positionSystem.AddForce(h.controlledEntity, mgl32.Vec2{0.0, accelAmt})
		h.animationSystem.SetAnimation(h.controlledEntity, "walking-down", t)
	}

	if state.states[StateLeft] {
		h.positionSystem.AddForce(h.controlledEntity, mgl32.Vec2{accelAmt, 0.0})
		h.animationSystem.SetAnimation(h.controlledEntity, "walking-left", t)
	}

	if state.states[StateRight] {
		h.positionSystem.AddForce(h.controlledEntity, mgl32.Vec2{-accelAmt, 0.0})
		h.animationSystem.SetAnimation(h.controlledEntity, "walking-right", t)
	}

	for _, x := range state.actions {
		switch x {
		case ActionFireWeapon:
			fmt.Println("Fired weapon!")
		}
	}
}
