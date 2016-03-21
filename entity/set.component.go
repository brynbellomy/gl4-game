package entity

import (
	"errors"
	"fmt"
)

type (
	ComponentSet struct {
		components IComponentSlice
		idxMap     map[ID]int
	}

	IComponentSlice interface {
		Get(idx int) (IComponent, bool)
		Set(idx int, cmpt IComponent) bool
		Append(cmpt IComponent) IComponentSlice
		Remove(idx int) IComponentSlice
	}

	IComponentSet interface {
		Add(eid ID, cmpt IComponent) error
		Remove(eid ID) error
		Get(eid ID) (IComponent, error)
		Set(eid ID, cmpt IComponent) error
		Index(eid ID) (int, bool)
		Indices(eids []ID) ([]int, error)
		IDForIndex(i int) (ID, bool)
		Slice() interface{}
	}
)

func NewComponentSet(cmptSlice IComponentSlice) IComponentSet {
	return &ComponentSet{
		components: cmptSlice,
		idxMap:     map[ID]int{},
	}
}

func (cs *ComponentSet) IDForIndex(idx int) (ID, bool) {
	for id, midx := range cs.idxMap {
		if idx == midx {
			return id, true
		}
	}
	return InvalidID, false
}

func (cs *ComponentSet) Add(eid ID, cmpt IComponent) error {
	if _, exists := cs.idxMap[eid]; exists {
		return errors.New("component already exists")
	}

	cs.components = cs.components.Append(cmpt)
	cs.idxMap[eid] = len(cs.idxMap)
	return nil
}

func (cs *ComponentSet) Remove(eid ID) error {
	idx, exists := cs.idxMap[eid]
	if !exists {
		return fmt.Errorf("component '%v' does not exist", eid)
	}

	cs.components = cs.components.Remove(idx)

	for key, curidx := range cs.idxMap {
		if curidx > idx {
			cs.idxMap[key] = curidx - 1
		}
	}

	delete(cs.idxMap, eid)
	return nil
}

func (cs *ComponentSet) Get(eid ID) (IComponent, error) {
	idx, exists := cs.idxMap[eid]
	if !exists {
		return nil, fmt.Errorf("entity.ComponentSet.Get: unknown entity ID '%v'", eid)
	}

	cmpt, exists := cs.components.Get(idx)
	if !exists {
		return nil, fmt.Errorf("entity.ComponentSet.Get: unknown component index '%v'", idx)
	}

	return cmpt, nil
}

func (cs *ComponentSet) Set(eid ID, cmpt IComponent) error {
	idx, exists := cs.idxMap[eid]
	if !exists {
		return fmt.Errorf("entity.ComponentSet.Set: unknown entity ID '%v'", eid)
	}

	exists = cs.components.Set(idx, cmpt)
	if !exists {
		return fmt.Errorf("entity.ComponentSet.Set: unknown entity ID '%v'", eid)
	}

	return nil
}

func (cs *ComponentSet) Index(eid ID) (int, bool) {
	i, exists := cs.idxMap[eid]
	return i, exists
}

func (cs *ComponentSet) Indices(entityIDs []ID) ([]int, error) {
	indices := make([]int, len(entityIDs))
	for i, eid := range entityIDs {
		idx, exists := cs.Index(eid)
		if !exists {
			return nil, fmt.Errorf("entity.ComponentSet.Indices: unknown entity ID '%v'", eid)
		}
		indices[i] = idx
	}
	return indices, nil
}

func (cs *ComponentSet) Slice() interface{} {
	return cs.components
}
