package triggersys

import (
	"github.com/brynbellomy/gl4-game/entity"
)

type (
	Component struct {
		Condition ICondition `config:"condition"`
		Effect    IEffect    `config:"effect"`

		entity.ComponentKind `config:"-"`
	}
)

// func (c *Component) GetTag() string {
// 	return c.Tag
// }

// func (c *Component) SetTag(tag string) {
// 	c.Tag = tag
// }

func (c *Component) Clone() entity.IComponent {
	x := *c
	return &x
}
