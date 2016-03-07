package gameobjsys

type (
	Component struct {
		Action     Action                          `config:"action"`
		Direction  Direction                       `config:"direction"`
		Animations map[Action]map[Direction]string `config:"animations"`
	}
)

func NewComponent(action Action, direction Direction, animations map[Action]map[Direction]string) *Component {
	return &Component{
		Action:     action,
		Direction:  direction,
		Animations: animations,
	}
}
