package entity

import (
	"errors"
	"fmt"
)

type (
	ComponentSet struct {
		components []IComponent
		idxMap     map[ID]int
	}

	IComponentSet interface {
		Add(eid ID, cmpt IComponent) error
		Remove(eid ID) error
		Visitor(ids []ID) (*ComponentSetVisitor, error)

		Len() int
		Get(idx int) IComponent
		Set(idx int, cmpt IComponent)
	}

	ComponentSetVisitor struct {
		cmptSet    IComponentSet
		currentIdx int
		indices    []int
	}
)

func NewComponentSet() *ComponentSet {
	return &ComponentSet{
		components: []IComponent{},
		idxMap:     map[ID]int{},
	}
}

func (cs *ComponentSet) Add(eid ID, cmpt IComponent) error {
	if _, exists := cs.idxMap[eid]; exists {
		return errors.New("component already exists")
	}

	cs.components = append(cs.components, cmpt)
	cs.idxMap[eid] = len(cs.components) - 1
	return nil
}

func (cs *ComponentSet) Remove(eid ID) error {
	idx, exists := cs.idxMap[eid]
	if !exists {
		return fmt.Errorf("component '%v' does not exist", eid)
	}

	cs.components = append(cs.components[:idx], cs.components[idx+1:]...)
	delete(cs.idxMap, eid)
	return nil
}

func (cs *ComponentSet) IndexOf(eid ID) (int, bool) {
	idx, exists := cs.idxMap[eid]
	return idx, exists
}

func (cs *ComponentSet) Get(idx int) IComponent {
	return cs.components[idx]
}

func (cs *ComponentSet) Set(idx int, cmpt IComponent) {
	cs.components[idx] = cmpt
}

func (cs *ComponentSet) Len() int {
	return len(cs.components)
}

func (cs *ComponentSet) Visitor(entityIDs []ID) (*ComponentSetVisitor, error) {
	indices := make([]int, len(entityIDs))
	for i, eid := range entityIDs {
		idx, exists := cs.idxMap[eid]
		if !exists {
			return nil, fmt.Errorf("entity.ComponentSet.Visitor: unknown entity ID '%v'", eid)
		}
		indices[i] = idx
	}
	return &ComponentSetVisitor{cs, 0, indices}, nil
}

func (v *ComponentSetVisitor) Advance() {
	v.currentIdx++
}

func (v *ComponentSetVisitor) Get() IComponent {
	idx := v.indices[v.currentIdx]
	return v.cmptSet.Get(idx)
}

func (v *ComponentSetVisitor) Set(cmpt IComponent) {
	idx := v.indices[v.currentIdx]
	v.cmptSet.Set(idx, cmpt)
}

func (v *ComponentSetVisitor) Len() int {
	return len(v.indices)
}

func (v *ComponentSetVisitor) Indices() []int {
	return v.indices
}
