package entity_test

import (
	"github.com/brynbellomy/gl4-game/entity"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type (
	NameCmpt struct {
		Name string
		entity.ComponentKind
	}

	NameCmptSlice []NameCmpt
)

func (c NameCmpt) Clone() entity.IComponent {
	return NameCmpt{Name: c.Name, ComponentKind: c.ComponentKind}
}

func (c *NameCmpt) SetName(n string) {
	c.Name = n
}

func (cs NameCmptSlice) Append(cmpt entity.IComponent) entity.IComponentSlice {
	return append(cs, cmpt.(NameCmpt))
}

func (cs NameCmptSlice) Remove(idx int) entity.IComponentSlice {
	return append(cs[:idx], cs[idx+1:]...)
}

var _ = Describe("ComponentSet", func() {
	Context("when components are added", func() {
		var cmptSet entity.IComponentSet

		expected := []NameCmpt{
			{"bryn", entity.ComponentKind(0)},
			{"duke nukem", entity.ComponentKind(0)},
			{"lo wang", entity.ComponentKind(0)},
		}

		BeforeEach(func() {
			cmptSet = entity.NewComponentSet(NameCmptSlice{})
			for i, cmpt := range expected {
				err := cmptSet.Add(entity.ID(i), cmpt)
				if err != nil {
					Fail(err.Error())
				}
			}
		})

		It("should return the proper .Indices for each entity ID", func() {
			idxs, err := cmptSet.Indices([]entity.ID{0, 1, 2})
			if err != nil {
				Fail(err.Error())
			}

			Expect(idxs[0]).To(Equal(0))
			Expect(idxs[1]).To(Equal(1))
			Expect(idxs[2]).To(Equal(2))
		})

		It("should keep track of the proper .Indices after components are .Removed", func() {
			cmptSet.Remove(entity.ID(1))

			idxs, err := cmptSet.Indices([]entity.ID{0, 2})
			if err != nil {
				Fail(err.Error())
			}

			Expect(idxs[0]).To(Equal(0))
			Expect(idxs[1]).To(Equal(1))
		})

		It("should return an error if asked for .Indices of .Removed entities", func() {
			cmptSet.Remove(entity.ID(1))
			_, err := cmptSet.Indices([]entity.ID{1})
			Expect(err).To(Not(BeNil()))
		})
	})

	// Context(".Visitor", func() {
	// 	expected := []*NameCmpt{
	// 		{"bryn", entity.ComponentKind(0)},
	// 		{"duke nukem", entity.ComponentKind(0)},
	// 		{"lo wang", entity.ComponentKind(0)},
	// 	}

	// 	var cmptSet *entity.ComponentSet
	// 	var v *entity.ComponentSetVisitor
	// 	var err error

	// 	BeforeEach(func() {
	// 		cmptSet = entity.NewComponentSet()
	// 		for i, cmpt := range expected {
	// 			cmptSet.Add(entity.ID(i), cmpt)
	// 		}

	// 		v, err = cmptSet.Visitor([]entity.ID{1, 2})
	// 		if err != nil {
	// 			Fail(err.Error())
	// 		}
	// 	})

	// 	It("should return the proper length", func() {
	// 		Expect(v.Len()).To(Equal(2))
	// 	})

	// 	It("should select the proper ComponentSet indices", func() {
	// 		Expect(v.Indices()).To(Equal([]int{1, 2}))
	// 	})

	// 	It("should enumerate items in the proper order", func() {
	// 		Expect(v.Get().(*NameCmpt).Name).To(Equal("duke nukem"))
	// 	})

	// 	It("should .Set items correctly", func() {
	// 		cmpt := v.Get().(NameCmpt)
	// 		cmpt.SetName("brynsk")
	// 		v.Set(cmpt)

	// 		Expect(cmptSet.Get(0)).To(Equal(NameCmpt{"brynsk", entity.ComponentKind(0)}))
	// 	})
	// })
})
