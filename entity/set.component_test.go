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

	MassCmpt struct {
		Mass int
		entity.ComponentKind
	}
)

func (c *NameCmpt) Clone() entity.IComponent {
	return &NameCmpt{Name: c.Name, ComponentKind: c.ComponentKind}
}

func (c *MassCmpt) Clone() entity.IComponent {
	return &MassCmpt{Mass: c.Mass, ComponentKind: c.ComponentKind}
}

var _ = Describe("ComponentSet", func() {
	Context("when components are added", func() {
		var cmptSet *entity.ComponentSet

		expected := []*NameCmpt{
			{"bryn", entity.ComponentKind(0)},
			{"duke nukem", entity.ComponentKind(0)},
			{"lo wang", entity.ComponentKind(0)},
		}

		BeforeEach(func() {
			cmptSet = entity.NewComponentSet()
			for i, cmpt := range expected {
				cmptSet.Add(entity.ID(i), cmpt)
			}
		})

		It("should return the proper .Len", func() {
			Expect(cmptSet.Len()).To(Equal(3))
		})

		It("should return the proper .IndexOf for each entity ID", func() {
			i, exists := cmptSet.IndexOf(entity.ID(0))
			Expect(i).To(Equal(0))
			Expect(exists).To(BeTrue())

			i, exists = cmptSet.IndexOf(entity.ID(1))
			Expect(i).To(Equal(1))
			Expect(exists).To(BeTrue())

			i, exists = cmptSet.IndexOf(entity.ID(2))
			Expect(i).To(Equal(2))
			Expect(exists).To(BeTrue())
		})

		It("should return the proper component when .Get is called", func() {
			for i := range expected {
				Expect(cmptSet.Get(i)).To(Equal(expected[i]))
			}
		})
	})

	Context(".Visitor", func() {
		expected := []*NameCmpt{
			{"bryn", entity.ComponentKind(0)},
			{"duke nukem", entity.ComponentKind(0)},
			{"lo wang", entity.ComponentKind(0)},
		}

		var cmptSet *entity.ComponentSet
		var v *entity.ComponentSetVisitor
		var err error

		BeforeEach(func() {
			cmptSet = entity.NewComponentSet()
			for i, cmpt := range expected {
				cmptSet.Add(entity.ID(i), cmpt)
			}

			v, err = cmptSet.Visitor([]entity.ID{1, 2})
			if err != nil {
				Fail(err.Error())
			}
		})

		It("should return the proper length", func() {
			Expect(v.Len()).To(Equal(2))
		})

		It("should select the proper ComponentSet indices", func() {
			Expect(v.Indices()).To(Equal([]int{1, 2}))
		})

		It("should enumerate items in the proper order", func() {
			Expect(v.Get().(*NameCmpt).Name).To(Equal("duke nukem"))
		})
	})
})
