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

	if h.controlledEntity == entity.InvalidID {
		return
	}

	moveCmptSet, err := h.entityManager.GetComponentSet("move")
	if err != nil {
		panic(err)
	}

	e, err := moveCmptSet.Get(h.controlledEntity)
	if err != nil {
		panic(err)
	}

	moveCmpt := e.(movesys.Component)

	var heading mgl32.Vec2
	var isMoving bool
	if state.states[StateUp] {
		heading = heading.Add(mgl32.Vec2{0.0, -1.0})
		isMoving = true
	}

	if state.states[StateDown] {
		heading = heading.Add(mgl32.Vec2{0.0, 1.0})
		isMoving = true
	}

	if state.states[StateLeft] {
		heading = heading.Add(mgl32.Vec2{1.0, 0.0})
		isMoving = true
	}

	if state.states[StateRight] {
		heading = heading.Add(mgl32.Vec2{-1.0, 0.0})
		isMoving = true
	}

	if isMoving {
		if state.states[StateSprint] {
			moveCmpt.SetMovementType(movesys.MvmtSprinting)
		} else {
			moveCmpt.SetMovementType(movesys.MvmtWalking)
		}
	} else {
		moveCmpt.SetMovementType(movesys.MvmtNone)
	}

	moveCmpt.SetVector(heading)
	moveCmptSet.Set(h.controlledEntity, moveCmpt)

	for _, x := range state.actions {
		switch x := x.(type) {
		case ActionFireWeapon:
			if h.onFireWeapon != nil {
				h.onFireWeapon(h.controlledEntity, x)
			}
		}
	}
}
