package triggersys

type (
    System struct {
        entities  []entityAspect
        entityMap map[entity.ID]*entityAspect
    }

    entityAspect struct {
        id      entity.ID
        triggerCmpt *Component
    }
)


func (s *System) Update(t common.Time) {
    for i, ent := range s.entities {
        if ent.triggerCmpt.ConditionValue(t, ent.id) == true {
            ent.triggerCmpt.Execute(t, ent.id)
        }
    }
}

func (s *System) ComponentTypes() map[string]entity.IComponent {
    return map[string]entity.IComponent{
        "trigger": &Component{},
    }
}

func (s *System) WillJoinManager(em *entity.Manager) {
    // no-op
}

func (s *System) EntityComponentsChanged(eid entity.ID, components []entity.IComponent) {
    var triggerCmpt *Component

    for _, cmpt := range components {
        if ac, is := cmpt.(*Component); is {
            triggerCmpt = ac
        }

        if triggerCmpt != nil {
            break
        }
    }

    if triggerCmpt != nil {
        if _, exists := s.entityMap[eid]; !exists {
            s.entities = append(s.entities, entityAspect{
                id:      eid,
                triggerCmpt: triggerCmpt,
            })

            s.entityMap[eid] = &s.entities[len(s.entities)-1]
        }

    } else {
        if _, exists := s.entityMap[eid]; exists {
            idx := -1
            for i := range s.entities {
                if s.entities[i].id == eid {
                    idx = i
                    break
                }
            }

            if idx >= 0 {
                s.entities = append(s.entities[:idx], s.entities[idx+1:]...)
            }

            delete(s.entityMap, eid)
        }
    }
}