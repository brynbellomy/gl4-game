package mainscene

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/inputsys"
	"github.com/brynbellomy/gl4-game/systems/movesys"
)

type InputHandler struct {
	moveSystem *movesys.System
	// moveCmptSet      entity.IComponentSet
	entityManager    *entity.Manager
	controlledEntity entity.ID
	onFireWeapon     func(controlledEntity entity.ID, x ActionFireWeapon)
}

func NewInputHandler(moveSystem *movesys.System) *InputHandler {
	return &InputHandler{
		moveSystem: moveSystem,
		// moveCmptSet:      moveCmptSet,
		controlledEntity: entity.InvalidID,
	}
}

func (h *InputHandler) SetEntityManager(entityManager *entity.Manager) {
	h.entityManager = entityManager
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
		moveCmptSet, err := h.entityManager.GetComponentSet("move")
		if err != nil {
			panic(err)
		}

		e, err := moveCmptSet.Get(h.controlledEntity)
		if err != nil {
			panic(err)
		}

		ent := e.(movesys.Component)
		ent.SetVector(totalAccel)
		moveCmptSet.Set(h.controlledEntity, ent)
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
