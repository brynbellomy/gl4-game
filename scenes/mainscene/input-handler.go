package mainscene

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/gameobjsys"
	"github.com/brynbellomy/gl4-game/systems/physicssys"
)

type InputHandler struct {
	physicsSystem    *physicssys.System
	gameobjSystem    *gameobjsys.System
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

	var totalVelocity mgl32.Vec2
	if state.states[StateUp] {
		totalVelocity = totalVelocity.Add(mgl32.Vec2{0.0, -accelAmt})
	}

	if state.states[StateDown] {
		totalVelocity = totalVelocity.Add(mgl32.Vec2{0.0, accelAmt})
	}

	if state.states[StateLeft] {
		totalVelocity = totalVelocity.Add(mgl32.Vec2{accelAmt, 0.0})
	}

	if state.states[StateRight] {
		totalVelocity = totalVelocity.Add(mgl32.Vec2{-accelAmt, 0.0})
	}

	h.physicsSystem.AddForce(h.controlledEntity, totalVelocity)

	// if !state.states[StateUp] && !state.states[StateDown] && !state.states[StateLeft] && !state.states[StateRight] {
	// 	h.physicsSystem.SetVelocity(h.controlledEntity, mgl32.Vec2{0, 0})
	// }

	for _, x := range state.actions {
		switch x {
		case ActionFireWeapon:
			fmt.Println("Fired weapon!")
		}
	}
}
