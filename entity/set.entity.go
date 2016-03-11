package entity

// type (
// 	EntitySet struct {
// 		componentQuery ComponentSet
// 		entities       []entityAspect
// 		entityMap      map[ID]*entityAspect
// 		fnMakeAspect   MakeAspectFunc
// 	}

// 	entityAspect struct {
// 		id         ID
// 		components []IComponent
// 	}

// 	MakeAspectFunc func(eid ID, cmpts []IComponent) interface{}
// )

// func NewEntitySet(cmptQuery ComponentSet, fnMakeAspect MakeAspectFunc) *EntitySet {
// 	return &EntitySet{
// 		entities:     []entityAspect{},
// 		entityMap:    map[ID]*entityAspect{},
// 		fnMakeAspect: fnMakeAspect,
// 	}
// }

// func (s *EntitySet) GetEntity(eid ID) ([]IComponent, bool) {
// 	cmpts, exists := s.entityMap[eid]
// 	return cmpts, exists
// }

// func (s *EntitySet) FilterQueriedCmpts(eid ID, components []IComponent) (aspect interface{}, queryMatch bool) {
// 	var got ComponentSet

// 	for _, cmpt := range components {
// 		if s.componentQuery.Has(cmpt.Kind()) {
// 			cmpts = append(cmpts, cmpt)
// 			got = got.Add(cmpt.Kind())
// 		}
// 	}

// 	// if we got everything we need
// 	if got.HasAll(s.componentQuery) {
// 		return s.fnMakeAspect(eid, cmpts), true
// 	} else {
// 		return nil, false
// 	}
// }

// func (s *EntitySet) EntityComponentsChanged(eid ID, components []IComponent) {
// 	cmpts := []IComponent{}
// 	var got ComponentSet

// 	for _, cmpt := range components {
// 		if s.componentQuery.Has(cmpt.Kind()) {
// 			cmpts = append(cmpts, cmpt)
// 			got = got.Add(cmpt.Kind())
// 		}
// 	}

// 	// if we got everything we need
// 	if got.HasAll(s.componentQuery) {
// 		if _, exists := s.entityMap[eid]; !exists {
// 			s.entities = append(s.entities, entityAspect{
// 				id:         eid,
// 				components: cmpts,
// 			})

// 			s.entityMap[eid] = &s.entities[len(s.entities)-1]
// 		}

// 	} else {
// 		if _, exists := s.entityMap[eid]; exists {
// 			idx := -1
// 			for i := range s.entities {
// 				if s.entities[i].id == eid {
// 					idx = i
// 					break
// 				}
// 			}

// 			if idx >= 0 {
// 				s.entities = append(s.entities[:idx], s.entities[idx+1:]...)
// 			}

// 			delete(s.entityMap, eid)
// 		}
// 	}
// }
