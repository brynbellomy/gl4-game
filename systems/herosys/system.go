package herosys

import (
    "github.com/go-gl/mathgl/mgl32"

    "github.com/brynbellomy/gl4-game/context"
    "github.com/brynbellomy/gl4-game/entity"
)

type (
    System struct {
        entities  []entityAspect
        entityMap map[entity.ID]*entityAspect
    }

    entityAspect struct {
        heroCmpt *Component
        positionCmpt *positionsys.Component
    }
)

func New() *System {
    return &System{
        entities:  make([]entityAspect, 0),
        entityMap: make(map[entity.ID]*entityAspect),
    }
}

func (s *System) GetPos(eid entity.ID) mgl32.Vec2 {
    if e, exists := s.entityMap[eid]; exists {
        return e.heroCmpt.Pos()
    } else {
        panic("entity does not exist")
    }
}

func (s *System) SetPos(eid entity.ID, pos mgl32.Vec2) {
    if e, exists := s.entityMap[eid]; exists {
        e.heroCmpt.SetPos(pos)
    } else {
        panic("entity does not exist")
    }
}

func (s *System) Update(c context.IContext) {
    for _, e := range s.entities {
        vel := e.heroCmpt.Velocity().Add(e.heroCmpt.Acceleration())

        newpos := e.heroCmpt.Pos().Add(vel)
        e.heroCmpt.SetPos(newpos)
    }
}

func (s *System) ComponentsWillJoin(eid entity.ID, components []entity.IComponent) {
    // if we find a *positionsys.Component on the entity, we keep track of it
    var heroCmpt *Component

    for _, cmpt := range components {
        if cmpt, is := cmpt.(*Component); is {
            heroCmpt = cmpt
            break
        }
    }

    if heroCmpt != nil {
        aspect := entityAspect{heroCmpt: heroCmpt}
        s.entities = append(s.entities, aspect)
        s.entityMap[eid] = &s.entities[len(s.entities)-1]
    }
}
