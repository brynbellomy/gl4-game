package entity

type (
	ComponentMask uint64
)

func (c ComponentMask) Has(n ComponentKind) bool {
	return (c & ComponentMask(n)) > 0
}

func (c ComponentMask) HasAll(n ComponentMask) bool {
	return (c & n) == n
}

// func (c ComponentMask) HasN(n uint) bool {
// 	return *c&(1<<n) > 0
// }

// func (c ComponentMask) And(other ComponentMask) ComponentMask {
// 	return c & other
// }

func (c ComponentMask) Add(other ComponentKind) ComponentMask {
	return c | ComponentMask(other)
}

// func (c *ComponentMask) Set(n uint, val bool) {
// 	if val == true {
// 		*c = *c | (1 << n)
// 	} else {
// 		*c = *c & ^(1 << n)
// 	}
// }
