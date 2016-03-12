package entity

type (
	ComponentMask ComponentKind
)

func (c ComponentMask) Has(n ComponentKind) bool {
	return c&n.KindMask() > 0
}

func (c ComponentMask) HasAll(n ComponentMask) bool {
	return c&n > 0
}

// func (c ComponentMask) HasN(n uint) bool {
// 	return *c&(1<<n) > 0
// }

// func (c ComponentMask) And(other ComponentMask) ComponentMask {
// 	return c & other
// }

func (c ComponentMask) Add(other ComponentKind) ComponentMask {
	return c | other.KindMask()
}

// func (c *ComponentMask) Set(n uint, val bool) {
// 	if val == true {
// 		*c = *c | (1 << n)
// 	} else {
// 		*c = *c & ^(1 << n)
// 	}
// }
