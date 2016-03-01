package gameobjsys

type (
	Component struct {
		action     Action
		direction  Direction
		animations map[Action]map[Direction]string
	}
)

func NewComponent(action Action, direction Direction, animations map[Action]map[Direction]string) *Component {
	return &Component{
		action:     action,
		direction:  direction,
		animations: animations,
	}
}
