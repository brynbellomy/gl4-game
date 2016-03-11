package entity

type (
	ComponentSet ComponentKind
)

func (c ComponentSet) Has(n ComponentKind) bool {
	return c&ComponentSet(n) > 0
}

func (c ComponentSet) HasAll(n ComponentSet) bool {
	return c&n > 0
}

// func (c ComponentSet) HasN(n uint) bool {
// 	return *c&(1<<n) > 0
// }

// func (c ComponentSet) And(other ComponentSet) ComponentSet {
// 	return c & other
// }

func (c ComponentSet) Add(other ComponentKind) ComponentSet {
	return c | ComponentSet(other)
}

// func (c *ComponentSet) Set(n uint, val bool) {
// 	if val == true {
// 		*c = *c | (1 << n)
// 	} else {
// 		*c = *c & ^(1 << n)
// 	}
// }
